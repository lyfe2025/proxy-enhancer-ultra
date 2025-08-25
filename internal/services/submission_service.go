package services

import (
	"encoding/json"
	"errors"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubmissionService 提交服务
type SubmissionService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewSubmissionService 创建新的提交服务
func NewSubmissionService(db *gorm.DB, logger logger.Logger) *SubmissionService {
	return &SubmissionService{
		db:     db,
		logger: logger,
	}
}

// CreateSubmissionRequest 创建提交请求
type CreateSubmissionRequest struct {
	PopupID    uuid.UUID              `json:"popup_id" binding:"required"`
	FormData   map[string]interface{} `json:"form_data" binding:"required"`
	UserAgent  string                 `json:"user_agent"`
	IPAddress  string                 `json:"ip_address"`
	Referrer   string                 `json:"referrer"`
	SubmittedAt *time.Time            `json:"submitted_at"`
}

// UpdateSubmissionRequest 更新提交请求
type UpdateSubmissionRequest struct {
	FormData   map[string]interface{} `json:"form_data"`
	UserAgent  string                 `json:"user_agent"`
	IPAddress  string                 `json:"ip_address"`
	Referrer   string                 `json:"referrer"`
	SubmittedAt *time.Time            `json:"submitted_at"`
}

// SubmissionStats 提交统计信息
type SubmissionStats struct {
	TotalSubmissions int64     `json:"total_submissions"`
	TodaySubmissions int64     `json:"today_submissions"`
	WeekSubmissions  int64     `json:"week_submissions"`
	MonthSubmissions int64     `json:"month_submissions"`
	LastSubmission   *time.Time `json:"last_submission"`
	PopupID          uuid.UUID `json:"popup_id"`
}

// CreateSubmission 创建提交
func (s *SubmissionService) CreateSubmission(req *CreateSubmissionRequest) (*models.Submission, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, req.PopupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	// 序列化表单数据
	formDataJSON, err := json.Marshal(req.FormData)
	if err != nil {
		return nil, errors.New("invalid form data")
	}

	// 创建提交记录
	submission := &models.Submission{
		PopupID:     req.PopupID, // PopupID是uuid.UUID类型
		FormData:    string(formDataJSON),
		UserAgent:   req.UserAgent,
		UserIP:     req.IPAddress,  // 使用UserIP字段
		ReferrerURL: req.Referrer,   // 使用ReferrerURL字段
	}

	if err := s.db.Create(submission).Error; err != nil {
		return nil, err
	}

	s.logger.WithFields(map[string]interface{}{
		"submission_id": submission.ID,
		"popup_id":      submission.PopupID,
		"user_ip":       submission.UserIP,
	}).Info("Submission created successfully")

	return submission, nil
}

// GetSubmission 获取提交
func (s *SubmissionService) GetSubmission(id uint) (*models.Submission, error) {
	var submission models.Submission
	if err := s.db.Preload("Popup").First(&submission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}

	return &submission, nil
}

// UpdateSubmission 更新提交
func (s *SubmissionService) UpdateSubmission(id uint, req *UpdateSubmissionRequest) error {
	// 检查提交是否存在
	var submission models.Submission
	if err := s.db.First(&submission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("submission not found")
		}
		return err
	}

	// 准备更新数据
	updates := make(map[string]interface{})

	if req.FormData != nil {
		formDataJSON, err := json.Marshal(req.FormData)
		if err != nil {
			return errors.New("invalid form data")
		}
		updates["form_data"] = string(formDataJSON)
	}

	if req.UserAgent != "" {
		updates["user_agent"] = req.UserAgent
	}

	if req.IPAddress != "" {
		updates["user_ip"] = req.IPAddress  // 使用user_ip字段
	}

	if req.Referrer != "" {
		updates["referrer_url"] = req.Referrer  // 使用referrer_url字段
	}

	if req.SubmittedAt != nil {
		updates["submitted_at"] = *req.SubmittedAt
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	updates["updated_at"] = time.Now()

	if err := s.db.Model(&submission).Updates(updates).Error; err != nil {
		return err
	}

	s.logger.WithFields(map[string]interface{}{
		"submission_id": id,
		"updates":       updates,
	}).Info("Submission updated successfully")

	return nil
}

// DeleteSubmission 删除提交
func (s *SubmissionService) DeleteSubmission(id uint) error {
	// 检查提交是否存在
	var submission models.Submission
	if err := s.db.First(&submission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("submission not found")
		}
		return err
	}

	// 软删除
	if err := s.db.Delete(&submission).Error; err != nil {
		return err
	}

	s.logger.WithFields(map[string]interface{}{
		"submission_id": id,
	}).Info("Submission deleted successfully")

	return nil
}

// ListSubmissions 获取提交列表
func (s *SubmissionService) ListSubmissions(page, pageSize int, popupID *uuid.UUID, startDate, endDate *time.Time) ([]*models.Submission, int64, error) {
	var submissions []*models.Submission
	var total int64

	// 构建查询
	query := s.db.Model(&models.Submission{})

	// 添加过滤条件
	if popupID != nil {
		query = query.Where("popup_id = ?", *popupID)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Popup").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

// GetSubmissionsByPopup 根据弹窗获取提交列表
func (s *SubmissionService) GetSubmissionsByPopup(popupID uuid.UUID) ([]*models.Submission, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	var submissions []*models.Submission
	if err := s.db.Where("popup_id = ?", popupID).Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

// GetSubmissionStats 获取提交统计信息
func (s *SubmissionService) GetSubmissionStats(popupID uuid.UUID) (*SubmissionStats, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	stats := &SubmissionStats{
		PopupID: popupID,
	}

	// 获取总提交数
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ?", popupID).Count(&stats.TotalSubmissions).Error; err != nil {
		return nil, err
	}

	// 获取今日提交数
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ? AND created_at < ?", popupID, today, tomorrow).Count(&stats.TodaySubmissions).Error; err != nil {
		return nil, err
	}

	// 获取本周提交数
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	weekEnd := weekStart.Add(7 * 24 * time.Hour)
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ? AND created_at < ?", popupID, weekStart, weekEnd).Count(&stats.WeekSubmissions).Error; err != nil {
		return nil, err
	}

	// 获取本月提交数
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	if err := s.db.Model(&models.Submission{}).Where("popup_id = ? AND created_at >= ? AND created_at < ?", popupID, monthStart, monthEnd).Count(&stats.MonthSubmissions).Error; err != nil {
		return nil, err
	}

	// 获取最近提交时间
	var lastSubmission models.Submission
	if err := s.db.Where("popup_id = ?", popupID).Order("created_at DESC").First(&lastSubmission).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// 没有提交记录，lastSubmission 保持为 nil
	} else {
		stats.LastSubmission = &lastSubmission.CreatedAt
	}

	return stats, nil
}

// GetSubmissionsByDateRange 根据日期范围获取提交
func (s *SubmissionService) GetSubmissionsByDateRange(popupID uuid.UUID, startDate, endDate time.Time) ([]*models.Submission, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	var submissions []*models.Submission
	if err := s.db.Where("popup_id = ? AND created_at >= ? AND created_at <= ?", popupID, startDate, endDate).Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

// ExportSubmissions 导出提交数据
func (s *SubmissionService) ExportSubmissions(popupID uuid.UUID, format string) ([]byte, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	// 获取所有提交数据
	var submissions []*models.Submission
	if err := s.db.Where("popup_id = ?", popupID).Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	// 根据格式导出数据
	switch format {
	case "json":
		return json.MarshalIndent(submissions, "", "  ")
	case "csv":
		// 这里可以实现CSV导出逻辑
		return nil, errors.New("CSV export not implemented yet")
	default:
		return nil, errors.New("unsupported export format")
	}
}

// DeleteSubmissionsByPopup 根据弹窗删除所有提交（软删除）
func (s *SubmissionService) DeleteSubmissionsByPopup(popupID uuid.UUID) error {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("popup not found")
		}
		return err
	}

	// 软删除所有相关提交
	if err := s.db.Where("popup_id = ?", popupID).Delete(&models.Submission{}).Error; err != nil {
		return err
	}

	s.logger.WithFields(map[string]interface{}{
		"popup_id": popupID,
	}).Info("All submissions for popup deleted successfully")

	return nil
}