package services

import (
	"errors"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubmissionQueryService 提交记录查询服务
type SubmissionQueryService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewSubmissionQueryService 创建新的提交查询服务
func NewSubmissionQueryService(db *gorm.DB, logger logger.Logger) *SubmissionQueryService {
	return &SubmissionQueryService{
		db:     db,
		logger: logger,
	}
}

// ListSubmissions 获取提交列表
func (s *SubmissionQueryService) ListSubmissions(page, pageSize int, popupID *uuid.UUID, startDate, endDate *time.Time) ([]*models.Submission, int64, error) {
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
func (s *SubmissionQueryService) GetSubmissionsByPopup(popupID uuid.UUID) ([]*models.Submission, error) {
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

// GetSubmissionsByDateRange 根据日期范围获取提交
func (s *SubmissionQueryService) GetSubmissionsByDateRange(popupID uuid.UUID, startDate, endDate time.Time) ([]*models.Submission, error) {
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

// SearchSubmissions 搜索提交记录
func (s *SubmissionQueryService) SearchSubmissions(keyword string, popupID *uuid.UUID, limit int) ([]*models.Submission, error) {
	var submissions []*models.Submission

	query := s.db.Model(&models.Submission{})

	// 添加弹窗过滤
	if popupID != nil {
		query = query.Where("popup_id = ?", *popupID)
	}

	// 在表单数据、用户代理、IP地址中搜索关键词
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("form_data LIKE ? OR user_agent LIKE ? OR user_ip LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Preload("Popup").Order("created_at DESC").Limit(limit).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}
