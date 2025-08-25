package services

// CreatePopupRequest 创建弹窗请求
type CreatePopupRequest struct {
	ProxyConfigID uint   `json:"proxy_config_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	PopupType     string `json:"popup_type" binding:"required"`
	TriggerType   string `json:"trigger_type" binding:"required"`
	TriggerValue  string `json:"trigger_value"`
	Position      string `json:"position"`
	Style         string `json:"style"`
	Enabled       *bool  `json:"enabled"`
}

// UpdatePopupRequest 更新弹窗请求
type UpdatePopupRequest struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	PopupType    string `json:"popup_type"`
	TriggerType  string `json:"trigger_type"`
	TriggerValue string `json:"trigger_value"`
	Position     string `json:"position"`
	Style        string `json:"style"`
	Enabled      *bool  `json:"enabled"`
}
