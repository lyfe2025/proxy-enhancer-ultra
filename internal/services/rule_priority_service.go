package services

import (
	"fmt"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RulePriorityService 规则优先级管理服务
type RulePriorityService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewRulePriorityService 创建新的规则优先级服务
func NewRulePriorityService(db *gorm.DB, logger logger.Logger) *RulePriorityService {
	return &RulePriorityService{
		db:     db,
		logger: logger,
	}
}

// UpdateRulePriorities 批量更新规则优先级
func (s *RulePriorityService) UpdateRulePriorities(updates []RulePriorityUpdate) error {
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

// ReorderRules 重新排序规则优先级
func (s *RulePriorityService) ReorderRules(ruleIDs []uuid.UUID) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 按照提供的顺序更新优先级
	for i, ruleID := range ruleIDs {
		priority := i + 1
		if err := tx.Model(&models.Rule{}).Where("id = ?", ruleID).Updates(map[string]interface{}{
			"priority":   priority,
			"updated_at": time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			s.logger.WithFields(map[string]interface{}{
				"rule_id": ruleID,
				"error":   err.Error(),
			}).Error("Failed to reorder rule")
			return fmt.Errorf("failed to reorder rule: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to commit rule reordering")
		return fmt.Errorf("failed to commit rule reordering: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"rule_count": len(ruleIDs),
	}).Info("Rules reordered successfully")

	return nil
}

// GetNextPriority 获取下一个可用的优先级
func (s *RulePriorityService) GetNextPriority() (int, error) {
	var maxPriority int
	if err := s.db.Model(&models.Rule{}).Select("COALESCE(MAX(priority), 0)").Scan(&maxPriority).Error; err != nil {
		return 0, fmt.Errorf("failed to get max priority: %w", err)
	}

	return maxPriority + 1, nil
}

// MovePriorityUp 将规则优先级上移
func (s *RulePriorityService) MovePriorityUp(ruleID uuid.UUID) error {
	// 获取当前规则
	var currentRule models.Rule
	if err := s.db.First(&currentRule, ruleID).Error; err != nil {
		return fmt.Errorf("failed to get current rule: %w", err)
	}

	// 查找上一个优先级的规则
	var prevRule models.Rule
	if err := s.db.Where("priority < ?", currentRule.Priority).
		Order("priority DESC").First(&prevRule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 已经是最高优先级
		}
		return fmt.Errorf("failed to find previous rule: %w", err)
	}

	// 交换优先级
	return s.swapPriorities(currentRule.ID, prevRule.ID, currentRule.Priority, prevRule.Priority)
}

// MovePriorityDown 将规则优先级下移
func (s *RulePriorityService) MovePriorityDown(ruleID uuid.UUID) error {
	// 获取当前规则
	var currentRule models.Rule
	if err := s.db.First(&currentRule, ruleID).Error; err != nil {
		return fmt.Errorf("failed to get current rule: %w", err)
	}

	// 查找下一个优先级的规则
	var nextRule models.Rule
	if err := s.db.Where("priority > ?", currentRule.Priority).
		Order("priority ASC").First(&nextRule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 已经是最低优先级
		}
		return fmt.Errorf("failed to find next rule: %w", err)
	}

	// 交换优先级
	return s.swapPriorities(currentRule.ID, nextRule.ID, currentRule.Priority, nextRule.Priority)
}

// swapPriorities 交换两个规则的优先级
func (s *RulePriorityService) swapPriorities(rule1ID, rule2ID uuid.UUID, priority1, priority2 int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新第一个规则
	if err := tx.Model(&models.Rule{}).Where("id = ?", rule1ID).Updates(map[string]interface{}{
		"priority":   priority2,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update rule1 priority: %w", err)
	}

	// 更新第二个规则
	if err := tx.Model(&models.Rule{}).Where("id = ?", rule2ID).Updates(map[string]interface{}{
		"priority":   priority1,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update rule2 priority: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit priority swap: %w", err)
	}

	return nil
}
