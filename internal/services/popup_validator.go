package services

import "errors"

// PopupValidator 弹窗验证器
type PopupValidator struct{}

// NewPopupValidator 创建新的弹窗验证器
func NewPopupValidator() *PopupValidator {
	return &PopupValidator{}
}

// ValidatePopupType 验证弹窗类型
func (v *PopupValidator) ValidatePopupType(popupType string) error {
	validPopupTypes := []string{"modal", "toast", "banner", "sidebar", "overlay"}
	if !containsString(validPopupTypes, popupType) {
		return errors.New("invalid popup type")
	}
	return nil
}

// ValidateTriggerType 验证触发类型
func (v *PopupValidator) ValidateTriggerType(triggerType string) error {
	validTriggerTypes := []string{"page_load", "time_delay", "scroll_percentage", "element_click", "exit_intent", "form_submit"}
	if !containsString(validTriggerTypes, triggerType) {
		return errors.New("invalid trigger type")
	}
	return nil
}

// ValidatePosition 验证位置
func (v *PopupValidator) ValidatePosition(position string) error {
	if position == "" {
		return nil // 位置是可选的
	}

	validPositions := []string{"top", "bottom", "left", "right", "center", "top-left", "top-right", "bottom-left", "bottom-right"}
	if !containsString(validPositions, position) {
		return errors.New("invalid position")
	}
	return nil
}

// ValidateCreateRequest 验证创建请求
func (v *PopupValidator) ValidateCreateRequest(req *CreatePopupRequest) error {
	if err := v.ValidatePopupType(req.PopupType); err != nil {
		return err
	}

	if err := v.ValidateTriggerType(req.TriggerType); err != nil {
		return err
	}

	if err := v.ValidatePosition(req.Position); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateRequest 验证更新请求
func (v *PopupValidator) ValidateUpdateRequest(req *UpdatePopupRequest) error {
	if req.PopupType != "" {
		if err := v.ValidatePopupType(req.PopupType); err != nil {
			return err
		}
	}

	if req.TriggerType != "" {
		if err := v.ValidateTriggerType(req.TriggerType); err != nil {
			return err
		}
	}

	if err := v.ValidatePosition(req.Position); err != nil {
		return err
	}

	return nil
}

// containsString 检查切片是否包含指定元素
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
