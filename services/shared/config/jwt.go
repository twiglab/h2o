package config

// JWTConfig JWT认证配置
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	AccessExpireMin  int    `mapstructure:"access_expire_min"`
	RefreshExpireDay int    `mapstructure:"refresh_expire_day"`
}
