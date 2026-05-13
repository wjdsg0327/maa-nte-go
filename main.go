package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/MaaXYZ/maa-framework-go/v4"
	"github.com/MaaXYZ/maa-framework-go/v4/controller/win32"

	"study/one/config"
)

func main() {
	cfg := config.DefaultConfig()

	// 将相对路径转为绝对路径
	libDir, err := filepath.Abs(cfg.LibDir)
	if err != nil {
		log.Fatalf("解析库目录失败: %v", err)
	}

	// 1. 初始化框架
	if err := maa.Init(
		maa.WithLibDir(libDir),
		maa.WithStdoutLevel(maa.LoggingLevelInfo),
		maa.WithDebugMode(true),
	); err != nil {
		log.Fatalf("maa 初始化失败: %v", err)
	}
	defer func() {
		_ = maa.Release()
	}()

	fmt.Println("MaaFramework 版本:", maa.Version())

	// 2. 配置工具包
	userDir, _ := filepath.Abs(".")
	if err := maa.ConfigInitOption(userDir, "{}"); err != nil {
		log.Fatalf("工具包配置失败: %v", err)
	}

	// 3. 获取异环窗口
	fmt.Println("正在查找桌面窗口...")
	windows, err := maa.FindDesktopWindows()
	if err != nil {
		log.Fatalf("查找窗口失败: %v", err)
	}

	var targetWindow *maa.DesktopWindow
	for _, w := range windows {
		if w.WindowName != "" {
			fmt.Printf("窗口: [%s] [%s]\n", w.ClassName, w.WindowName)
		}
		if strings.Contains(w.WindowName, "异环") {
			targetWindow = w
			break
		}
	}

	if targetWindow == nil {
		log.Fatal("未找到异环游戏窗口，请确保游戏已启动")
	}

	fmt.Printf("\n找到异环窗口:\n")
	fmt.Printf("  窗口句柄: %v\n", targetWindow.Handle)
	fmt.Printf("  类名: %s\n", targetWindow.ClassName)
	fmt.Printf("  窗口标题: %s\n", targetWindow.WindowName)

	// 4. 创建 Tasker
	tasker, err := maa.NewTasker()
	if err != nil {
		log.Fatalf("创建 Tasker 失败: %v", err)
	}
	defer tasker.Destroy()

	// 5. 创建 Win32 控制器
	ctrl, err := maa.NewWin32Controller(
		unsafe.Pointer(targetWindow.Handle),
		win32.ScreencapPrintWindow,
		win32.InputSeize,
		win32.InputSeize,
	)
	if err != nil {
		log.Fatalf("创建控制器失败: %v", err)
	}
	defer ctrl.Destroy()

	// 绑定控制器
	ctrl.PostConnect().Wait()
	if err := tasker.BindController(ctrl); err != nil {
		log.Fatalf("绑定控制器失败: %v", err)
	}

	// 6. 加载资源
	res, err := maa.NewResource()
	if err != nil {
		log.Fatalf("创建资源失败: %v", err)
	}
	defer res.Destroy()

	resPath, _ := filepath.Abs("./resource")
	fmt.Printf("加载资源路径: %s\n", resPath)

	// 使用 PostPipeline 只加载 pipeline 目录
	pipelinePath := filepath.Join(resPath, "pipeline")
	job := res.PostPipeline(pipelinePath)
	status := job.Wait()
	fmt.Printf("资源加载状态: %v\n", status)

	if !res.Loaded() {
		log.Fatal("资源加载失败")
	}

	// 获取节点列表
	nodes, err := res.GetNodeList()
	if err != nil {
		log.Fatalf("获取节点列表失败: %v", err)
	}
	fmt.Printf("已加载节点: %v\n", nodes)

	if err := tasker.BindResource(res); err != nil {
		log.Fatalf("绑定资源失败: %v", err)
	}

	if !tasker.Initialized() {
		log.Fatal("任务器初始化失败")
	}

	// 直接测试按键 - ESC = 27
	fmt.Println("\n发送 ESC 键 (KeyDown + KeyUp)...")
	ctrl.PostKeyDown(27).Wait()
	fmt.Println("KeyDown 完成")
	time.Sleep(100 * time.Millisecond)
	ctrl.PostKeyUp(27).Wait()
	fmt.Println("KeyUp 完成")

	return

	// 7. 执行 Startup 任务
	fmt.Println("\n开始执行 Startup 任务...")
	taskJob := tasker.PostTask("Startup")
	taskStatus := taskJob.Wait()
	fmt.Printf("任务执行状态: %v\n", taskStatus)
	fmt.Println("任务完成")
}
