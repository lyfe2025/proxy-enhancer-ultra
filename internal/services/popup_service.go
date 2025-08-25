package services

import (
	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"gorm.io/gorm"
)

// PopupService 弹窗服务 - 组合各个专门的服务以保持向后兼容性
type PopupService struct {
	db     *gorm.DB
	logger logger.Logger

	// 组合的专门服务
	crudService  *PopupCRUDService
	statsService *PopupStatsService
}

// NewPopupService 创建新的弹窗服务
func NewPopupService(db *gorm.DB, logger logger.Logger) *PopupService {
	return &PopupService{
		db:     db,
		logger: logger,

		// 初始化专门的服务
		crudService:  NewPopupCRUDService(db, logger),
		statsService: NewPopupStatsService(db, logger),
	}
}

// 为了向后兼容，这些类型现在从popup_types.go重新导出
// 实际定义已移动到专门的类型文件中

// CreatePopup 创建弹窗 - 委托给CRUD服务
func (s *PopupService) CreatePopup(req *CreatePopupRequest) (*models.Popup, error) {
	return s.crudService.CreatePopup(req)
}

// GetPopup 获取弹窗 - 委托给CRUD服务
func (s *PopupService) GetPopup(id uint) (*models.Popup, error) {
	return s.crudService.GetPopup(id)
}

// UpdatePopup 更新弹窗 - 委托给CRUD服务
func (s *PopupService) UpdatePopup(id uint, req *UpdatePopupRequest) error {
	return s.crudService.UpdatePopup(id, req)
}

// DeletePopup 删除弹窗 - 委托给CRUD服务
func (s *PopupService) DeletePopup(id uint) error {
	return s.crudService.DeletePopup(id)
}

// ListPopups 获取弹窗列表 - 委托给CRUD服务
func (s *PopupService) ListPopups(page, pageSize int, proxyConfigID *uint, popupType string, enabled *bool) ([]*models.Popup, int64, error) {
	return s.crudService.ListPopups(page, pageSize, proxyConfigID, popupType, enabled)
}

// TogglePopupStatus 切换弹窗状态 - 委托给CRUD服务
func (s *PopupService) TogglePopupStatus(id uint) error {
	return s.crudService.TogglePopupStatus(id)
}

// GetPopupsByProxyConfig 根据代理配置获取弹窗 - 委托给CRUD服务
func (s *PopupService) GetPopupsByProxyConfig(proxyConfigID uint) ([]*models.Popup, error) {
	return s.crudService.GetPopupsByProxyConfig(proxyConfigID)
}

// GetPopupStats 获取弹窗统计信息 - 委托给统计服务
func (s *PopupService) GetPopupStats(id uint) (map[string]interface{}, error) {
	return s.statsService.GetPopupStats(id)
}
