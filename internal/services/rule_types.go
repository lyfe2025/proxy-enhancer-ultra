package services

import "github.com/google/uuid"

// CreateRuleRequest 创建规则请求
type CreateRuleRequest struct {
	ProxyConfigID uuid.UUID `json:"proxy_config_id" binding:"required"`
	RuleType      string    `json:"rule_type" binding:"required"`
	Selector      string    `json:"selector" binding:"required"`
	Action        string    `json:"action" binding:"required"`
	Content       string    `json:"content"`
	Position      string    `json:"position"`
	Priority      int       `json:"priority"`
	Enabled       *bool     `json:"enabled"`
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

// RulePriorityUpdate 规则优先级更新
type RulePriorityUpdate struct {
	ID       uuid.UUID `json:"id"`
	Priority int       `json:"priority"`
}
