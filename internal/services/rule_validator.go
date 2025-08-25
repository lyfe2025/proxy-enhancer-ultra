package services

import "errors"

// RuleValidator 规则验证器
type RuleValidator struct{}

// NewRuleValidator 创建新的规则验证器
func NewRuleValidator() *RuleValidator {
	return &RuleValidator{}
}

// ValidateRuleType 验证规则类型
func (v *RuleValidator) ValidateRuleType(ruleType string) error {
	validRuleTypes := []string{"inject", "replace", "remove", "modify"}
	if !containsString(validRuleTypes, ruleType) {
		return errors.New("invalid rule type")
	}
	return nil
}

// ValidateAction 验证动作类型
func (v *RuleValidator) ValidateAction(action string) error {
	validActions := []string{"append", "prepend", "replace", "remove", "modify_attribute", "modify_text"}
	if !containsString(validActions, action) {
		return errors.New("invalid action type")
	}
	return nil
}

// ValidatePosition 验证位置
func (v *RuleValidator) ValidatePosition(position string) error {
	if position == "" {
		return nil // 位置是可选的
	}

	validPositions := []string{"before", "after", "inside", "replace"}
	if !containsString(validPositions, position) {
		return errors.New("invalid position")
	}
	return nil
}

// ValidateCreateRequest 验证创建请求
func (v *RuleValidator) ValidateCreateRequest(req *CreateRuleRequest) error {
	if err := v.ValidateRuleType(req.RuleType); err != nil {
		return err
	}

	if err := v.ValidateAction(req.Action); err != nil {
		return err
	}

	if err := v.ValidatePosition(req.Position); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateRequest 验证更新请求
func (v *RuleValidator) ValidateUpdateRequest(req *UpdateRuleRequest) error {
	if req.RuleType != "" {
		if err := v.ValidateRuleType(req.RuleType); err != nil {
			return err
		}
	}

	if req.Action != "" {
		if err := v.ValidateAction(req.Action); err != nil {
			return err
		}
	}

	if err := v.ValidatePosition(req.Position); err != nil {
		return err
	}

	return nil
}
