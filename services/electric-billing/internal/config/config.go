package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	MySQL   MySQLConfig   `mapstructure:"mysql"`
	MQTT    MQTTConfig    `mapstructure:"mqtt"`
	Storage StorageConfig `mapstructure:"storage"`
	Log     LogConfig     `mapstructure:"log"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// DSN 生成 MySQL 连接字符串
func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

// MQTTConfig MQTT配置
type MQTTConfig struct {
	Broker   string `mapstructure:"broker"`
	ClientID string `mapstructure:"client_id"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Topic    string `mapstructure:"topic"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	DataDir       string `mapstructure:"data_dir"`       // 数据存储目录
	WALDir        string `mapstructure:"wal_dir"`        // WAL 日志目录
	RetentionDays int    `mapstructure:"retention_days"` // 数据保留天数
	Enabled       bool   `mapstructure:"enabled"`        // 是否启用本地存储
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 环境变量覆盖
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 环境变量绑定
	v.BindEnv("mysql.host", "MYSQL_HOST")
	v.BindEnv("mysql.port", "MYSQL_PORT")
	v.BindEnv("mysql.user", "MYSQL_USER")
	v.BindEnv("mysql.password", "MYSQL_PASSWORD")
	v.BindEnv("mysql.database", "MYSQL_DATABASE")
	v.BindEnv("mqtt.broker", "MQTT_BROKER")
	v.BindEnv("mqtt.username", "MQTT_USERNAME")
	v.BindEnv("mqtt.password", "MQTT_PASSWORD")
	v.BindEnv("storage.data_dir", "DATA_DIR")
	v.BindEnv("storage.wal_dir", "WAL_DIR")
	v.BindEnv("storage.retention_days", "RETENTION_DAYS")
	v.BindEnv("storage.enabled", "STORAGE_ENABLED")

	// 默认值
	v.SetDefault("storage.data_dir", "/app/data")
	v.SetDefault("storage.wal_dir", "/app/data/wal")
	v.SetDefault("storage.retention_days", 30)
	v.SetDefault("storage.enabled", true)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
