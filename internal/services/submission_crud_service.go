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

// SubmissionCRUDService 提交记录CRUD操作服务
type SubmissionCRUDService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewSubmissionCRUDService 创建新的提交CRUD服务
func NewSubmissionCRUDService(db *gorm.DB, logger logger.Logger) *SubmissionCRUDService {
	return &SubmissionCRUDService{
		db:     db,
		logger: logger,
	}
}

// CreateSubmission 创建提交
func (s *SubmissionCRUDService) CreateSubmission(req *CreateSubmissionRequest) (*models.Submission, error) {
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
		UserIP:      req.IPAddress, // 使用UserIP字段
		ReferrerURL: req.Referrer,  // 使用ReferrerURL字段
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
func (s *SubmissionCRUDService) GetSubmission(id uint) (*models.Submission, error) {
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
func (s *SubmissionCRUDService) UpdateSubmission(id uint, req *UpdateSubmissionRequest) error {
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
		updates["user_ip"] = req.IPAddress // 使用user_ip字段
	}

	if req.Referrer != "" {
		updates["referrer_url"] = req.Referrer // 使用referrer_url字段
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
func (s *SubmissionCRUDService) DeleteSubmission(id uint) error {
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

// DeleteSubmissionsByPopup 根据弹窗删除所有提交（软删除）
func (s *SubmissionCRUDService) DeleteSubmissionsByPopup(popupID uuid.UUID) error {
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
