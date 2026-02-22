package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"electric-billing/internal/config"
	"electric-billing/internal/mqtt"
	"electric-billing/internal/repository"
	"electric-billing/internal/service"
	"electric-billing/internal/storage"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var configPath = flag.String("config", "/config.yaml", "config file path")

func main() {
	flag.Parse()

	// 初始化日志
	zapLogger := initLogger("info")
	defer zapLogger.Sync()

	zapLogger.Info("starting electric billing service")

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		zapLogger.Fatal("failed to load config", zap.Error(err))
	}

	// 更新日志级别
	zapLogger = initLogger(cfg.Log.Level)

	zapLogger.Info("config loaded",
		zap.String("mysql_host", cfg.MySQL.Host),
		zap.String("mqtt_broker", cfg.MQTT.Broker))

	// 连接数据库
	db, err := initDB(cfg)
	if err != nil {
		zapLogger.Fatal("failed to connect database", zap.Error(err))
	}
	zapLogger.Info("database connected")

	// 初始化仓库
	meterRepo := repository.NewMeterRepository(db)
	electricRateRepo := repository.NewElectricRateRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	readingRepo := repository.NewReadingRepository(db)

	// 初始化存储
	var wal *storage.WAL
	var fileStore *storage.FileStore

	if cfg.Storage.Enabled {
		var err error
		wal, err = storage.NewWAL(cfg.Storage.WALDir, zapLogger)
		if err != nil {
			zapLogger.Fatal("failed to initialize WAL", zap.Error(err))
		}
		defer wal.Close()
		zapLogger.Info("WAL initialized", zap.String("dir", cfg.Storage.WALDir))

		fileStore, err = storage.NewFileStore(cfg.Storage.DataDir, zapLogger)
		if err != nil {
			zapLogger.Fatal("failed to initialize file store", zap.Error(err))
		}
		zapLogger.Info("file store initialized", zap.String("dir", cfg.Storage.DataDir))

		// 初始化恢复服务并执行 WAL 恢复
		recoveryService := service.NewRecoveryService(
			wal, fileStore, meterRepo, readingRepo, accountRepo, db, zapLogger,
		)

		if err := recoveryService.RecoverFromWAL(); err != nil {
			zapLogger.Error("failed to recover from WAL", zap.Error(err))
		}
	}

	// 初始化计费服务
	billingService := service.NewBillingService(
		db, meterRepo, electricRateRepo, accountRepo, readingRepo, wal, fileStore, zapLogger,
	)

	// 初始化MQTT客户端
	mqttClient := mqtt.NewClient(mqtt.ClientConfig{
		Broker:   cfg.MQTT.Broker,
		ClientID: cfg.MQTT.ClientID,
		Username: cfg.MQTT.Username,
		Password: cfg.MQTT.Password,
	}, zapLogger)

	// 初始化消息处理器
	mqttHandler := mqtt.NewHandler(zapLogger, billingService)
	mqttClient.SetHandler(mqttHandler)

	// 连接MQTT
	if err := mqttClient.Connect(); err != nil {
		zapLogger.Fatal("failed to connect MQTT", zap.Error(err))
	}

	// 订阅主题
	if err := mqttClient.Subscribe(cfg.MQTT.Topic, 1); err != nil {
		zapLogger.Fatal("failed to subscribe topic", zap.Error(err))
	}

	zapLogger.Info("service started, waiting for messages...",
		zap.String("topic", cfg.MQTT.Topic))

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("shutting down...")

	// 清理资源
	mqttClient.Disconnect()

	zapLogger.Info("service stopped")
}

// initLogger 初始化日志
func initLogger(level string) *zap.Logger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

// initDB 初始化数据库连接
func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.MySQL.DSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
