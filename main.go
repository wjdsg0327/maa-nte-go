package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"study/one/api"
	"study/one/config"
	"study/one/service"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/*
var embeddedFiles embed.FS

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化 MaaFramework
	if err := service.Service.InitWithConfig(cfg); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 配置CORS中间件 - 支持前后端分离
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// API 路由
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/config", api.GetConfig)

		// Pipeline 管理
		apiGroup.GET("/pipelines", api.ListPipelines)
		apiGroup.GET("/pipelines/:name", api.GetPipeline)
		apiGroup.POST("/pipelines", api.CreatePipeline)
		apiGroup.PUT("/pipelines/:name", api.UpdatePipeline)
		apiGroup.DELETE("/pipelines/:name", api.DeletePipeline)

		// 任务执行
		apiGroup.POST("/tasks", api.ExecuteTask)
		apiGroup.POST("/tasks/run", api.RunTask)
		apiGroup.GET("/tasks/status", api.GetTaskStatus)
		apiGroup.POST("/tasks/stop", api.StopTask)

		// 窗口管理
		apiGroup.GET("/windows", api.GetWindows)
		apiGroup.POST("/windows/connect", api.ConnectWindow)

		// 节点列表
		apiGroup.GET("/nodes", api.GetNodeList)
		apiGroup.GET("/screenshot", api.GetScreenshot)
		apiGroup.POST("/resources/reload", api.ReloadResources)

		// 图片管理
		apiGroup.GET("/images", api.ListImages)
		apiGroup.POST("/images", api.UploadImage)

		// 模型列表
		apiGroup.GET("/models/ocr", api.ListOcrModels)
		apiGroup.GET("/models/detect", api.ListDetectModels)
	}

	// 前端静态文件 - 开发模式使用本地文件，生产模式使用嵌入文件
	if _, err := os.Stat("./web/dist"); err == nil {
		// 开发模式：直接使用本地文件
		r.Static("/assets", "./web/dist/assets")
		r.NoRoute(func(c *gin.Context) {
			c.File("./web/dist/index.html")
		})
		log.Println("开发模式：使用本地 web/dist 文件")
	} else {
		// 生产模式：使用嵌入的文件系统
		distFS, err := fs.Sub(embeddedFiles, "web/dist")
		if err != nil {
			log.Fatalf("无法创建嵌入文件系统: %v", err)
		}

		// 静态资源
		assetsFS, err := fs.Sub(distFS, "assets")
		if err == nil {
			r.StaticFS("/assets", http.FS(assetsFS))
		}

		// 所有其他路由返回 index.html (SPA支持)
		r.NoRoute(func(c *gin.Context) {
			data, err := fs.ReadFile(distFS, "index.html")
			if err != nil {
				c.String(http.StatusNotFound, "Not Found")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		})
		log.Println("生产模式：使用嵌入的前端资源")
	}

	log.Printf("服务器启动: %s", cfg.ServerAddr)
	if err := r.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
