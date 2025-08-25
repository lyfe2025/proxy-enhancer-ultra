package services

import (
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubmissionService 提交服务 - 组合各个专门的服务以保持向后兼容性
type SubmissionService struct {
	db     *gorm.DB
	logger logger.Logger

	// 组合的专门服务
	crudService   *SubmissionCRUDService
	queryService  *SubmissionQueryService
	statsService  *SubmissionStatsService
	exportService *SubmissionExportService
}

// NewSubmissionService 创建新的提交服务
func NewSubmissionService(db *gorm.DB, logger logger.Logger) *SubmissionService {
	return &SubmissionService{
		db:     db,
		logger: logger,

		// 初始化专门的服务
		crudService:   NewSubmissionCRUDService(db, logger),
		queryService:  NewSubmissionQueryService(db, logger),
		statsService:  NewSubmissionStatsService(db, logger),
		exportService: NewSubmissionExportService(db, logger),
	}
}

// 为了向后兼容，这些类型现在从submission_types.go重新导出
// 实际定义已移动到专门的类型文件中

// CreateSubmission 创建提交 - 委托给CRUD服务
func (s *SubmissionService) CreateSubmission(req *CreateSubmissionRequest) (*models.Submission, error) {
	return s.crudService.CreateSubmission(req)
}

// GetSubmission 获取提交 - 委托给CRUD服务
func (s *SubmissionService) GetSubmission(id uint) (*models.Submission, error) {
	return s.crudService.GetSubmission(id)
}

// UpdateSubmission 更新提交 - 委托给CRUD服务
func (s *SubmissionService) UpdateSubmission(id uint, req *UpdateSubmissionRequest) error {
	return s.crudService.UpdateSubmission(id, req)
}

// DeleteSubmission 删除提交 - 委托给CRUD服务
func (s *SubmissionService) DeleteSubmission(id uint) error {
	return s.crudService.DeleteSubmission(id)
}

// ListSubmissions 获取提交列表 - 委托给查询服务
func (s *SubmissionService) ListSubmissions(page, pageSize int, popupID *uuid.UUID, startDate, endDate *time.Time) ([]*models.Submission, int64, error) {
	return s.queryService.ListSubmissions(page, pageSize, popupID, startDate, endDate)
}

// GetSubmissionsByPopup 根据弹窗获取提交列表 - 委托给查询服务
func (s *SubmissionService) GetSubmissionsByPopup(popupID uuid.UUID) ([]*models.Submission, error) {
	return s.queryService.GetSubmissionsByPopup(popupID)
}

// GetSubmissionStats 获取提交统计信息 - 委托给统计服务
func (s *SubmissionService) GetSubmissionStats(popupID uuid.UUID) (*SubmissionStats, error) {
	return s.statsService.GetSubmissionStats(popupID)
}

// GetSubmissionsByDateRange 根据日期范围获取提交 - 委托给查询服务
func (s *SubmissionService) GetSubmissionsByDateRange(popupID uuid.UUID, startDate, endDate time.Time) ([]*models.Submission, error) {
	return s.queryService.GetSubmissionsByDateRange(popupID, startDate, endDate)
}

// ExportSubmissions 导出提交数据 - 委托给导出服务
func (s *SubmissionService) ExportSubmissions(popupID uuid.UUID, format string) ([]byte, error) {
	return s.exportService.ExportSubmissions(popupID, format)
}

// DeleteSubmissionsByPopup 根据弹窗删除所有提交 - 委托给CRUD服务
func (s *SubmissionService) DeleteSubmissionsByPopup(popupID uuid.UUID) error {
	return s.crudService.DeleteSubmissionsByPopup(popupID)
}
