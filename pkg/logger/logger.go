package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger 日志接口
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// LogrusLogger logrus实现的日志器
type LogrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// NewLogrusLogger 创建新的logrus日志器
func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	
	// 设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	
	// 设置日志级别
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	
	// 设置输出
	logger.SetOutput(os.Stdout)
	
	return &LogrusLogger{
		logger: logger,
		entry:  logrus.NewEntry(logger),
	}
}

// Debug 调试日志
func (l *LogrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Info 信息日志
func (l *LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Warn 警告日志
func (l *LogrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Error 错误日志
func (l *LogrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Fatal 致命错误日志
func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Panic panic日志
func (l *LogrusLogger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

// Debugf 格式化调试日志
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Infof 格式化信息日志
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warnf 格式化警告日志
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Errorf 格式化错误日志
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// Panicf 格式化panic日志
func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

// WithField 添加字段
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

// WithFields 添加多个字段
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(logrus.Fields(fields)),
	}
}

// WithError 添加错误字段
func (l *LogrusLogger) WithError(err error) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithError(err),
	}
}

// StandardLogger 标准日志器（用于兼容标准库）
type StandardLogger struct {
	logger Logger
}

// NewStandardLogger 创建标准日志器
func NewStandardLogger(logger Logger) *StandardLogger {
	return &StandardLogger{logger: logger}
}

// Write 实现io.Writer接口
func (s *StandardLogger) Write(p []byte) (n int, err error) {
	s.logger.Info(string(p))
	return len(p), nil
}

// Printf 格式化打印
func (s *StandardLogger) Printf(format string, v ...interface{}) {
	s.logger.Infof(format, v...)
}

// Print 打印
func (s *StandardLogger) Print(v ...interface{}) {
	s.logger.Info(v...)
}

// Println 打印行
func (s *StandardLogger) Println(v ...interface{}) {
	s.logger.Info(v...)
}

// SetupGlobalLogger 设置全局日志器
func SetupGlobalLogger(logger Logger) {
	stdLogger := NewStandardLogger(logger)
	log.SetOutput(stdLogger)
	log.SetFlags(0) // 禁用标准库的时间戳，使用我们自己的格式
}

// GetLogLevel 获取日志级别
func GetLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return "info"
	}
	return level
}

// SetLogLevel 设置日志级别
func SetLogLevel(level string) error {
	switch level {
	case "debug", "info", "warn", "error":
		os.Setenv("LOG_LEVEL", level)
		return nil
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}
}

// LogMiddleware HTTP日志中间件
type LogMiddleware struct {
	logger Logger
}

// NewLogMiddleware 创建日志中间件
func NewLogMiddleware(logger Logger) *LogMiddleware {
	return &LogMiddleware{logger: logger}
}

// LogRequest 记录HTTP请求
func (m *LogMiddleware) LogRequest(method, path, clientIP string, statusCode int, duration time.Duration, userAgent string) {
	m.logger.WithFields(map[string]interface{}{
		"method":      method,
		"path":        path,
		"client_ip":   clientIP,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
		"user_agent":  userAgent,
	}).Info("HTTP Request")
}

// LogError 记录HTTP错误
func (m *LogMiddleware) LogError(method, path, clientIP string, err error, userAgent string) {
	m.logger.WithFields(map[string]interface{}{
		"method":     method,
		"path":       path,
		"client_ip":  clientIP,
		"user_agent": userAgent,
	}).WithError(err).Error("HTTP Error")
}