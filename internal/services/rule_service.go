package services

import (
	"errors"
	"fmt"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"gorm.io/gorm"
)

// RuleService 规则服务
type RuleService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewRuleService 创建新的规则服务
func NewRuleService(db *gorm.DB, logger logger.Logger) *RuleService {
	return &RuleService{
		db:     db,
		logger: logger,
	}
}

// CreateRuleRequest 创建规则请求
type CreateRuleRequest struct {
	ProxyConfigID uint   `json:"proxy_config_id" binding:"required"`
	RuleType      string `json:"rule_type" binding:"required"`
	Selector      string `json:"selector" binding:"required"`
	Action        string `json:"action" binding:"required"`
	Content       string `json:"content"`
	Position      string `json:"position"`
	Priority      int    `json:"priority"`
	Enabled       *bool  `json:"enabled"`
}

// UpdateRuleRequest 更新规则请求
type UpdateRuleRequest struct {
	RuleType string `json:"rule_type"`
	Selector string `json:"selector"`
	Action   string `json:"action"`
	Content  string `json:"content"`
	Position string `json:"position"`
	Priority *int   `json:"priority"`
	Enabled  *bool  `json:"enabled"`
}

// CreateRule 创建规则
func (s *RuleService) CreateRule(req *CreateRuleRequest) (*models.Rule, error) {
	// 验证代理配置是否存在
	var proxyConfig models.ProxyConfig
	if err := s.db.First(&proxyConfig, req.ProxyConfigID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("proxy config not found")
		}
		return nil, fmt.Errorf("failed to check proxy config: %w", err)
	}

	// 验证规则类型
	validRuleTypes := []string{"inject", "replace", "remove", "modify"}
	if !contains(validRuleTypes, req.RuleType) {
		return nil, errors.New("invalid rule type")
	}

	// 验证动作类型
	validActions := []string{"append", "prepend", "replace", "remove", "modify_attribute", "modify_text"}
	if !contains(validActions, req.Action) {
		return nil, errors.New("invalid action type")
	}

	// 验证位置（如果提供）
	if req.Position != "" {
		validPositions := []string{"before", "after", "inside", "replace"}
		if !contains(validPositions, req.Position) {
			return nil, errors.New("invalid position")
		}
	}

	// 设置默认值
	if req.Enabled == nil {
		enabled := true
		req.Enabled = &enabled
	}

	// 如果没有指定优先级，设置为最低优先级
	if req.Priority == 0 {
		var maxPriority int
		s.db.Model(&models.Rule{}).Where("proxy_config_id = ?", req.ProxyConfigID).Select("COALESCE(MAX(priority), 0)").Scan(&maxPriority)
		req.Priority = maxPriority + 1
	}

	rule := &models.Rule{
		Name:              req.RuleType, // 使用RuleType作为Name
		Description:       req.Content,  // 使用Content作为Description
		TriggerConditions: fmt.Sprintf(`{"selector":"%s","action":"%s","position":"%s"}`, req.Selector, req.Action, req.Position),
		Actions:           fmt.Sprintf(`{"content":"%s"}`, req.Content),
		IsActive:          *req.Enabled,
		Priority:          req.Priority,
	}

	if err := s.db.Create(rule).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"proxy_config_id": req.ProxyConfigID,
			"rule_type":       req.RuleType,
			"error":           err.Error(),
		}).Error("Failed to create rule")
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"rule_id":         rule.ID,
		"proxy_config_id": req.ProxyConfigID,
		"rule_type":       req.RuleType,
	}).Info("Rule created successfully")

	return rule, nil
}

// GetRule 获取规则
func (s *RuleService) GetRule(id uint) (*models.Rule, error) {
	var rule models.Rule
	if err := s.db.Preload("ProxyConfig").First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("rule not found")
		}
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}

	return &rule, nil
}

// UpdateRule 更新规则
func (s *RuleService) UpdateRule(id uint, req *UpdateRuleRequest) error {
	// 检查规则是否存在
	var rule models.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("rule not found")
		}
		return fmt.Errorf("failed to check rule: %w", err)
	}

	// 准备更新数据
	updateData := make(map[string]interface{})

	if req.RuleType != "" {
		validRuleTypes := []string{"inject", "replace", "remove", "modify"}
		if !contains(validRuleTypes, req.RuleType) {
			return errors.New("invalid rule type")
		}
		updateData["rule_type"] = req.RuleType
	}

	if req.Selector != "" {
		updateData["selector"] = req.Selector
	}

	if req.Action != "" {
		validActions := []string{"append", "prepend", "replace", "remove", "modify_attribute", "modify_text"}
		if !contains(validActions, req.Action) {
			return errors.New("invalid action type")
		}
		updateData["action"] = req.Action
	}

	if req.Content != "" {
		updateData["content"] = req.Content
	}

	if req.Position != "" {
		validPositions := []string{"before", "after", "inside", "replace"}
		if !contains(validPositions, req.Position) {
			return errors.New("invalid position")
		}
		updateData["position"] = req.Position
	}

	if req.Priority != nil {
		updateData["priority"] = *req.Priority
	}

	if req.Enabled != nil {
		updateData["enabled"] = *req.Enabled
	}

	updateData["updated_at"] = time.Now()

	if err := s.db.Model(&rule).Updates(updateData).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"rule_id": id,
			"error":   err.Error(),
		}).Error("Failed to update rule")
		return fmt.Errorf("failed to update rule: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"rule_id": id,
	}).Info("Rule updated successfully")

	return nil
}

// DeleteRule 删除规则
func (s *RuleService) DeleteRule(id uint) error {
	// 检查规则是否存在
	var rule models.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("rule not found")
		}
		return fmt.Errorf("failed to check rule: %w", err)
	}

	// 软删除规则
	if err := s.db.Delete(&rule).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"rule_id": id,
			"error":   err.Error(),
		}).Error("Failed to delete rule")
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"rule_id": id,
	}).Info("Rule deleted successfully")

	return nil
}

// ListRules 获取规则列表
func (s *RuleService) ListRules(page, pageSize int, proxyConfigID *uint, ruleType string, enabled *bool) ([]*models.Rule, int64, error) {
	var rules []*models.Rule
	var total int64

	// 构建查询
	query := s.db.Model(&models.Rule{})

	// 添加过滤条件
	if proxyConfigID != nil {
		query = query.Where("proxy_config_id = ?", *proxyConfigID)
	}

	if ruleType != "" {
		query = query.Where("rule_type = ?", ruleType)
	}

	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rules: %w", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("ProxyConfig").Order("priority ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get rules: %w", err)
	}

	return rules, total, nil
}

// ToggleRuleStatus 切换规则状态
func (s *RuleService) ToggleRuleStatus(id uint) error {
	// 检查规则是否存在
	var rule models.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("rule not found")
		}
		return fmt.Errorf("failed to check rule: %w", err)
	}

	// 切换状态
	newStatus := !rule.IsActive
	if err := s.db.Model(&rule).Updates(map[string]interface{}{
		"is_active":  newStatus,
		"updated_at": time.Now(),
	}).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"rule_id": id,
			"error":   err.Error(),
		}).Error("Failed to toggle rule status")
		return fmt.Errorf("failed to toggle rule status: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"rule_id": id,
		"enabled": newStatus,
	}).Info("Rule status toggled successfully")

	return nil
}

// GetRulesByProxyConfig 根据代理配置获取规则
func (s *RuleService) GetRulesByProxyConfig(proxyConfigID uint) ([]*models.Rule, error) {
	var rules []*models.Rule
	if err := s.db.Where("proxy_config_id = ? AND enabled = ?", proxyConfigID, true).Order("priority ASC").Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("failed to get rules by proxy config: %w", err)
	}

	return rules, nil
}

// UpdateRulePriorities 批量更新规则优先级
func (s *RuleService) UpdateRulePriorities(updates []struct {
	ID       uint `json:"id"`
	Priority int  `json:"priority"`
}) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, update := range updates {
		if err := tx.Model(&models.Rule{}).Where("id = ?", update.ID).Updates(map[string]interface{}{
			"priority":   update.Priority,
			"updated_at": time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			s.logger.WithFields(map[string]interface{}{
				"rule_id": update.ID,
				"error":   err.Error(),
			}).Error("Failed to update rule priority")
			return fmt.Errorf("failed to update rule priority: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to commit rule priority updates")
		return fmt.Errorf("failed to commit rule priority updates: %w", err)
	}

	s.logger.Info("Rule priorities updated successfully")
	return nil
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}