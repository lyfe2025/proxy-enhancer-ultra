package services

import (
	"errors"
	"fmt"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RuleCRUDService 规则CRUD操作服务
type RuleCRUDService struct {
	db        *gorm.DB
	logger    logger.Logger
	validator *RuleValidator
}

// NewRuleCRUDService 创建新的规则CRUD服务
func NewRuleCRUDService(db *gorm.DB, logger logger.Logger) *RuleCRUDService {
	return &RuleCRUDService{
		db:        db,
		logger:    logger,
		validator: NewRuleValidator(),
	}
}

// CreateRule 创建规则
func (s *RuleCRUDService) CreateRule(req *CreateRuleRequest) (*models.Rule, error) {
	// 验证请求数据
	if err := s.validator.ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	// 设置默认值
	if req.Enabled == nil {
		enabled := true
		req.Enabled = &enabled
	}

	// 如果没有指定优先级，设置为最低优先级
	if req.Priority == 0 {
		var maxPriority int
		s.db.Model(&models.Rule{}).Select("COALESCE(MAX(priority), 0)").Scan(&maxPriority)
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
			"rule_type": req.RuleType,
			"error":     err.Error(),
		}).Error("Failed to create rule")
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"rule_id":   rule.ID,
		"rule_type": req.RuleType,
	}).Info("Rule created successfully")

	return rule, nil
}

// GetRule 获取规则
func (s *RuleCRUDService) GetRule(id uuid.UUID) (*models.Rule, error) {
	var rule models.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("rule not found")
		}
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}

	return &rule, nil
}

// UpdateRule 更新规则
func (s *RuleCRUDService) UpdateRule(id uuid.UUID, req *UpdateRuleRequest) error {
	// 检查规则是否存在
	var rule models.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("rule not found")
		}
		return fmt.Errorf("failed to check rule: %w", err)
	}

	// 验证请求数据
	if err := s.validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	// 准备更新数据
	updateData := make(map[string]interface{})

	if req.RuleType != "" {
		updateData["rule_type"] = req.RuleType
	}
	if req.Selector != "" {
		updateData["selector"] = req.Selector
	}
	if req.Action != "" {
		updateData["action"] = req.Action
	}
	if req.Content != "" {
		updateData["content"] = req.Content
	}
	if req.Position != "" {
		updateData["position"] = req.Position
	}
	if req.Priority != nil {
		updateData["priority"] = *req.Priority
	}
	if req.Enabled != nil {
		updateData["is_active"] = *req.Enabled
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
func (s *RuleCRUDService) DeleteRule(id uuid.UUID) error {
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
func (s *RuleCRUDService) ListRules(page, pageSize int, ruleType string, enabled *bool) ([]*models.Rule, int64, error) {
	var rules []*models.Rule
	var total int64

	// 构建查询
	query := s.db.Model(&models.Rule{})

	// 添加过滤条件
	if ruleType != "" {
		query = query.Where("name = ?", ruleType) // 使用name字段代替rule_type
	}
	if enabled != nil {
		query = query.Where("is_active = ?", *enabled)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rules: %w", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Order("priority ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get rules: %w", err)
	}

	return rules, total, nil
}

// ToggleRuleStatus 切换规则状态
func (s *RuleCRUDService) ToggleRuleStatus(id uuid.UUID) error {
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

// GetActiveRules 获取所有活跃规则
func (s *RuleCRUDService) GetActiveRules() ([]*models.Rule, error) {
	var rules []*models.Rule
	if err := s.db.Where("is_active = ?", true).Order("priority ASC").Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("failed to get active rules: %w", err)
	}

	return rules, nil
}

// GetRulesByProxyConfig 根据代理配置获取规则
func (s *RuleCRUDService) GetRulesByProxyConfig(proxyConfigID uuid.UUID) ([]*models.Rule, error) {
	var rules []*models.Rule
	if err := s.db.Where("proxy_config_id = ? AND is_active = ?", proxyConfigID, true).Order("priority ASC").Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("failed to get rules by proxy config: %w", err)
	}

	return rules, nil
}
