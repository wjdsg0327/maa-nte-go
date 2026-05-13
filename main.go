package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// 使用 PostPipeline 加载 pipeline 目录
	pipelinePath := filepath.Join(resPath, "pipeline")
	res.PostPipeline(pipelinePath).Wait()

	// 加载图片目录
	imagePath := filepath.Join(resPath, "image")
	res.PostImage(imagePath).Wait()

	// 加载 OCR 模型
	ocrPath := filepath.Join(resPath, "ocr")
	res.PostOcrModel(ocrPath).Wait()

	if !res.Loaded() {
		log.Fatal("资源加载失败")
	}

	// 获取节点列表
	nodes, err := res.GetNodeList()
	if err != nil {
		log.Fatalf("获取节点列表失败: %v", err)
	}
	fmt.Printf("已加载任务: %v\n", nodes)

	if err := tasker.BindResource(res); err != nil {
		log.Fatalf("绑定资源失败: %v", err)
	}

	if !tasker.Initialized() {
		log.Fatal("任务器初始化失败")
	}

	// 7. 选择并执行任务
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== 可用任务 ===")
		for i, node := range nodes {
			fmt.Printf("%d. %s\n", i+1, node)
		}
		fmt.Println("0. 退出")
		fmt.Print("\n请选择任务编号: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			fmt.Println("退出程序")
			return
		}

		// 解析选择
		var choice int
		fmt.Sscanf(input, "%d", &choice)
		if choice < 1 || choice > len(nodes) {
			fmt.Println("无效选择，请重新输入")
			continue
		}

		taskName := nodes[choice-1]
		fmt.Printf("\n开始执行任务: %s\n", taskName)
		taskJob := tasker.PostTask(taskName)
		status := taskJob.Wait()

		// 获取最新节点详情
		nodeDetail, err := tasker.GetLatestNode(taskName)

		if err != nil {
			log.Fatalf("获取节点详情失败: %v", err)
		}

		//ocr处理
		if nodeDetail.Recognition != nil && nodeDetail.Recognition.Results != nil {
			// 所有识别结果
			for i, result := range nodeDetail.Recognition.Results.All {
				if ocrResult, ok := result.AsOCR(); ok {
					fmt.Printf("结果%d: 文本=%s, 置信度=%.2f, 位置=%v\n",
						i, ocrResult.Text, ocrResult.Score, ocrResult.Box)
				}
			}

			// 最佳结果
			if nodeDetail.Recognition.Results.Best != nil {
				if ocrResult, ok := nodeDetail.Recognition.Results.Best.AsOCR(); ok {
					fmt.Printf("最佳结果: 文本=%s, 置信度=%.2f\n", ocrResult.Text, ocrResult.Score)
				}
			}

			// 所有匹配结果
			for i, result := range nodeDetail.Recognition.Results.All {
				if tmResult, ok := result.AsTemplateMatch(); ok {
					fmt.Printf("结果%d: 置信度=%.2f, 位置=%v\n",
						i, tmResult.Score, tmResult.Box)
				}
			}
			// 最佳结果
			if nodeDetail.Recognition.Results.Best != nil {
				if tmResult, ok := nodeDetail.Recognition.Results.Best.AsTemplateMatch(); ok {
					fmt.Printf("最佳匹配: 置信度=%.2f, 位置=%v\n", tmResult.Score, tmResult.Box)
				}
			}

		}

		fmt.Printf("任务执行完成，状态: %v\n", status)
	}
}
