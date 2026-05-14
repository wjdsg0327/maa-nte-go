package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"study/one/service"

	"github.com/gin-gonic/gin"
)

const pipelineDir = "./resource/pipeline"

func validatePipelineName(name string) error {
	if name == "" || strings.ContainsAny(name, "/\\:*?\"<>|") || name == "." || name == ".." {
		return http.ErrNotSupported
	}
	return nil
}

// ListPipelines 获取所有 pipeline 列表
func ListPipelines(c *gin.Context) {
	files, err := os.ReadDir(pipelineDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pipelines []map[string]interface{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			name := strings.TrimSuffix(f.Name(), ".json")
			pipelines = append(pipelines, map[string]interface{}{
				"name": name,
				"file": f.Name(),
			})
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
	filePath := filepath.Join(pipelineDir, name+".json")

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

	// 文件名安全检查
	if err := validatePipelineName(req.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}

	filePath := filepath.Join(pipelineDir, req.Name+".json")

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
	filePath := filepath.Join(pipelineDir, name+".json")

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
	filePath := filepath.Join(pipelineDir, name+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载资源
	if err := service.Service.ReloadResources(); err != nil {
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
	dst := filepath.Join("./resource/image", filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded", "filename": filename})
}

// ListImages 获取图片列表
func ListImages(c *gin.Context) {
	files, err := os.ReadDir("./resource/image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var images []string
	for _, f := range files {
		if !f.IsDir() {
			images = append(images, f.Name())
		}
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
	ocrDir := "./resource/ocr"
	files, err := os.ReadDir(ocrDir)
	if err != nil {
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
	detectDir := "./resource/detect"
	files, err := os.ReadDir(detectDir)
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
