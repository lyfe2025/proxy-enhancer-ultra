package services

import (
	"time"

	"github.com/google/uuid"
)

// CreateSubmissionRequest 创建提交请求
type CreateSubmissionRequest struct {
	PopupID     uuid.UUID              `json:"popup_id" binding:"required"`
	FormData    map[string]interface{} `json:"form_data" binding:"required"`
	UserAgent   string                 `json:"user_agent"`
	IPAddress   string                 `json:"ip_address"`
	Referrer    string                 `json:"referrer"`
	SubmittedAt *time.Time             `json:"submitted_at"`
}

// UpdateSubmissionRequest 更新提交请求
type UpdateSubmissionRequest struct {
	FormData    map[string]interface{} `json:"form_data"`
	UserAgent   string                 `json:"user_agent"`
	IPAddress   string                 `json:"ip_address"`
	Referrer    string                 `json:"referrer"`
	SubmittedAt *time.Time             `json:"submitted_at"`
}

// SubmissionStats 提交统计信息
type SubmissionStats struct {
	TotalSubmissions int64      `json:"total_submissions"`
	TodaySubmissions int64      `json:"today_submissions"`
	WeekSubmissions  int64      `json:"week_submissions"`
	MonthSubmissions int64      `json:"month_submissions"`
	LastSubmission   *time.Time `json:"last_submission"`
	PopupID          uuid.UUID  `json:"popup_id"`
}
