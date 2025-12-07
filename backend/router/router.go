package router

import (
	"nas-navigation/config"
	"nas-navigation/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API路由组
	api := r.Group("/api")
	{
		// 分类
		categoryHandler := handlers.NewCategoryHandler()
		api.GET("/categories", categoryHandler.GetAll)
		api.GET("/categories/:id", categoryHandler.GetByID)
		api.POST("/categories", categoryHandler.Create)
		api.PUT("/categories/:id", categoryHandler.Update)
		api.DELETE("/categories/:id", categoryHandler.Delete)

		// 链接
		linkHandler := handlers.NewLinkHandler()
		api.GET("/links", linkHandler.GetAll)
		api.GET("/links/search", linkHandler.Search)
		api.GET("/links/:id", linkHandler.GetByID)
		api.POST("/links", linkHandler.Create)
		api.PUT("/links/:id", linkHandler.Update)
		api.DELETE("/links/:id", linkHandler.Delete)

		// 脚本
		scriptHandler := handlers.NewScriptHandler()
		api.GET("/scripts", scriptHandler.GetAll)
		api.GET("/scripts/:id", scriptHandler.GetByID)
		api.POST("/scripts", scriptHandler.Create)
		api.PUT("/scripts/:id", scriptHandler.Update)
		api.DELETE("/scripts/:id", scriptHandler.Delete)
		api.POST("/scripts/:id/execute", scriptHandler.Execute)

		// 设置
		settingHandler := handlers.NewSettingHandler()
		api.GET("/settings", settingHandler.GetAll)
		api.PUT("/settings", settingHandler.Update)

		// 系统操作
		systemHandler := handlers.NewSystemHandler()
		api.GET("/system/info", systemHandler.GetInfo)
		api.POST("/system/restart", systemHandler.Restart)
		api.POST("/system/restart-all", systemHandler.RestartAllServices)
		api.POST("/system/restart-link/:id", systemHandler.RestartLink)
		api.POST("/system/quick/:action", systemHandler.ExecuteQuickAction)
	}

	// 静态文件服务
	r.Static("/assets", cfg.StaticDir+"/assets")
	r.StaticFile("/", cfg.StaticDir+"/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File(cfg.StaticDir + "/index.html")
	})

	return r
}
