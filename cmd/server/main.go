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
	"proxy-enhancer-ultra/internal/proxy"
	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	loggerInstance := logger.NewLogrusLogger()
	logger.SetupGlobalLogger(loggerInstance)

	// 数据库连接（允许失败，以便先测试代理功能）
	var db *database.Database
	var dbConnected bool

	db, err = database.NewDatabase(cfg, loggerInstance)
	if err != nil {
		loggerInstance.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Warn("Failed to connect to database, continuing without database features")
		dbConnected = false
	} else {
		dbConnected = true
		defer db.Close()

		// 执行数据库迁移
		if err := db.AutoMigrate(
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
		); err != nil {
			loggerInstance.WithFields(map[string]interface{}{
				"error": err.Error(),
			}).Warn("Failed to migrate database, continuing without database features")
			dbConnected = false
		}
	}

	// 初始化JWT管理器
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// 初始化服务（根据数据库连接状态）
	var userService *services.UserService
	var proxyService *services.ProxyService
	var ruleService *services.RuleService
	var popupService *services.PopupService
	var submissionService *services.SubmissionService
	var monitoringService *services.MonitoringService
	var proxyServer *proxy.ProxyServer

	if dbConnected && db != nil {
		userService = services.NewUserService(db.DB, jwtManager, loggerInstance)
		proxyService = services.NewProxyService(db.DB, loggerInstance)
		ruleService = services.NewRuleService(db.DB, loggerInstance)
		popupService = services.NewPopupService(db.DB, loggerInstance)
		submissionService = services.NewSubmissionService(db.DB, loggerInstance)
		monitoringService = services.NewMonitoringService(db.DB, loggerInstance)
		proxyServer = proxy.NewProxyServer(db.DB, loggerInstance, cfg)
	} else {
		// 创建无数据库版本的代理服务器
		proxyServer = proxy.NewProxyServer(nil, loggerInstance, cfg)
		loggerInstance.Info("Running in proxy-only mode without database features")
	}

	// 初始化处理器（根据数据库连接状态）
	var authHandler *handlers.AuthHandler
	var proxyHandler *handlers.ProxyHandler
	var ruleHandler *handlers.RuleHandler
	var popupHandler *handlers.PopupHandler
	var submissionHandler *handlers.SubmissionHandler
	var monitoringHandler *handlers.MonitoringHandler

	if dbConnected && db != nil {
		authHandler = handlers.NewAuthHandler(userService, loggerInstance)
		proxyHandler = handlers.NewProxyHandler(proxyService, loggerInstance)
		ruleHandler = handlers.NewRuleHandler(ruleService, loggerInstance)
		popupHandler = handlers.NewPopupHandler(popupService, loggerInstance)
		submissionHandler = handlers.NewSubmissionHandler(submissionService, loggerInstance)
		monitoringHandler = handlers.NewMonitoringHandler(monitoringService, loggerInstance)

		// 启动监控服务
		if cfg.Monitoring.MetricsEnabled {
			go monitoringService.StartMetricsCollection(5 * time.Minute)
		}
	}

	// 设置路由
	router := mux.NewRouter()

	// 只有在数据库连接时才设置API路由
	if dbConnected && db != nil {
		// 认证路由
		auth := router.PathPrefix("/api/auth").Subrouter()
		auth.HandleFunc("/login", authHandler.Login).Methods("POST")
		auth.HandleFunc("/register", authHandler.Register).Methods("POST")
		auth.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

		// 受保护的API路由
		api := router.PathPrefix("/api").Subrouter()
		api.Use(middleware.AuthMiddleware(jwtManager))

		// 用户相关路由
		api.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
		api.HandleFunc("/profile", authHandler.UpdateProfile).Methods("PUT")
		api.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")

		// 管理后台API（需要管理员权限）
		admin := api.PathPrefix("/admin").Subrouter()
		admin.Use(middleware.AdminMiddleware)

		// 用户管理
		admin.HandleFunc("/users", authHandler.CreateUser).Methods("POST")
		admin.HandleFunc("/users", authHandler.ListUsers).Methods("GET")
		admin.HandleFunc("/users/{id}", authHandler.GetUser).Methods("GET")
		admin.HandleFunc("/users/{id}", authHandler.UpdateUser).Methods("PUT")
		admin.HandleFunc("/users/{id}", authHandler.DeleteUser).Methods("DELETE")

		// 代理配置管理
		admin.HandleFunc("/proxies", proxyHandler.CreateProxyConfig).Methods("POST")
		admin.HandleFunc("/proxies", proxyHandler.ListProxyConfigs).Methods("GET")
		admin.HandleFunc("/proxies/{id}", proxyHandler.GetProxyConfig).Methods("GET")
		admin.HandleFunc("/proxies/{id}", proxyHandler.UpdateProxyConfig).Methods("PUT")
		admin.HandleFunc("/proxies/{id}", proxyHandler.DeleteProxyConfig).Methods("DELETE")
		admin.HandleFunc("/proxies/{id}/toggle", proxyHandler.ToggleProxyConfig).Methods("POST")
		admin.HandleFunc("/proxies/{id}/stats", proxyHandler.GetProxyStats).Methods("GET")

		// 规则管理
		admin.HandleFunc("/rules", ruleHandler.CreateRule).Methods("POST")
		admin.HandleFunc("/rules", ruleHandler.ListRules).Methods("GET")
		admin.HandleFunc("/rules/{id}", ruleHandler.GetRule).Methods("GET")
		admin.HandleFunc("/rules/{id}", ruleHandler.UpdateRule).Methods("PUT")
		admin.HandleFunc("/rules/{id}", ruleHandler.DeleteRule).Methods("DELETE")
		admin.HandleFunc("/rules/{id}/toggle", ruleHandler.ToggleRuleStatus).Methods("POST")
		admin.HandleFunc("/rules/proxy/{proxy_config_id}", ruleHandler.GetRulesByProxyConfig).Methods("GET")
		admin.HandleFunc("/rules/batch/priority", ruleHandler.UpdateRulePriorities).Methods("PUT")

		// 弹窗管理
		admin.HandleFunc("/popups", popupHandler.CreatePopup).Methods("POST")
		admin.HandleFunc("/popups", popupHandler.ListPopups).Methods("GET")
		admin.HandleFunc("/popups/{id}", popupHandler.GetPopup).Methods("GET")
		admin.HandleFunc("/popups/{id}", popupHandler.UpdatePopup).Methods("PUT")
		admin.HandleFunc("/popups/{id}", popupHandler.DeletePopup).Methods("DELETE")
		admin.HandleFunc("/popups/{id}/toggle", popupHandler.TogglePopupStatus).Methods("POST")
		admin.HandleFunc("/popups/proxy/{proxy_config_id}", popupHandler.GetPopupsByProxyConfig).Methods("GET")
		admin.HandleFunc("/popups/{id}/stats", popupHandler.GetPopupStats).Methods("GET")

		// 提交数据管理
		admin.HandleFunc("/submissions", submissionHandler.ListSubmissions).Methods("GET")
		admin.HandleFunc("/submissions/{id}", submissionHandler.GetSubmission).Methods("GET")
		admin.HandleFunc("/submissions/{id}", submissionHandler.UpdateSubmission).Methods("PUT")
		admin.HandleFunc("/submissions/{id}", submissionHandler.DeleteSubmission).Methods("DELETE")
		admin.HandleFunc("/submissions/popup/{popup_id}", submissionHandler.GetSubmissionsByPopup).Methods("GET")
		admin.HandleFunc("/submissions/stats", submissionHandler.GetSubmissionStats).Methods("GET")
		admin.HandleFunc("/submissions/export", submissionHandler.ExportSubmissions).Methods("GET")
		admin.HandleFunc("/submissions/date-range", submissionHandler.GetSubmissionsByDateRange).Methods("GET")

		// 监控相关路由
		admin.HandleFunc("/monitoring/system", monitoringHandler.GetSystemMetrics).Methods("GET")
		admin.HandleFunc("/monitoring/proxy", monitoringHandler.GetProxyMetrics).Methods("GET")
		admin.HandleFunc("/monitoring/stats", monitoringHandler.GetOverallStats).Methods("GET")
		admin.HandleFunc("/monitoring/system/history", monitoringHandler.GetSystemMetricsHistory).Methods("GET")
		admin.HandleFunc("/monitoring/proxy/history", monitoringHandler.GetProxyMetricsHistory).Methods("GET")
		admin.HandleFunc("/monitoring/proxy/{proxy_config_id}/stats", monitoringHandler.GetProxyStats).Methods("GET")
		admin.HandleFunc("/monitoring/cleanup", monitoringHandler.CleanupOldMetrics).Methods("POST")
		admin.HandleFunc("/monitoring/dashboard", monitoringHandler.GetDashboardData).Methods("GET")

		// 公共API路由（不需要认证）
		public := router.PathPrefix("/api/public").Subrouter()
		public.HandleFunc("/submissions", submissionHandler.CreateSubmission).Methods("POST")

		// 健康检查和监控端点
		router.HandleFunc(cfg.Monitoring.HealthCheckPath, monitoringHandler.GetHealthCheck).Methods("GET")
		router.HandleFunc(cfg.Monitoring.MetricsPath, monitoringHandler.GetSystemMetrics).Methods("GET")
	}

	// 静态文件服务（前端）
	staticDir := "./web/dist"
	if _, err := os.Stat(staticDir); err == nil {
		fs := http.FileServer(http.Dir(staticDir))
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	}

	// 代理处理（处理所有其他请求）
	router.PathPrefix("/").HandlerFunc(proxyServer.ServeHTTP)

	// 设置CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.Security.CORS.AllowedOrigins,
		AllowedMethods:   cfg.Security.CORS.AllowedMethods,
		AllowedHeaders:   cfg.Security.CORS.AllowedHeaders,
		AllowCredentials: cfg.Security.CORS.AllowCredentials,
	})

	// 应用中间件
	handler := middleware.LoggingMiddleware(loggerInstance)(middleware.RecoveryMiddleware(loggerInstance)(c.Handler(router)))

	// 如果启用了限流，应用限流中间件
	handler = middleware.RateLimitMiddleware(cfg.Security.RateLimit.RequestsPerMinute)(handler)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         cfg.GetAddr(),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 启动服务器
	go func() {
		loggerInstance.WithFields(map[string]interface{}{
			"address": cfg.GetAddr(),
		}).Info("Starting HTTP server")

		loggerInstance.Info("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loggerInstance.WithFields(map[string]interface{}{
				"error": err.Error(),
			}).Fatal("Failed to start HTTP server")
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	loggerInstance.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		loggerInstance.Fatal("Server forced to shutdown:", err)
	}

	loggerInstance.Info("Server exited")
}
