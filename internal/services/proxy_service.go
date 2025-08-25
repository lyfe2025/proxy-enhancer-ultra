package services

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"gorm.io/gorm"
)

// ProxyService 代理服务
type ProxyService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewProxyService 创建新的代理服务
func NewProxyService(db *gorm.DB, logger logger.Logger) *ProxyService {
	return &ProxyService{
		db:     db,
		logger: logger,
	}
}

// CreateProxyConfig 创建代理配置
func (s *ProxyService) CreateProxyConfig(config *models.ProxyConfig) error {
	// 验证配置
	if err := s.validateProxyConfig(config); err != nil {
		return err
	}

	// 检查域名是否已存在
	var existingConfig models.ProxyConfig
	err := s.db.Where("proxy_domain = ?", config.ProxyDomain).First(&existingConfig).Error
	if err == nil {
		return fmt.Errorf("domain %s already exists", config.ProxyDomain)
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 设置默认值
	config.IsActive = true
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	// 创建配置
	if err := s.db.Create(config).Error; err != nil {
		return fmt.Errorf("failed to create proxy config: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"domain": config.ProxyDomain,
		"target_url": config.TargetURL,
	}).Info("Proxy config created successfully")

	return nil
}

// GetProxyConfig 获取代理配置
func (s *ProxyService) GetProxyConfig(id uint) (*models.ProxyConfig, error) {
	var config models.ProxyConfig
	err := s.db.First(&config, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("proxy config not found")
		}
		return nil, err
	}
	return &config, nil
}

// GetProxyConfigByDomain 根据域名获取代理配置
func (s *ProxyService) GetProxyConfigByDomain(domain string) (*models.ProxyConfig, error) {
	var config models.ProxyConfig
	err := s.db.Where("proxy_domain = ? AND is_active = ?", domain, true).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("proxy config not found for domain %s", domain)
		}
		return nil, err
	}
	return &config, nil
}

// UpdateProxyConfig 更新代理配置
func (s *ProxyService) UpdateProxyConfig(id uint, updates *models.ProxyConfig) error {
	// 验证更新数据
	if err := s.validateProxyConfig(updates); err != nil {
		return err
	}

	// 检查配置是否存在
	var existingConfig models.ProxyConfig
	err := s.db.First(&existingConfig, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("proxy config not found")
		}
		return err
	}

	// 如果域名发生变化，检查新域名是否已存在
	if updates.ProxyDomain != existingConfig.ProxyDomain {
		var duplicateConfig models.ProxyConfig
		err := s.db.Where("proxy_domain = ? AND id != ?", updates.ProxyDomain, id).First(&duplicateConfig).Error
		if err == nil {
			return fmt.Errorf("domain %s already exists", updates.ProxyDomain)
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	// 更新时间戳
	updates.UpdatedAt = time.Now()

	// 执行更新
	if err := s.db.Model(&existingConfig).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update proxy config: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"id": id,
		"domain": updates.ProxyDomain,
	}).Info("Proxy config updated successfully")

	return nil
}

// DeleteProxyConfig 删除代理配置
func (s *ProxyService) DeleteProxyConfig(id uint) error {
	// 检查配置是否存在
	var config models.ProxyConfig
	err := s.db.First(&config, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("proxy config not found")
		}
		return err
	}

	// 软删除相关数据
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除相关的弹窗配置
	if err := tx.Where("proxy_config_id = ?", id).Delete(&models.Popup{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete related popups: %w", err)
	}

	// 删除相关的规则
	if err := tx.Where("proxy_config_id = ?", id).Delete(&models.Rule{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete related rules: %w", err)
	}

	// 删除代理配置
	if err := tx.Delete(&config).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete proxy config: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"id": id,
		"domain": config.ProxyDomain,
	}).Info("Proxy config deleted successfully")

	return nil
}

// ListProxyConfigs 获取代理配置列表
func (s *ProxyService) ListProxyConfigs(page, pageSize int, enabled *bool) ([]*models.ProxyConfig, int64, error) {
	var configs []*models.ProxyConfig
	var total int64

	query := s.db.Model(&models.ProxyConfig{})

	// 添加过滤条件
	if enabled != nil {
		query = query.Where("is_active = ?", *enabled)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// ToggleProxyConfig 切换代理配置状态
func (s *ProxyService) ToggleProxyConfig(id uint) error {
	var config models.ProxyConfig
	err := s.db.First(&config, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("proxy config not found")
		}
		return err
	}

	// 切换状态
	config.IsActive = !config.IsActive
	config.UpdatedAt = time.Now()

	if err := s.db.Save(&config).Error; err != nil {
		return fmt.Errorf("failed to toggle proxy config: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"id": id,
		"enabled": config.IsActive,
	}).Info("Proxy config status toggled")

	return nil
}

// GetProxyStats 获取代理统计信息
func (s *ProxyService) GetProxyStats(configID uint, startTime, endTime time.Time) (*ProxyStats, error) {
	stats := &ProxyStats{}

	// 基础查询
	query := s.db.Model(&models.ProxyLog{}).Where("proxy_config_id = ?", configID)

	// 时间范围过滤
	if !startTime.IsZero() {
		query = query.Where("timestamp >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("timestamp <= ?", endTime)
	}

	// 总请求数
	if err := query.Count(&stats.TotalRequests).Error; err != nil {
		return nil, err
	}

	// 成功请求数（2xx状态码）
	if err := query.Where("status_code >= 200 AND status_code < 300").Count(&stats.SuccessRequests).Error; err != nil {
		return nil, err
	}

	// 错误请求数（4xx和5xx状态码）
	if err := query.Where("status_code >= 400").Count(&stats.ErrorRequests).Error; err != nil {
		return nil, err
	}

	// 平均响应时间
	var avgResponseTime float64
	if err := query.Select("AVG(response_time)").Scan(&avgResponseTime).Error; err != nil {
		return nil, err
	}
	stats.AvgResponseTime = int(avgResponseTime)

	// 独立访客数（基于IP）
	if err := query.Distinct("client_ip").Count(&stats.UniqueVisitors).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// validateProxyConfig 验证代理配置
func (s *ProxyService) validateProxyConfig(config *models.ProxyConfig) error {
	if config.ProxyDomain == "" {
		return fmt.Errorf("domain is required")
	}

	if config.TargetURL == "" {
		return fmt.Errorf("target URL is required")
	}

	// 验证目标URL格式
	parsedURL, err := url.Parse(config.TargetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("target URL must include scheme (http or https)")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("target URL must include host")
	}

	// 验证域名格式
	if !s.isValidDomain(config.ProxyDomain) {
		return fmt.Errorf("invalid domain format")
	}

	return nil
}

// isValidDomain 验证域名格式
func (s *ProxyService) isValidDomain(domain string) bool {
	// 简单的域名格式验证
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// 检查是否包含非法字符
	for _, char := range domain {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
			(char >= '0' && char <= '9') || char == '.' || char == '-') {
			return false
		}
	}

	// 检查是否以点或连字符开头/结尾
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") ||
		strings.HasPrefix(domain, "-") || strings.HasSuffix(domain, "-") {
		return false
	}

	return true
}

// ProxyStats 代理统计信息
type ProxyStats struct {
	TotalRequests     int64 `json:"total_requests"`
	SuccessRequests   int64 `json:"success_requests"`
	ErrorRequests     int64 `json:"error_requests"`
	UniqueVisitors    int64 `json:"unique_visitors"`
	AvgResponseTime   int   `json:"avg_response_time"`
}