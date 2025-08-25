package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Server ServerConfig `json:"server"`
	
	// 数据库配置
	Database DatabaseConfig `json:"database"`
	
	// JWT配置
	JWT JWTConfig `json:"jwt"`
	
	// 日志配置
	Log LogConfig `json:"log"`
	
	// 代理配置
	Proxy ProxyConfig `json:"proxy"`
	
	// 监控配置
	Monitoring MonitoringConfig `json:"monitoring"`
	
	// CORS配置
	CORS CORSConfig `json:"cors"`
	
	// 限流配置
	RateLimit RateLimitConfig `json:"rate_limit"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	TLSCertFile  string        `json:"tls_cert_file"`
	TLSKeyFile   string        `json:"tls_key_file"`
	EnableTLS    bool          `json:"enable_tls"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	DBName          string        `json:"db_name"`
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	LogLevel        string        `json:"log_level"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string        `json:"secret"`
	Expiration       time.Duration `json:"expiration"`
	RefreshExpiration time.Duration `json:"refresh_expiration"`
	Issuer           string        `json:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Timeout         time.Duration `json:"timeout"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	IdleConnTimeout time.Duration `json:"idle_conn_timeout"`
	UserAgent       string        `json:"user_agent"`
	FollowRedirects bool          `json:"follow_redirects"`
	MaxRedirects    int           `json:"max_redirects"`
	BufferSize      int           `json:"buffer_size"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled           bool          `json:"enabled"`
	CollectInterval   time.Duration `json:"collect_interval"`
	RetentionDays     int           `json:"retention_days"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`
	MetricsPath       string        `json:"metrics_path"`
	HealthCheckPath   string        `json:"health_check_path"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled          bool     `json:"enabled"`
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	ExposedHeaders   []string `json:"exposed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled    bool          `json:"enabled"`
	RPS        int           `json:"rps"`
	Burst      int           `json:"burst"`
	WindowSize time.Duration `json:"window_size"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
}

// Load 加载配置
func Load() (*Config, error) {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		// .env文件不存在不是错误，继续使用环境变量
		fmt.Printf("Warning: .env file not found, using environment variables\n")
	}

	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", "30s"),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", "60s"),
			TLSCertFile:  getEnv("TLS_CERT_FILE", ""),
			TLSKeyFile:   getEnv("TLS_KEY_FILE", ""),
			EnableTLS:    getEnvAsBool("ENABLE_TLS", false),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "proxy_enhancer_ultra"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "5m"),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", "5m"),
			LogLevel:        getEnv("DB_LOG_LEVEL", "warn"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "your-secret-key"),
			Expiration:       getEnvAsDuration("JWT_EXPIRATION", "24h"),
			RefreshExpiration: getEnvAsDuration("JWT_REFRESH_EXPIRATION", "168h"), // 7 days
			Issuer:           getEnv("JWT_ISSUER", "proxy-enhancer-ultra"),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			Output:     getEnv("LOG_OUTPUT", "stdout"),
			FilePath:   getEnv("LOG_FILE_PATH", "logs/app.log"),
			MaxSize:    getEnvAsInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 3),
			MaxAge:     getEnvAsInt("LOG_MAX_AGE", 28),
			Compress:   getEnvAsBool("LOG_COMPRESS", true),
		},
		Proxy: ProxyConfig{
			Timeout:         getEnvAsDuration("PROXY_TIMEOUT", "30s"),
			MaxIdleConns:    getEnvAsInt("PROXY_MAX_IDLE_CONNS", 100),
			IdleConnTimeout: getEnvAsDuration("PROXY_IDLE_CONN_TIMEOUT", "90s"),
			UserAgent:       getEnv("PROXY_USER_AGENT", "ProxyEnhancerUltra/1.0"),
			FollowRedirects: getEnvAsBool("PROXY_FOLLOW_REDIRECTS", true),
			MaxRedirects:    getEnvAsInt("PROXY_MAX_REDIRECTS", 10),
			BufferSize:      getEnvAsInt("PROXY_BUFFER_SIZE", 32768), // 32KB
		},
		Monitoring: MonitoringConfig{
			Enabled:           getEnvAsBool("MONITORING_ENABLED", true),
			CollectInterval:   getEnvAsDuration("MONITORING_COLLECT_INTERVAL", "1m"),
			RetentionDays:     getEnvAsInt("MONITORING_RETENTION_DAYS", 30),
			CleanupInterval:   getEnvAsDuration("MONITORING_CLEANUP_INTERVAL", "24h"),
			MetricsPath:       getEnv("MONITORING_METRICS_PATH", "/metrics"),
			HealthCheckPath:   getEnv("MONITORING_HEALTH_CHECK_PATH", "/health"),
		},
		CORS: CORSConfig{
			Enabled:          getEnvAsBool("CORS_ENABLED", true),
			AllowedOrigins:   getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods:   getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"*"}),
			ExposedHeaders:   getEnvAsSlice("CORS_EXPOSED_HEADERS", []string{}),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 86400), // 24 hours
		},
		RateLimit: RateLimitConfig{
			Enabled:         getEnvAsBool("RATE_LIMIT_ENABLED", true),
			RPS:             getEnvAsInt("RATE_LIMIT_RPS", 100),
			Burst:           getEnvAsInt("RATE_LIMIT_BURST", 200),
			WindowSize:      getEnvAsDuration("RATE_LIMIT_WINDOW_SIZE", "1m"),
			CleanupInterval: getEnvAsDuration("RATE_LIMIT_CLEANUP_INTERVAL", "5m"),
		},
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// 验证JWT配置
	if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key" {
		return fmt.Errorf("JWT secret must be set and not use default value")
	}

	// 验证数据库配置
	if c.Database.Host == "" {
		return fmt.Errorf("database host must be set")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	// 验证TLS配置
	if c.Server.EnableTLS {
		if c.Server.TLSCertFile == "" || c.Server.TLSKeyFile == "" {
			return fmt.Errorf("TLS cert and key files must be set when TLS is enabled")
		}
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	// 使用PostgreSQL连接
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetServerAddress 获取服务器地址
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsProduction 是否为生产环境
func (c *Config) IsProduction() bool {
	return strings.ToLower(getEnv("ENV", "development")) == "production"
}

// IsDevelopment 是否为开发环境
func (c *Config) IsDevelopment() bool {
	return !c.IsProduction()
}

// 辅助函数

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsDuration 获取环境变量并转换为时间间隔
func getEnvAsDuration(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	// 如果解析失败，尝试解析默认值
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	// 如果都失败，返回1分钟
	return time.Minute
}

// getEnvAsSlice 获取环境变量并转换为字符串切片
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}