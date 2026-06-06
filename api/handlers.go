package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"study/one/service"

	"github.com/gin-gonic/gin"
)

func pipelineDir() string {
	return service.Service.ResourcePath("pipeline")
}

func imageDir() string {
	return service.Service.ResourcePath("image")
}

func ocrDir() string {
	return service.Service.ResourcePath("ocr")
}

func detectDir() string {
	return service.Service.ResourcePath("detect")
}

func validatePipelineName(name string) error {
	if name == "" || strings.ContainsAny(name, "/\\:*?\"<>|") || name == "." || name == ".." {
		return http.ErrNotSupported
	}
	return nil
}

func validatePipelineContent(content map[string]interface{}) error {
	return validatePipelineValue(content, "")
}

func validatePipelineValue(value interface{}, path string) error {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, child := range typed {
			childPath := joinPipelinePath(path, key)
			switch key {
			case "roi":
				if err := validateRectField(child, true); err != nil {
					return fmt.Errorf("%s: %w", childPath, err)
				}
			case "roi_offset":
				if err := validateRectField(child, false); err != nil {
					return fmt.Errorf("%s: %w", childPath, err)
				}
			case "wait_freezes":
				return fmt.Errorf("%s: wait_freezes is not a Maa route field; use pre_wait_freezes, post_wait_freezes, or repeat_wait_freezes for delays", childPath)
			case "reverse":
				return fmt.Errorf("%s: reverse is not a Maa route field; use on_error for failure routing", childPath)
			}
			if err := validatePipelineValue(child, childPath); err != nil {
				return err
			}
		}
	case []interface{}:
		for idx, child := range typed {
			if err := validatePipelineValue(child, fmt.Sprintf("%s[%d]", path, idx)); err != nil {
				return err
			}
		}
	}
	return nil
}

func joinPipelinePath(path string, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

func validateRectField(value interface{}, allowString bool) error {
	if value == nil {
		return nil
	}

	if text, ok := value.(string); ok {
		if allowString && strings.TrimSpace(text) != "" {
			return nil
		}
		if allowString {
			return fmt.Errorf("must be a non-empty string reference or [x,y,w,h] with 4 integer values")
		}
		return fmt.Errorf("must be [x,y,w,h] with 4 integer values")
	}

	arr, ok := value.([]interface{})
	if !ok || len(arr) != 4 {
		return fmt.Errorf("must be [x,y,w,h] with 4 integer values")
	}

	for _, item := range arr {
		num, ok := numberValue(item)
		if !ok || math.Trunc(num) != num {
			return fmt.Errorf("must be [x,y,w,h] with 4 integer values")
		}
	}
	return nil
}

func numberValue(value interface{}) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, isFiniteNumber(typed)
	case float32:
		num := float64(typed)
		return num, isFiniteNumber(num)
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case json.Number:
		num, err := typed.Float64()
		return num, err == nil && isFiniteNumber(num)
	default:
		return 0, false
	}
}

func isFiniteNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func rollbackPipelineFile(filePath string, backup []byte, existed bool) error {
	if existed {
		if err := os.WriteFile(filePath, backup, 0644); err != nil {
			return err
		}
	} else if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return service.Service.ReloadResources()
}

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, service.Service.ConfigSnapshot())
}

func ReloadResources(c *gin.Context) {
	if err := service.Service.ReloadResources(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Resources reloaded"})
}

// ListPipelines 获取所有 pipeline 列表
func ListPipelines(c *gin.Context) {
	files, err := os.ReadDir(pipelineDir())
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"pipelines": []map[string]interface{}{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pipelines []map[string]interface{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			name := strings.TrimSuffix(f.Name(), ".json")
			item := map[string]interface{}{
				"name": name,
				"file": f.Name(),
			}
			if entry, err := service.ResolvePipelineEntryName(name); err == nil {
				item["entry"] = entry
			}
			pipelines = append(pipelines, item)
		}
	}

	c.JSON(http.StatusOK, gin.H{"pipelines": pipelines})
}

// GetPipeline 获取单个 pipeline 内容
func GetPipeline(c *gin.Context) {
	name := c.Param("name")
	if err := validatePipelineName(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}
	filePath := filepath.Join(pipelineDir(), name+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON format"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":    name,
		"content": content,
	})
}

// CreatePipeline 创建新 pipeline
func CreatePipeline(c *gin.Context) {
	var req struct {
		Name    string                 `json:"name" binding:"required"`
		Content map[string]interface{} `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validatePipelineContent(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline content: " + err.Error()})
		return
	}

	// 文件名安全检查
	if err := validatePipelineName(req.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}

	if err := os.MkdirAll(pipelineDir(), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filePath := filepath.Join(pipelineDir(), req.Name+".json")

	// 检查是否已存在
	if _, err := os.Stat(filePath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Pipeline already exists"})
		return
	}

	data, err := json.MarshalIndent(req.Content, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载资源
	if err := service.Service.ReloadResources(); err != nil {
		_ = rollbackPipelineFile(filePath, nil, false)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pipeline created but resource reload failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pipeline created", "name": req.Name})
}

// UpdatePipeline 更新 pipeline
func UpdatePipeline(c *gin.Context) {
	name := c.Param("name")
	if err := validatePipelineName(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}
	filePath := filepath.Join(pipelineDir(), name+".json")

	// 检查是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	var req struct {
		Content map[string]interface{} `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validatePipelineContent(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline content: " + err.Error()})
		return
	}

	data, err := json.MarshalIndent(req.Content, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	backup, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载资源
	if err := service.Service.ReloadResources(); err != nil {
		_ = rollbackPipelineFile(filePath, backup, true)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pipeline updated but resource reload failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline updated", "name": name})
}

// DeletePipeline 删除 pipeline
func DeletePipeline(c *gin.Context) {
	name := c.Param("name")
	if err := validatePipelineName(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}
	filePath := filepath.Join(pipelineDir(), name+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	backup, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载资源
	if err := service.Service.ReloadResources(); err != nil {
		_ = rollbackPipelineFile(filePath, backup, true)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pipeline deleted but resource reload failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline deleted", "name": name})
}

// ExecuteTask 执行任务
func ExecuteTask(c *gin.Context) {
	var req struct {
		TaskName string `json:"task" binding:"required"`
		NodeName string `json:"node"` // 可选：指定从哪个节点开始执行
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.Service.StartTask(req.TaskName, req.NodeName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NodeName != "" {
		c.JSON(http.StatusAccepted, gin.H{"message": "Task started", "task": req.TaskName, "node": req.NodeName})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"message": "Task started", "task": req.TaskName})
	}
}

// RunTask 同步执行任务并返回 Maa 任务/节点调试详情
func RunTask(c *gin.Context) {
	var req struct {
		TaskName string `json:"task" binding:"required"`
		NodeName string `json:"node"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := service.Service.ExecuteTask(req.TaskName, req.NodeName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTaskStatus 获取任务状态
func GetTaskStatus(c *gin.Context) {
	status := service.Service.GetStatus()
	c.JSON(http.StatusOK, status)
}

// StopTask 停止任务
func StopTask(c *gin.Context) {
	if err := service.Service.StopTask(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stop signal sent"})
}

// GetWindows 获取窗口列表
func GetWindows(c *gin.Context) {
	windows, err := service.Service.GetWindowList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"windows": windows})
}

// ConnectWindow 连接窗口
func ConnectWindow(c *gin.Context) {
	var req struct {
		Handle uintptr `json:"handle" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.Service.ConnectWindow(req.Handle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connected"})
}

// GetScreenshot 获取当前连接窗口截图，用于 ROI 框选
func GetScreenshot(c *gin.Context) {
	capture, err := service.Service.CaptureScreenshot()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, capture)
}

// ServeIndex 提供前端页面
func ServeIndex(c *gin.Context) {
	c.File("./web/index.html")
}

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	filename := filepath.Base(file.Filename)
	if filename == "." || filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}

	dir := imageDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dst := filepath.Join(dir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded", "filename": filename})
}

// ListImages 获取图片列表
func ListImages(c *gin.Context) {
	var images []string
	dir := imageDir()
	if err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		images = append(images, filepath.ToSlash(rel))
		return nil
	}); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"images": []string{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// GetNodeList 获取节点列表
func GetNodeList(c *gin.Context) {
	nodes, err := service.Service.GetNodeList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

// ReadFile 读取文件内容
func ReadFile(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	json.Unmarshal(data, &req)

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path required"})
		return
	}
}

// ListOcrModels 获取OCR模型列表
func ListOcrModels(c *gin.Context) {
	files, err := os.ReadDir(ocrDir())
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"models": []string{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var models []string
	for _, f := range files {
		if f.IsDir() {
			// 子文件夹作为模型名
			models = append(models, f.Name())
		} else if strings.HasSuffix(f.Name(), ".onnx") {
			// 根目录的 onnx 文件
			models = append(models, f.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// ListDetectModels 获取Detect模型列表
func ListDetectModels(c *gin.Context) {
	files, err := os.ReadDir(detectDir())
	if err != nil {
		// 目录不存在返回空列表
		c.JSON(http.StatusOK, gin.H{"models": []string{}})
		return
	}

	var models []string
	for _, f := range files {
		if f.IsDir() {
			models = append(models, f.Name())
		} else if strings.HasSuffix(f.Name(), ".onnx") {
			models = append(models, f.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}
