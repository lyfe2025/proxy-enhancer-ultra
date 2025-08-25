package services

import (
	"database/sql"
	"errors"
	"runtime"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"gorm.io/gorm"
)

// MonitoringService 监控服务
type MonitoringService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewMonitoringService 创建新的监控服务
func NewMonitoringService(db *gorm.DB, logger logger.Logger) *MonitoringService {
	return &MonitoringService{
		db:     db,
		logger: logger,
	}
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Goroutines  int     `json:"goroutines"`
	Timestamp   time.Time `json:"timestamp"`
}

// ProxyMetrics 代理指标
type ProxyMetrics struct {
	TotalRequests    int64   `json:"total_requests"`
	SuccessRequests  int64   `json:"success_requests"`
	ErrorRequests    int64   `json:"error_requests"`
	AverageResponse  float64 `json:"average_response_time"`
	UniqueVisitors   int64   `json:"unique_visitors"`
	ActiveProxies    int64   `json:"active_proxies"`
	Timestamp        time.Time `json:"timestamp"`
}

// DatabaseMetrics 数据库指标
type DatabaseMetrics struct {
	Connections     int     `json:"connections"`
	ActiveQueries   int     `json:"active_queries"`
	SlowQueries     int64   `json:"slow_queries"`
	AverageQueryTime float64 `json:"average_query_time"`
	Timestamp       time.Time `json:"timestamp"`
}

// OverallStats 总体统计
type OverallStats struct {
	TotalUsers       int64 `json:"total_users"`
	TotalProxies     int64 `json:"total_proxies"`
	ActiveProxies    int64 `json:"active_proxies"`
	TotalPopups      int64 `json:"total_popups"`
	ActivePopups     int64 `json:"active_popups"`
	TotalSubmissions int64 `json:"total_submissions"`
	TodaySubmissions int64 `json:"today_submissions"`
}

// CollectSystemMetrics 收集系统指标
func (s *MonitoringService) CollectSystemMetrics() (*SystemMetrics, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := &SystemMetrics{
		CPUUsage:    0.0, // 简化实现，实际应该使用系统调用获取CPU使用率
		MemoryUsage: float64(m.Alloc) / 1024 / 1024, // MB
		DiskUsage:   0.0, // 简化实现，实际应该获取磁盘使用率
		Goroutines:  runtime.NumGoroutine(),
		Timestamp:   time.Now(),
	}

	// 存储到数据库
	systemMetric := &models.SystemMetric{
		MetricName:  "overall",
		MetricValue: metrics.MemoryUsage,
		Timestamp:   metrics.Timestamp,
	}

	if err := s.db.Create(systemMetric).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to store system metrics")
		return metrics, err
	}

	return metrics, nil
}

// CollectProxyMetrics 收集代理指标
func (s *MonitoringService) CollectProxyMetrics() (*ProxyMetrics, error) {
	metrics := &ProxyMetrics{
		Timestamp: time.Now(),
	}

	// 获取总请求数（从代理日志表）
	if err := s.db.Model(&models.ProxyLog{}).Count(&metrics.TotalRequests).Error; err != nil {
		return nil, err
	}

	// 获取成功请求数
	if err := s.db.Model(&models.ProxyLog{}).Where("status_code >= 200 AND status_code < 400").Count(&metrics.SuccessRequests).Error; err != nil {
		return nil, err
	}

	// 获取错误请求数
	if err := s.db.Model(&models.ProxyLog{}).Where("status_code >= 400").Count(&metrics.ErrorRequests).Error; err != nil {
		return nil, err
	}

	// 获取平均响应时间
	var avgResponse sql.NullFloat64
	if err := s.db.Model(&models.ProxyLog{}).Select("AVG(response_time)").Scan(&avgResponse).Error; err != nil {
		return nil, err
	}
	if avgResponse.Valid {
		metrics.AverageResponse = avgResponse.Float64
	}

	// 获取独立访客数（基于IP地址）
	if err := s.db.Model(&models.ProxyLog{}).Distinct("user_ip").Count(&metrics.UniqueVisitors).Error; err != nil {
		return nil, err
	}

	// 获取活跃代理数
	if err := s.db.Model(&models.ProxyConfig{}).Where("enabled = ?", true).Count(&metrics.ActiveProxies).Error; err != nil {
		return nil, err
	}

	// 存储到数据库
	proxyMetric := &models.SystemMetric{
		MetricName:  "overall",
		MetricValue: float64(metrics.TotalRequests),
		Timestamp:   metrics.Timestamp,
	}

	if err := s.db.Create(proxyMetric).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to store proxy metrics")
		return metrics, err
	}

	return metrics, nil
}

// GetOverallStats 获取总体统计
func (s *MonitoringService) GetOverallStats() (*OverallStats, error) {
	stats := &OverallStats{}

	// 获取总用户数
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// 获取总代理数
	if err := s.db.Model(&models.ProxyConfig{}).Count(&stats.TotalProxies).Error; err != nil {
		return nil, err
	}

	// 获取活跃代理数
	if err := s.db.Model(&models.ProxyConfig{}).Where("enabled = ?", true).Count(&stats.ActiveProxies).Error; err != nil {
		return nil, err
	}

	// 获取总弹窗数
	if err := s.db.Model(&models.Popup{}).Count(&stats.TotalPopups).Error; err != nil {
		return nil, err
	}

	// 获取活跃弹窗数
	if err := s.db.Model(&models.Popup{}).Where("enabled = ?", true).Count(&stats.ActivePopups).Error; err != nil {
		return nil, err
	}

	// 获取总提交数
	if err := s.db.Model(&models.Submission{}).Count(&stats.TotalSubmissions).Error; err != nil {
		return nil, err
	}

	// 获取今日提交数
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	if err := s.db.Model(&models.Submission{}).Where("submitted_at >= ? AND submitted_at < ?", today, tomorrow).Count(&stats.TodaySubmissions).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetSystemMetricsHistory 获取系统指标历史
func (s *MonitoringService) GetSystemMetricsHistory(hours int) ([]*models.SystemMetric, error) {
	var metrics []*models.SystemMetric
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	if err := s.db.Where("metric_type = ? AND collected_at >= ?", "system", startTime).Order("collected_at DESC").Find(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}

// GetProxyMetricsHistory 获取代理指标历史
func (s *MonitoringService) GetProxyMetricsHistory(hours int) ([]*models.SystemMetric, error) {
	var metrics []*models.SystemMetric
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	if err := s.db.Where("metric_type = ? AND collected_at >= ?", "proxy", startTime).Order("collected_at DESC").Find(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}

// GetProxyStats 获取特定代理的统计信息
func (s *MonitoringService) GetProxyStats(proxyConfigID uint) (map[string]interface{}, error) {
	// 验证代理配置是否存在
	var proxyConfig models.ProxyConfig
	if err := s.db.First(&proxyConfig, proxyConfigID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("proxy config not found")
		}
		return nil, err
	}

	stats := make(map[string]interface{})

	// 获取总请求数
	var totalRequests int64
	if err := s.db.Model(&models.ProxyLog{}).Where("proxy_config_id = ?", proxyConfigID).Count(&totalRequests).Error; err != nil {
		return nil, err
	}
	stats["total_requests"] = totalRequests

	// 获取成功请求数
	var successRequests int64
	if err := s.db.Model(&models.ProxyLog{}).Where("proxy_config_id = ? AND status_code >= 200 AND status_code < 400", proxyConfigID).Count(&successRequests).Error; err != nil {
		return nil, err
	}
	stats["success_requests"] = successRequests

	// 获取错误请求数
	var errorRequests int64
	if err := s.db.Model(&models.ProxyLog{}).Where("proxy_config_id = ? AND status_code >= 400", proxyConfigID).Count(&errorRequests).Error; err != nil {
		return nil, err
	}
	stats["error_requests"] = errorRequests

	// 获取独立访客数
	var uniqueVisitors int64
	if err := s.db.Model(&models.ProxyLog{}).Where("proxy_config_id = ?", proxyConfigID).Distinct("user_ip").Count(&uniqueVisitors).Error; err != nil {
		return nil, err
	}
	stats["unique_visitors"] = uniqueVisitors

	// 获取平均响应时间
	var avgResponse sql.NullFloat64
	if err := s.db.Model(&models.ProxyLog{}).Where("proxy_config_id = ?", proxyConfigID).Select("AVG(response_time)").Scan(&avgResponse).Error; err != nil {
		return nil, err
	}
	if avgResponse.Valid {
		stats["average_response_time"] = avgResponse.Float64
	} else {
		stats["average_response_time"] = 0.0
	}

	// 获取最近访问时间
	var lastAccess time.Time
	if err := s.db.Model(&models.ProxyLog{}).Where("proxy_config_id = ?", proxyConfigID).Select("MAX(created_at)").Scan(&lastAccess).Error; err != nil {
		return nil, err
	}
	stats["last_access"] = lastAccess

	return stats, nil
}

// CleanupOldMetrics 清理旧的指标数据
func (s *MonitoringService) CleanupOldMetrics(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)

	// 删除旧的系统指标
	if err := s.db.Where("collected_at < ?", cutoffTime).Delete(&models.SystemMetric{}).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"cutoff_time": cutoffTime,
			"error":       err.Error(),
		}).Error("Failed to cleanup old system metrics")
		return err
	}

	// 删除旧的代理日志
	if err := s.db.Where("created_at < ?", cutoffTime).Delete(&models.ProxyLog{}).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"cutoff_time": cutoffTime,
			"error":       err.Error(),
		}).Error("Failed to cleanup old proxy logs")
		return err
	}

	s.logger.WithFields(map[string]interface{}{
		"cutoff_time": cutoffTime,
		"days":        days,
	}).Info("Old metrics cleaned up successfully")

	return nil
}

// StartMetricsCollection 启动指标收集
func (s *MonitoringService) StartMetricsCollection(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			// 收集系统指标
			if _, err := s.CollectSystemMetrics(); err != nil {
				s.logger.WithFields(map[string]interface{}{
					"error": err.Error(),
				}).Error("Failed to collect system metrics")
			}

			// 收集代理指标
			if _, err := s.CollectProxyMetrics(); err != nil {
				s.logger.WithFields(map[string]interface{}{
					"error": err.Error(),
				}).Error("Failed to collect proxy metrics")
			}
		}
	}()

	s.logger.WithFields(map[string]interface{}{
		"interval": interval,
	}).Info("Metrics collection started")
}