package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件到指定结构体
// configPath: 配置文件路径
// cfg: 配置结构体指针
// envBindings: 环境变量绑定映射，key为配置路径，value为环境变量名
func LoadConfig(configPath string, cfg interface{}, envBindings map[string]string) error {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 自动绑定环境变量
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 绑定指定的环境变量
	for configKey, envVar := range envBindings {
		v.BindEnv(configKey, envVar)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}

// DefaultMySQLEnvBindings 返回MySQL的默认环境变量绑定
func DefaultMySQLEnvBindings() map[string]string {
	return map[string]string{
		"mysql.host":     "MYSQL_HOST",
		"mysql.port":     "MYSQL_PORT",
		"mysql.user":     "MYSQL_USER",
		"mysql.password": "MYSQL_PASSWORD",
		"mysql.database": "MYSQL_DATABASE",
	}
}

// DefaultJWTEnvBindings 返回JWT的默认环境变量绑定
func DefaultJWTEnvBindings() map[string]string {
	return map[string]string{
		"jwt.secret": "JWT_SECRET",
	}
}

// MergeEnvBindings 合并多个环境变量绑定映射
func MergeEnvBindings(bindings ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, b := range bindings {
		for k, v := range b {
			result[k] = v
		}
	}
	return result
}
