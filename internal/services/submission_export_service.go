package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubmissionExportService 提交数据导出服务
type SubmissionExportService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewSubmissionExportService 创建新的提交导出服务
func NewSubmissionExportService(db *gorm.DB, logger logger.Logger) *SubmissionExportService {
	return &SubmissionExportService{
		db:     db,
		logger: logger,
	}
}

// ExportSubmissions 导出提交数据
func (s *SubmissionExportService) ExportSubmissions(popupID uuid.UUID, format string) ([]byte, error) {
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
		return s.exportToCSV(submissions)
	default:
		return nil, errors.New("unsupported export format")
	}
}

// exportToCSV 导出为CSV格式
func (s *SubmissionExportService) exportToCSV(submissions []*models.Submission) ([]byte, error) {
	if len(submissions) == 0 {
		return []byte("id,popup_id,form_data,user_agent,user_ip,referrer_url,created_at\n"), nil
	}

	// CSV头部
	csv := "id,popup_id,form_data,user_agent,user_ip,referrer_url,created_at\n"

	// 添加数据行
	for _, submission := range submissions {
		line := fmt.Sprintf("%d,%s,\"%s\",\"%s\",\"%s\",\"%s\",%s\n",
			submission.ID,
			submission.PopupID.String(),
			escapeCSVField(submission.FormData),
			escapeCSVField(submission.UserAgent),
			escapeCSVField(submission.UserIP),
			escapeCSVField(submission.ReferrerURL),
			submission.CreatedAt.Format("2006-01-02 15:04:05"),
		)
		csv += line
	}

	return []byte(csv), nil
}

// escapeCSVField 转义CSV字段中的特殊字符
func escapeCSVField(field string) string {
	// 替换双引号为两个双引号（CSV标准转义方式）
	escaped := ""
	for _, char := range field {
		if char == '"' {
			escaped += "\"\""
		} else {
			escaped += string(char)
		}
	}
	return escaped
}

// ExportSubmissionsByDateRange 按日期范围导出提交数据
func (s *SubmissionExportService) ExportSubmissionsByDateRange(popupID uuid.UUID, startDate, endDate string, format string) ([]byte, error) {
	// 验证弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, popupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, err
	}

	// 获取指定日期范围的提交数据
	var submissions []*models.Submission
	query := s.db.Where("popup_id = ?", popupID)

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	if err := query.Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	// 根据格式导出数据
	switch format {
	case "json":
		return json.MarshalIndent(submissions, "", "  ")
	case "csv":
		return s.exportToCSV(submissions)
	default:
		return nil, errors.New("unsupported export format")
	}
}
