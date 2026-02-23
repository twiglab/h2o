package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"admin-server/internal/config"
	"admin-server/internal/handler"
	"admin-server/internal/repository"
	"admin-server/internal/service"
	"shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var configPath = flag.String("config", "/app/config.yaml", "config file path")

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 连接数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	log.Println("Database connected")

	// 初始化 JWT 中间件
	jwtMiddleware := middleware.NewJWTMiddleware(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpireMin,
		cfg.JWT.RefreshExpireDay,
	)

	// 初始化 repositories
	userRepo := repository.NewUserRepository(db)
	deptRepo := repository.NewDeptRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	meterRepo := repository.NewMeterRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	rateRepo := repository.NewRateRepository(db)
	readingRepo := repository.NewReadingRepository(db)
	statsRepo := repository.NewStatsRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	shopRepo := repository.NewShopRepository(db)

	// 初始化 services
	authService := service.NewAuthService(userRepo, jwtMiddleware)
	userService := service.NewUserService(userRepo)
	deptService := service.NewDeptService(deptRepo)
	permService := service.NewPermissionService(permRepo)
	meterService := service.NewMeterService(meterRepo, readingRepo)
	accountService := service.NewAccountService(accountRepo, meterRepo)
	rateService := service.NewRateService(rateRepo, meterRepo)
	dashboardService := service.NewDashboardService(meterRepo, accountRepo, readingRepo, statsRepo)
	merchantService := service.NewMerchantService(merchantRepo)
	shopService := service.NewShopService(shopRepo)

	// 初始化 handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	deptHandler := handler.NewDeptHandler(deptService)
	permHandler := handler.NewPermissionHandler(permService)
	meterHandler := handler.NewMeterHandler(meterService)
	rateHandler := handler.NewRateHandler(rateService)
	accountHandler := handler.NewAccountHandler(accountService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	systemHandler := handler.NewSystemHandler()
	merchantHandler := handler.NewMerchantHandler(merchantService)
	shopHandler := handler.NewShopHandler(shopService)

	// 创建路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由
	api := r.Group("/api")
	{
		// 认证相关 (无需登录)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(jwtMiddleware.AuthMiddleware())
		{
			// 认证相关
			authorized.POST("/auth/logout", authHandler.Logout)
			authorized.GET("/auth/profile", authHandler.Profile)
			authorized.PUT("/auth/password", authHandler.ChangePassword)

			// 仪表盘
			authorized.GET("/dashboard", dashboardHandler.Dashboard)

			// 报表
			reports := authorized.Group("/reports")
			{
				reports.GET("/consumption", dashboardHandler.ConsumptionReport)
				reports.GET("/revenue", dashboardHandler.RevenueReport)
				reports.GET("/collection", dashboardHandler.CollectionStats)
			}

			// 商户管理
			merchants := authorized.Group("/merchants")
			{
				merchants.GET("", merchantHandler.List)
				merchants.GET("/all", merchantHandler.ListAll)
				merchants.GET("/:id", merchantHandler.Get)
				merchants.POST("", merchantHandler.Create)
				merchants.PUT("/:id", merchantHandler.Update)
				merchants.DELETE("/:id", merchantHandler.Delete)
				merchants.GET("/:id/stats", merchantHandler.GetStats)
			}

			// 店铺管理
			shops := authorized.Group("/shops")
			{
				shops.GET("", shopHandler.List)
				shops.GET("/all", shopHandler.ListAll)
				shops.GET("/:id", shopHandler.Get)
				shops.POST("", shopHandler.Create)
				shops.PUT("/:id", shopHandler.Update)
				shops.DELETE("/:id", shopHandler.Delete)
				shops.GET("/:id/stats", shopHandler.GetStats)
			}

			// 电表管理
			electricMeters := authorized.Group("/electric-meters")
			{
				electricMeters.GET("", meterHandler.ListElectric)
				electricMeters.GET("/:id", meterHandler.GetElectric)
				electricMeters.POST("", meterHandler.CreateElectric)
				electricMeters.PUT("/:id", meterHandler.UpdateElectric)
				electricMeters.DELETE("/:id", meterHandler.DeleteElectric)
				electricMeters.POST("/:id/reading", meterHandler.ManualReadingElectric)
				electricMeters.GET("/:id/readings", meterHandler.GetElectricReadings)
			}

			// 水表管理
			waterMeters := authorized.Group("/water-meters")
			{
				waterMeters.GET("", meterHandler.ListWater)
				waterMeters.GET("/:id", meterHandler.GetWater)
				waterMeters.POST("", meterHandler.CreateWater)
				waterMeters.PUT("/:id", meterHandler.UpdateWater)
				waterMeters.DELETE("/:id", meterHandler.DeleteWater)
				waterMeters.POST("/:id/reading", meterHandler.ManualReadingWater)
				waterMeters.GET("/:id/readings", meterHandler.GetWaterReadings)
			}

			// 手工抄表记录
			authorized.GET("/manual-readings", meterHandler.ListAllReadings)

			// 费率管理
			rates := authorized.Group("/rates")
			{
				rates.GET("", rateHandler.List)
				rates.GET("/:id", rateHandler.Get)
				rates.POST("", rateHandler.Create)
				rates.PUT("/:id", rateHandler.Update)
				rates.DELETE("/:id", rateHandler.Delete)
			}

			// 账户管理
			accounts := authorized.Group("/accounts")
			{
				accounts.GET("", accountHandler.List)
				accounts.GET("/all", accountHandler.ListAll)
				accounts.GET("/arrears", accountHandler.GetArrears)
				accounts.GET("/recharges", accountHandler.ListRecharges)
				accounts.GET("/electric-deductions", accountHandler.ListElectricDeductions)
				accounts.GET("/water-deductions", accountHandler.ListWaterDeductions)
				accounts.GET("/:id", accountHandler.Get)
				accounts.POST("", accountHandler.Create)
				accounts.PUT("/:id", accountHandler.Update)
				accounts.POST("/:id/recharge", accountHandler.Recharge)
				accounts.GET("/:id/recharges", accountHandler.GetRecharges)
				accounts.GET("/:id/electric-deductions", accountHandler.GetElectricDeductions)
				accounts.GET("/:id/water-deductions", accountHandler.GetWaterDeductions)
			}

			// 用户管理
			users := authorized.Group("/users")
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.Get)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
				users.PUT("/:id/password", userHandler.ResetPassword)
				users.PUT("/:id/password/reset", userHandler.ResetPassword) // 兼容前端
			}

			// 部门管理
			depts := authorized.Group("/depts")
			{
				depts.GET("", deptHandler.List)
				depts.GET("/:id", deptHandler.Get)
				depts.POST("", deptHandler.Create)
				depts.PUT("/:id", deptHandler.Update)
				depts.DELETE("/:id", deptHandler.Delete)
			}

			// 权限管理
			permissions := authorized.Group("/permissions")
			{
				permissions.GET("", permHandler.GetTree)
				permissions.GET("/:id", permHandler.Get)
				permissions.POST("", permHandler.Create)
				permissions.PUT("/:id", permHandler.Update)
				permissions.DELETE("/:id", permHandler.Delete)
			}

			// 系统管理
			authorized.GET("/roles", systemHandler.GetRoles)
			authorized.GET("/logs/operation", systemHandler.GetOperationLogs)
		}
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.MySQL.DSN()

	logLevel := logger.Warn
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
