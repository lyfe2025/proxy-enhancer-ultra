package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"proxy-enhancer-ultra/internal/config"
	"proxy-enhancer-ultra/pkg/logger"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Database 数据库连接管理器
type Database struct {
	DB     *gorm.DB
	SQLDB  *sql.DB
	logger logger.Logger
	config *config.Config
}

// getEnvOrDefault 获取环境变量或返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewDatabase 创建新的数据库连接
func NewDatabase(cfg *config.Config, log logger.Logger) (*Database, error) {
	// 使用配置的GetDSN方法获取数据库连接字符串
	dsn := cfg.GetDSN()

	// 配置GORM日志
	logLevel := gormLogger.Silent
	if os.Getenv("DB_DEBUG") == "true" {
		logLevel = gormLogger.Info
	}

	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层SQL DB连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100) // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connected successfully")

	return &Database{
		DB:     db,
		SQLDB:  sqlDB,
		logger: log,
		config: cfg,
	}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	if d.SQLDB != nil {
		return d.SQLDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移数据库表结构
func (d *Database) AutoMigrate(models ...interface{}) error {
	d.logger.Info("Starting database migration...")
	
	for _, model := range models {
		if err := d.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}
	
	d.logger.Info("Database migration completed successfully")
	return nil
}

// Health 检查数据库健康状态
func (d *Database) Health() error {
	return d.SQLDB.Ping()
}