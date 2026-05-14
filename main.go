package main

import (
	"log"

	"study/one/api"
	"study/one/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 MaaFramework
	if err := service.Service.Init(); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 静态文件
	r.Static("/static", "./web")
	r.GET("/", api.ServeIndex)

	// API 路由
	apiGroup := r.Group("/api")
	{
		// Pipeline 管理
		apiGroup.GET("/pipelines", api.ListPipelines)
		apiGroup.GET("/pipelines/:name", api.GetPipeline)
		apiGroup.POST("/pipelines", api.CreatePipeline)
		apiGroup.PUT("/pipelines/:name", api.UpdatePipeline)
		apiGroup.DELETE("/pipelines/:name", api.DeletePipeline)

		// 任务执行
		apiGroup.POST("/tasks", api.ExecuteTask)
		apiGroup.GET("/tasks/status", api.GetTaskStatus)
		apiGroup.POST("/tasks/stop", api.StopTask)

		// 窗口管理
		apiGroup.GET("/windows", api.GetWindows)
		apiGroup.POST("/windows/connect", api.ConnectWindow)

		// 节点列表
		apiGroup.GET("/nodes", api.GetNodeList)

		// 图片管理
		apiGroup.GET("/images", api.ListImages)
		apiGroup.POST("/images", api.UploadImage)

		// 模型列表
		apiGroup.GET("/models/ocr", api.ListOcrModels)
		apiGroup.GET("/models/detect", api.ListDetectModels)
	}

	log.Println("服务器启动: http://localhost:8080")
	r.Run(":8080")
}
