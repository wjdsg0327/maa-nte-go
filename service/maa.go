package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"github.com/MaaXYZ/maa-framework-go/v4"
	"github.com/MaaXYZ/maa-framework-go/v4/controller/win32"
)

type MaaService struct {
	mu         sync.Mutex
	tasker     *maa.Tasker
	controller *maa.Controller
	resource   *maa.Resource
	inited     bool
	running    bool
	lastTask   string
	lastError  string
}

var Service = &MaaService{}

// Init 初始化 MaaFramework
func (s *MaaService) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.inited {
		return nil
	}

	libDir, err := filepath.Abs("./MaaFramework/bin")
	if err != nil {
		return fmt.Errorf("解析库目录失败: %v", err)
	}

	// 设置日志目录，带标记的识别图片会保存到 log/vision 目录
	logDir, _ := filepath.Abs("./log")
	// 确保日志目录存在
	os.MkdirAll(logDir, 0755)
	os.MkdirAll(filepath.Join(logDir, "vision"), 0755)

	if err := maa.Init(
		maa.WithLibDir(libDir),
		maa.WithLogDir(logDir),
		maa.WithSaveDraw(true), // 启用保存带标记的识别结果图片
		maa.WithStdoutLevel(maa.LoggingLevelInfo),
		maa.WithDebugMode(true),
	); err != nil {
		return fmt.Errorf("maa 初始化失败: %v", err)
	}

	visionDir := filepath.Join(logDir, "vision")
	log.Printf("识别结果图片将保存到: %s", visionDir)

	userDir, _ := filepath.Abs(".")
	if err := maa.ConfigInitOption(userDir, "{}"); err != nil {
		return fmt.Errorf("工具包配置失败: %v", err)
	}

	// 加载资源
	s.resource, err = maa.NewResource()
	if err != nil {
		return fmt.Errorf("创建资源失败: %v", err)
	}
	if err := s.loadResourcesLocked(); err != nil {
		s.resource.Destroy()
		s.resource = nil
		return err
	}

	// 创建 Tasker
	s.tasker, err = maa.NewTasker()
	if err != nil {
		return fmt.Errorf("创建 Tasker 失败: %v", err)
	}

	if err := s.tasker.BindResource(s.resource); err != nil {
		s.tasker.Destroy()
		s.tasker = nil
		s.resource.Destroy()
		s.resource = nil
		return fmt.Errorf("绑定资源失败: %v", err)
	}

	s.inited = true
	log.Println("MaaFramework 初始化成功")

	// 自动连接目标窗口
	go s.autoConnectTargetWindow("异环")

	return nil
}

// autoConnectTargetWindow 自动连接指定名称的窗口
func (s *MaaService) autoConnectTargetWindow(windowTitle string) {
	// 等待一下让初始化完全完成
	// time.Sleep(time.Second)

	s.mu.Lock()
	if !s.inited {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	windows, err := maa.FindDesktopWindows()
	if err != nil {
		log.Printf("获取窗口列表失败: %v", err)
		return
	}

	for _, w := range windows {
		// 去除空格后匹配
		if strings.TrimSpace(w.WindowName) == windowTitle {
			log.Printf("找到目标窗口: %s (句柄: %v)", w.WindowName, w.Handle)
			if err := s.ConnectWindow(uintptr(w.Handle)); err != nil {
				log.Printf("自动连接窗口失败: %v", err)
			} else {
				log.Printf("已自动连接到窗口: %s", w.WindowName)
			}
			return
		}
	}

	log.Printf("未找到名为 '%s' 的窗口", windowTitle)
}

// ConnectWindow 连接到指定窗口
func (s *MaaService) ConnectWindow(hWnd uintptr) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.inited {
		return errors.New("MaaFramework 未初始化")
	}

	// 销毁旧控制器
	if s.controller != nil {
		s.controller.Destroy()
	}

	var err error
	s.controller, err = maa.NewWin32Controller(
		unsafe.Pointer(hWnd),
		win32.ScreencapPrintWindow,
		win32.InputSeize,
		win32.InputSeize,
	)
	if err != nil {
		return fmt.Errorf("创建控制器失败: %v", err)
	}

	if err := waitMaaJob("连接窗口", s.controller.PostConnect()); err != nil {
		s.controller.Destroy()
		s.controller = nil
		return err
	}
	if !s.controller.Connected() {
		s.controller.Destroy()
		s.controller = nil
		return errors.New("控制器连接失败")
	}
	if err := s.tasker.BindController(s.controller); err != nil {
		s.controller.Destroy()
		s.controller = nil
		return fmt.Errorf("绑定控制器失败: %v", err)
	}

	if !s.tasker.Initialized() {
		return errors.New("任务器初始化失败")
	}

	log.Printf("已连接到窗口: %v", hWnd)
	return nil
}

// GetWindowList 获取桌面窗口列表
func (s *MaaService) GetWindowList() ([]map[string]interface{}, error) {
	if !s.inited {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	windows, err := maa.FindDesktopWindows()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, w := range windows {
		if w.WindowName != "" {
			result = append(result, map[string]interface{}{
				"handle": uintptr(w.Handle),
				"class":  w.ClassName,
				"title":  w.WindowName,
			})
		}
	}
	return result, nil
}

// GetNodeList 获取所有任务节点
func (s *MaaService) GetNodeList() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.inited || s.resource == nil {
		return nil, errors.New("MaaFramework 未初始化")
	}

	return s.resource.GetNodeList()
}

// ReloadResources 重新加载资源
func (s *MaaService) ReloadResources() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.inited || s.resource == nil {
		return errors.New("MaaFramework 未初始化")
	}

	if err := s.resource.Clear(); err != nil {
		return fmt.Errorf("清理资源失败: %v", err)
	}
	if err := s.loadResourcesLocked(); err != nil {
		return err
	}

	log.Println("资源已重新加载")
	return nil
}

func (s *MaaService) loadResourcesLocked() error {
	resPath, err := filepath.Abs("./resource")
	if err != nil {
		return fmt.Errorf("解析资源目录失败: %v", err)
	}
	if _, err := os.Stat(resPath); err != nil {
		return fmt.Errorf("资源目录不可用: %v", err)
	}

	// Prefer the official bundle loading path so default_pipeline.json and the
	// standard resource layout are handled consistently by MaaFramework.
	if err := waitMaaJob("加载资源包", s.resource.PostBundle(resPath)); err != nil {
		return err
	}

	// Keep compatibility with the current project layout, where OCR models live
	// in resource/ocr instead of the documented resource/model/ocr directory.
	ocrPath := filepath.Join(resPath, "ocr")
	if info, err := os.Stat(ocrPath); err == nil && info.IsDir() {
		if err := waitMaaJob("加载 OCR 模型", s.resource.PostOcrModel(ocrPath)); err != nil {
			return err
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("检查 OCR 模型目录失败: %v", err)
	}

	if !s.resource.Loaded() {
		return errors.New("资源加载未完成或失败")
	}
	return nil
}

func waitMaaJob(name string, job *maa.Job) error {
	if job == nil {
		return fmt.Errorf("%s失败: Maa job 为空", name)
	}
	status := job.Wait().Status()
	if !status.Success() {
		return fmt.Errorf("%s失败: %s", name, status)
	}
	return nil
}

// ExecuteTask 执行任务
func (s *MaaService) ExecuteTask(taskName string, nodeName string) (result map[string]interface{}, err error) {
	if err := s.reserveTask(taskName); err != nil {
		return nil, err
	}
	defer func() {
		s.finishTask(err)
	}()

	return s.executeReservedTask(taskName, nodeName)
}

func (s *MaaService) StartTask(taskName string, nodeName string) error {
	if err := s.reserveTask(taskName); err != nil {
		return err
	}

	go func() {
		_, err := s.executeReservedTask(taskName, nodeName)
		if err != nil {
			log.Printf("任务执行失败: %v", err)
		}
		s.finishTask(err)
	}()

	return nil
}

func (s *MaaService) reserveTask(taskName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.inited || s.tasker == nil {
		return errors.New("MaaFramework 未初始化")
	}
	if s.controller == nil || !s.controller.Connected() {
		return errors.New("请先连接窗口")
	}
	if s.running {
		return errors.New("任务正在执行中")
	}

	s.running = true
	s.lastTask = taskName
	s.lastError = ""
	return nil
}

func (s *MaaService) finishTask(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.running = false
	if err != nil {
		s.lastError = err.Error()
	}
}

func (s *MaaService) executeReservedTask(taskName string, nodeName string) (map[string]interface{}, error) {
	// 如果指定了节点名，使用节点名作为入口执行
	entryTask := taskName
	if nodeName != "" {
		entryTask = nodeName
		log.Printf("开始执行节点: %s (pipeline: %s)", nodeName, taskName)
	} else {
		log.Printf("开始执行任务: %s", taskName)
	}

	taskJob := s.tasker.PostTask(entryTask)
	taskJob.Wait()
	if err := taskJob.Error(); err != nil {
		return nil, err
	}
	status := taskJob.Status()

	// 获取结果
	result := map[string]interface{}{
		"task":   taskName,
		"status": status.String(),
	}

	// 获取节点详情 - 注意：GetLatestNode 对中文节点名可能有编码问题
	// 使用 TaskJob.GetDetail() 获取任务详情，再获取节点信息
	taskDetail, err := taskJob.GetDetail()
	if err == nil && taskDetail != nil && len(taskDetail.Nodes) > 0 {
		// 获取最后一个节点的详情
		lastNodeRef := taskDetail.Nodes[len(taskDetail.Nodes)-1]
		nodeDetail, err := lastNodeRef.GetDetail()
		if err == nil && nodeDetail != nil {
			nodeInfo := map[string]interface{}{
				"name":    nodeDetail.Name,
				"success": nodeDetail.RunCompleted,
			}

			// 识别结果
			if nodeDetail.Recognition != nil && nodeDetail.Recognition.Results != nil {
				log.Printf("[%s] 识别算法: %s, 结果数量: %d", nodeDetail.Name, nodeDetail.Recognition.Algorithm, len(nodeDetail.Recognition.Results.All))

				var allResults []map[string]interface{}
				for i, r := range nodeDetail.Recognition.Results.All {
					if ocr, ok := r.AsOCR(); ok {
						allResults = append(allResults, map[string]interface{}{
							"type":  "OCR",
							"text":  ocr.Text,
							"score": ocr.Score,
							"box":   ocr.Box,
						})
						log.Printf("[%s] OCR结果[%d]: text=%s, score=%.4f, box=(%d,%d,%d,%d)",
							nodeDetail.Name, i, ocr.Text, ocr.Score,
							ocr.Box.X(), ocr.Box.Y(), ocr.Box.Width(), ocr.Box.Height())
					} else if tm, ok := r.AsTemplateMatch(); ok {
						allResults = append(allResults, map[string]interface{}{
							"type":  "TemplateMatch",
							"score": tm.Score,
							"box":   tm.Box,
						})
						log.Printf("[%s] 模板匹配结果[%d]: score=%.4f", nodeDetail.Name, i, tm.Score)
					} else {
						log.Printf("[%s] 其他识别结果[%d]: type=%s", nodeDetail.Name, i, r.Type())
					}
				}
				nodeInfo["results"] = allResults

				// 使用 OCR 处理模块打印 OCR 结果
				ocrResults := ExtractOCRResults(nodeDetail.Recognition)
				if len(ocrResults) > 0 {
					PrintOCRResults(ocrResults, nodeDetail.Name)
				} else {
					log.Printf("[%s] 无 OCR 识别结果", nodeDetail.Name)
				}

				// 保存带标记的图片
				if err := SaveDrawImage(nodeDetail.Recognition, nodeDetail.Name); err != nil {
					log.Printf("[%s] 保存标记图片失败: %v", nodeDetail.Name, err)
				}
			} else {
				log.Printf("[%s] 无识别结果或识别结果为空", nodeDetail.Name)
			}

			result["node"] = nodeInfo
		}
	} else {
		log.Printf("获取节点详情失败 (可能是因为中文节点名编码问题): %v", err)
	}

	log.Printf("任务完成: %s, 状态: %v", entryTask, status)

	return result, nil
}

// IsRunning 是否正在执行任务
func (s *MaaService) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// StopTask 停止任务
func (s *MaaService) StopTask() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || s.tasker == nil {
		return errors.New("没有正在执行的任务")
	}

	s.tasker.PostStop()
	log.Println("已发送停止信号")
	return nil
}

// GetStatus 获取当前状态
func (s *MaaService) GetStatus() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return map[string]interface{}{
		"inited":    s.inited,
		"connected": s.controller != nil && s.controller.Connected(),
		"running":   s.running,
		"lastTask":  s.lastTask,
		"lastError": s.lastError,
	}
}

// Close 清理资源
func (s *MaaService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.controller != nil {
		s.controller.Destroy()
		s.controller = nil
	}
	if s.resource != nil {
		s.resource.Destroy()
		s.resource = nil
	}
	if s.tasker != nil {
		s.tasker.Destroy()
		s.tasker = nil
	}
	if s.inited {
		maa.Release()
		s.inited = false
	}
}
