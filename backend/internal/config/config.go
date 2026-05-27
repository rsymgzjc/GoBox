package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	SMTP     SMTPConfig     `mapstructure:"smtp"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}

type AuthConfig struct {
	JWTSecret                      string `mapstructure:"jwt_secret"`
	TokenDuration                  int    `mapstructure:"token_duration_hours"`
	VerificationCodeTTLMinutes     int    `mapstructure:"verification_code_ttl_minutes"`
	VerificationCodeCooldownSecond int    `mapstructure:"verification_code_cooldown_seconds"`
}

type SMTPConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

func (s SMTPConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type ZapField = zap.Field

func ZapString(key, value string) ZapField {
	return zap.String(key, value)
}

func ErrorField(err error) ZapField {
	return zap.Error(err)
}

func Load() (*Config, *zap.Logger, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./backend")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, nil, err
	}

	loggerCfg := zap.NewProductionConfig()
	if cfg.App.Env == "development" {
		loggerCfg = zap.NewDevelopmentConfig()
	}
	logger, err := loggerCfg.Build()
	if err != nil {
		return nil, nil, err
	}

	return &cfg, logger, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "GoBox")
	v.SetDefault("app.env", "development")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.dsn", "gobox:gobox123@tcp(127.0.0.1:3306)/gobox?charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("redis.addr", "redis:6379")
	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.token_duration_hours", 72)
	v.SetDefault("auth.verification_code_ttl_minutes", 10)
	v.SetDefault("auth.verification_code_cooldown_seconds", 60)
	v.SetDefault("smtp.enabled", false)
	v.SetDefault("smtp.host", "smtp.example.com")
	v.SetDefault("smtp.port", 587)
	v.SetDefault("smtp.username", "")
	v.SetDefault("smtp.password", "")
	v.SetDefault("smtp.from", "GoBox <noreply@example.com>")
}
