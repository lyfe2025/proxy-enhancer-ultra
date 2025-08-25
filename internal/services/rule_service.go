package services

import (
	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RuleService 规则服务 - 组合各个专门的服务以保持向后兼容性
type RuleService struct {
	db     *gorm.DB
	logger logger.Logger

	// 组合的专门服务
	crudService     *RuleCRUDService
	priorityService *RulePriorityService
}

// NewRuleService 创建新的规则服务
func NewRuleService(db *gorm.DB, logger logger.Logger) *RuleService {
	return &RuleService{
		db:     db,
		logger: logger,

		// 初始化专门的服务
		crudService:     NewRuleCRUDService(db, logger),
		priorityService: NewRulePriorityService(db, logger),
	}
}

// 为了向后兼容，这些类型现在从rule_types.go重新导出
// 实际定义已移动到专门的类型文件中

// CreateRule 创建规则 - 委托给CRUD服务
func (s *RuleService) CreateRule(req *CreateRuleRequest) (*models.Rule, error) {
	return s.crudService.CreateRule(req)
}

// GetRule 获取规则 - 委托给CRUD服务
func (s *RuleService) GetRule(id uuid.UUID) (*models.Rule, error) {
	return s.crudService.GetRule(id)
}

// UpdateRule 更新规则 - 委托给CRUD服务
func (s *RuleService) UpdateRule(id uuid.UUID, req *UpdateRuleRequest) error {
	return s.crudService.UpdateRule(id, req)
}

// DeleteRule 删除规则 - 委托给CRUD服务
func (s *RuleService) DeleteRule(id uuid.UUID) error {
	return s.crudService.DeleteRule(id)
}

// ListRules 获取规则列表 - 委托给CRUD服务
func (s *RuleService) ListRules(page, pageSize int, ruleType string, enabled *bool) ([]*models.Rule, int64, error) {
	return s.crudService.ListRules(page, pageSize, ruleType, enabled)
}

// ToggleRuleStatus 切换规则状态 - 委托给CRUD服务
func (s *RuleService) ToggleRuleStatus(id uuid.UUID) error {
	return s.crudService.ToggleRuleStatus(id)
}

// GetActiveRules 获取所有活跃规则 - 委托给CRUD服务
func (s *RuleService) GetActiveRules() ([]*models.Rule, error) {
	return s.crudService.GetActiveRules()
}

// GetRulesByProxyConfig 根据代理配置获取规则 - 委托给CRUD服务
func (s *RuleService) GetRulesByProxyConfig(proxyConfigID uuid.UUID) ([]*models.Rule, error) {
	return s.crudService.GetRulesByProxyConfig(proxyConfigID)
}

// UpdateRulePriorities 批量更新规则优先级 - 委托给优先级服务
func (s *RuleService) UpdateRulePriorities(updates []RulePriorityUpdate) error {
	return s.priorityService.UpdateRulePriorities(updates)
}
