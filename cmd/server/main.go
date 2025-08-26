// @title Proxy Enhancer Ultra API
// @version 1.0
// @description 智能反向代理增强平台API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"proxy-enhancer-ultra/internal/auth"
	"proxy-enhancer-ultra/internal/config"
	"proxy-enhancer-ultra/internal/database"
	"proxy-enhancer-ultra/internal/handlers"
	"proxy-enhancer-ultra/internal/middleware"
	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	_ "proxy-enhancer-ultra/docs" // 导入生成的swagger文档

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 初始化日志器
	logger := logger.NewLogrusLogger()

	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化JWT管理器
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// 初始化数据库
	db, err := database.NewDatabase(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.ProxyConfig{},
		&models.Domain{},
		&models.Rule{},
		&models.Popup{},
		&models.Submission{},
		&models.ProxyLog{},
		&models.SystemMetric{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建默认数据
	err = db.Seed()
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// 初始化服务
	userService := services.NewUserService(db.DB, jwtManager, logger)
	proxyService := services.NewProxyService(db.DB, logger)
	popupService := services.NewPopupService(db.DB, logger)
	ruleService := services.NewRuleService(db.DB, logger)
	submissionService := services.NewSubmissionService(db.DB, logger)
	monitoringService := services.NewMonitoringService(db.DB, logger)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(userService, logger)
	userAdminHandler := handlers.NewUserAdminHandler(userService, logger)
	proxyHandler := handlers.NewProxyHandler(proxyService, logger)
	popupHandler := handlers.NewPopupHandler(popupService, logger)
	ruleHandler := handlers.NewRuleHandler(ruleService, logger)
	submissionHandler := handlers.NewSubmissionHandler(submissionService, logger)
	monitoringHandler := handlers.NewMonitoringHandler(monitoringService, logger)
	profileHandler := handlers.NewProfileHandler(userService, logger)

	// 设置Gin模式（根据环境变量）
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin引擎
	r := gin.New()

	// 添加中间件
	r.Use(gin.LoggerWithWriter(os.Stdout))
	r.Use(gin.Recovery())

	// CORS配置
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.Security.CORS.AllowedOrigins
	corsConfig.AllowMethods = cfg.Security.CORS.AllowedMethods
	corsConfig.AllowHeaders = cfg.Security.CORS.AllowedHeaders
	corsConfig.AllowCredentials = cfg.Security.CORS.AllowCredentials
	r.Use(cors.New(corsConfig))

	// 公共路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", gin.WrapF(authHandler.Login))
		auth.POST("/register", gin.WrapF(authHandler.Register))
		auth.POST("/refresh", gin.WrapF(authHandler.RefreshToken))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
	})

	// Swagger API文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(func(c *gin.Context) {
		// 将gin.Context适配为标准的http.Handler
		handler := middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))
		handler.ServeHTTP(c.Writer, c.Request)
		if c.IsAborted() {
			return
		}
	})
	{
		// 认证相关
		auth := protected.Group("/auth")
		{
			auth.POST("/logout", gin.WrapF(authHandler.Logout))
		}

		// 用户相关
		users := protected.Group("/users")
		{
			users.GET("/profile", gin.WrapF(profileHandler.GetProfile))
			users.PUT("/profile", gin.WrapF(profileHandler.UpdateProfile))
			users.POST("/change-password", gin.WrapF(profileHandler.ChangePassword))
		}

		// 代理配置
		proxy := protected.Group("/proxy")
		{
			proxy.GET("/configs", gin.WrapF(proxyHandler.ListProxyConfigs))
			proxy.POST("/configs", gin.WrapF(proxyHandler.CreateProxyConfig))
			proxy.GET("/configs/:id", gin.WrapF(proxyHandler.GetProxyConfig))
			proxy.PUT("/configs/:id", gin.WrapF(proxyHandler.UpdateProxyConfig))
			proxy.DELETE("/configs/:id", gin.WrapF(proxyHandler.DeleteProxyConfig))
			proxy.POST("/configs/:id/toggle", gin.WrapF(proxyHandler.ToggleProxyConfig))
			proxy.GET("/configs/:id/stats", gin.WrapF(proxyHandler.GetProxyStats))
		}

		// 弹窗管理
		popups := protected.Group("/popups")
		{
			popups.GET("", gin.WrapF(popupHandler.ListPopups))
			popups.POST("", gin.WrapF(popupHandler.CreatePopup))
			popups.GET("/:id", gin.WrapF(popupHandler.GetPopup))
			popups.PUT("/:id", gin.WrapF(popupHandler.UpdatePopup))
			popups.DELETE("/:id", gin.WrapF(popupHandler.DeletePopup))
			popups.POST("/:id/toggle", gin.WrapF(popupHandler.TogglePopupStatus))
			popups.GET("/:id/stats", gin.WrapF(popupHandler.GetPopupStats))
		}

		// 规则管理
		rules := protected.Group("/rules")
		{
			rules.GET("", gin.WrapF(ruleHandler.ListRules))
			rules.POST("", gin.WrapF(ruleHandler.CreateRule))
			rules.GET("/:id", gin.WrapF(ruleHandler.GetRule))
			rules.PUT("/:id", gin.WrapF(ruleHandler.UpdateRule))
			rules.DELETE("/:id", gin.WrapF(ruleHandler.DeleteRule))
			rules.POST("/:id/toggle", gin.WrapF(ruleHandler.ToggleRuleStatus))
			rules.PUT("/priorities", gin.WrapF(ruleHandler.UpdateRulePriorities))
		}

		// 提交管理
		submissions := protected.Group("/submissions")
		{
			submissions.GET("", gin.WrapF(submissionHandler.ListSubmissions))
			submissions.POST("", gin.WrapF(submissionHandler.CreateSubmission))
			submissions.GET("/:id", gin.WrapF(submissionHandler.GetSubmission))
			submissions.PUT("/:id", gin.WrapF(submissionHandler.UpdateSubmission))
			submissions.DELETE("/:id", gin.WrapF(submissionHandler.DeleteSubmission))
			submissions.GET("/export", gin.WrapF(submissionHandler.ExportSubmissions))
			submissions.DELETE("/popup/:popup_id", gin.WrapF(submissionHandler.DeleteSubmissionsByPopup))
		}

		// 系统监控
		monitoring := protected.Group("/monitoring")
		{
			monitoring.GET("/health", gin.WrapF(monitoringHandler.GetHealthCheck))
			monitoring.GET("/dashboard", gin.WrapF(monitoringHandler.GetDashboardData))
			monitoring.GET("/metrics/system", gin.WrapF(monitoringHandler.GetSystemMetrics))
			monitoring.GET("/metrics/proxy", gin.WrapF(monitoringHandler.GetProxyMetrics))
		}

		// 用户管理（管理员）
		admin := users.Group("/admin")
		admin.Use(func(c *gin.Context) {
			// 创建一个适配器来处理AdminMiddleware
			handler := middleware.AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.Next()
			}))
			handler.ServeHTTP(c.Writer, c.Request)
		})
		{
			admin.POST("/", gin.WrapF(userAdminHandler.CreateUser))
			admin.GET("/:id", gin.WrapF(userAdminHandler.GetUser))
			admin.PUT("/:id", gin.WrapF(userAdminHandler.UpdateUser))
			admin.DELETE("/:id", gin.WrapF(userAdminHandler.DeleteUser))
			admin.GET("/", gin.WrapF(userAdminHandler.ListUsers))
		}
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.GetAddr(),
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on %s", cfg.GetAddr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
