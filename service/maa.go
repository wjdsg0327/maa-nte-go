package service

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/MaaXYZ/maa-framework-go/v4"
	"github.com/MaaXYZ/maa-framework-go/v4/controller/win32"
)

type MaaService struct {
	mu        sync.Mutex
	tasker    *maa.Tasker
	controller *maa.Controller
	resource  *maa.Resource
	inited    bool
	running   bool
	lastTask  string
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

	if err := maa.Init(
		maa.WithLibDir(libDir),
		maa.WithStdoutLevel(maa.LoggingLevelInfo),
		maa.WithDebugMode(true),
	); err != nil {
		return fmt.Errorf("maa 初始化失败: %v", err)
	}

	userDir, _ := filepath.Abs(".")
	if err := maa.ConfigInitOption(userDir, "{}"); err != nil {
		return fmt.Errorf("工具包配置失败: %v", err)
	}

	// 加载资源
	s.resource, err = maa.NewResource()
	if err != nil {
		return fmt.Errorf("创建资源失败: %v", err)
	}

	resPath, _ := filepath.Abs("./resource")
	s.resource.PostPipeline(filepath.Join(resPath, "pipeline")).Wait()
	s.resource.PostImage(filepath.Join(resPath, "image")).Wait()
	s.resource.PostOcrModel(filepath.Join(resPath, "ocr")).Wait()

	// 创建 Tasker
	s.tasker, err = maa.NewTasker()
	if err != nil {
		return fmt.Errorf("创建 Tasker 失败: %v", err)
	}

	s.tasker.BindResource(s.resource)

	s.inited = true
	log.Println("MaaFramework 初始化成功")
	return nil
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

	s.controller.PostConnect().Wait()
	s.tasker.BindController(s.controller)

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

	resPath, _ := filepath.Abs("./resource")
	s.resource.Clear()
	s.resource.PostPipeline(filepath.Join(resPath, "pipeline")).Wait()
	s.resource.PostImage(filepath.Join(resPath, "image")).Wait()
	s.resource.PostOcrModel(filepath.Join(resPath, "ocr")).Wait()

	log.Println("资源已重新加载")
	return nil
}

// ExecuteTask 执行任务
func (s *MaaService) ExecuteTask(taskName string) (map[string]interface{}, error) {
	s.mu.Lock()
	if !s.inited || s.tasker == nil {
		s.mu.Unlock()
		return nil, errors.New("MaaFramework 未初始化")
	}
	if s.controller == nil || !s.controller.Connected() {
		s.mu.Unlock()
		return nil, errors.New("请先连接窗口")
	}
	if s.running {
		s.mu.Unlock()
		return nil, errors.New("任务正在执行中")
	}
	s.running = true
	s.lastTask = taskName
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	log.Printf("开始执行任务: %s", taskName)
	taskJob := s.tasker.PostTask(taskName)
	status := taskJob.Wait()

	// 获取结果
	result := map[string]interface{}{
		"task":   taskName,
		"status": fmt.Sprintf("%v", status),
	}

	// 获取节点详情
	nodeDetail, err := s.tasker.GetLatestNode(taskName)
	if err == nil && nodeDetail != nil {
		nodeInfo := map[string]interface{}{
			"name":    nodeDetail.Name,
			"success": nodeDetail.RunCompleted,
		}

		// 识别结果
		if nodeDetail.Recognition != nil && nodeDetail.Recognition.Results != nil {
			var allResults []map[string]interface{}
			for _, r := range nodeDetail.Recognition.Results.All {
				if ocr, ok := r.AsOCR(); ok {
					allResults = append(allResults, map[string]interface{}{
						"type":  "OCR",
						"text":  ocr.Text,
						"score": ocr.Score,
						"box":   ocr.Box,
					})
				} else if tm, ok := r.AsTemplateMatch(); ok {
					allResults = append(allResults, map[string]interface{}{
						"type":  "TemplateMatch",
						"score": tm.Score,
						"box":   tm.Box,
					})
				}
			}
			nodeInfo["results"] = allResults
		}

		result["node"] = nodeInfo
	}

	log.Printf("任务完成: %s, 状态: %v", taskName, status)
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
