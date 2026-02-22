package db

import (
	"shared/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Options 数据库连接选项
type Options struct {
	LogLevel     logger.LogLevel
	MaxIdleConns int
	MaxOpenConns int
}

// DefaultOptions 返回默认连接选项
func DefaultOptions() Options {
	return Options{
		LogLevel:     logger.Warn,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}
}

// NewMySQL 创建 MySQL 数据库连接
func NewMySQL(cfg *config.MySQLConfig, opts ...Options) (*gorm.DB, error) {
	opt := DefaultOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	dsn := cfg.DSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(opt.LogLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(opt.MaxIdleConns)
	sqlDB.SetMaxOpenConns(opt.MaxOpenConns)

	return db, nil
}

// Close 关闭数据库连接
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
