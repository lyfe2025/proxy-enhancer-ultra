package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Proxy      ProxyConfig      `mapstructure:"proxy"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Security   SecurityConfig   `mapstructure:"security"`
	Cache      CacheConfig      `mapstructure:"cache"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

// PostgresConfig PostgreSQL配置
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	ExpiresIn        time.Duration `mapstructure:"expires_in"`
	RefreshExpiresIn time.Duration `mapstructure:"refresh_expires_in"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Timeout             time.Duration `mapstructure:"timeout"`
	MaxIdleConns        int           `mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost int           `mapstructure:"max_idle_conns_per_host"`
	IdleConnTimeout     time.Duration `mapstructure:"idle_conn_timeout"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	Burst             int `mapstructure:"burst"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Redis RedisConfig `mapstructure:"redis"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	MetricsEnabled  bool   `mapstructure:"metrics_enabled"`
	HealthCheckPath string `mapstructure:"health_check_path"`
	MetricsPath     string `mapstructure:"metrics_path"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// 设置环境变量替换
	viper.SetEnvPrefix("")

	// 服务器配置环境变量绑定
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")

	// 数据库配置环境变量绑定
	viper.BindEnv("database.postgres.host", "DB_HOST")
	viper.BindEnv("database.postgres.port", "DB_PORT")
	viper.BindEnv("database.postgres.user", "DB_USER")
	viper.BindEnv("database.postgres.password", "DB_PASSWORD")
	viper.BindEnv("database.postgres.dbname", "DB_NAME")
	viper.BindEnv("database.postgres.sslmode", "DB_SSL_MODE")

	// JWT配置环境变量绑定
	viper.BindEnv("jwt.secret", "JWT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetAddr 获取服务器地址
func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	// 使用PostgreSQL配置
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Database.Postgres.Host,
		c.Database.Postgres.Port,
		c.Database.Postgres.User,
		c.Database.Postgres.Password,
		c.Database.Postgres.DBName,
		c.Database.Postgres.SSLMode,
		c.Database.Postgres.Timezone,
	)
	fmt.Printf("Using PostgreSQL config: %s\n", dsn)
	return dsn
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}
