package services

import (
	"errors"
	"fmt"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"gorm.io/gorm"
)

// PopupCRUDService 弹窗CRUD操作服务
type PopupCRUDService struct {
	db        *gorm.DB
	logger    logger.Logger
	validator *PopupValidator
}

// NewPopupCRUDService 创建新的弹窗CRUD服务
func NewPopupCRUDService(db *gorm.DB, logger logger.Logger) *PopupCRUDService {
	return &PopupCRUDService{
		db:        db,
		logger:    logger,
		validator: NewPopupValidator(),
	}
}

// CreatePopup 创建弹窗
func (s *PopupCRUDService) CreatePopup(req *CreatePopupRequest) (*models.Popup, error) {
	// 验证代理配置是否存在
	var proxyConfig models.ProxyConfig
	if err := s.db.First(&proxyConfig, req.ProxyConfigID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("proxy config not found")
		}
		return nil, fmt.Errorf("failed to check proxy config: %w", err)
	}

	// 验证请求数据
	if err := s.validator.ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	// 设置默认值
	position := req.Position
	if position == "" {
		position = "center" // 默认居中
	}

	if req.Enabled == nil {
		enabled := true
		req.Enabled = &enabled
	}

	popup := &models.Popup{
		Title:       req.Title,
		Content:     req.Content,
		StyleConfig: req.Style,
		FormConfig:  "{}", // 默认空的表单配置
		IsActive:    *req.Enabled,
	}

	if err := s.db.Create(popup).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"proxy_config_id": req.ProxyConfigID,
			"title":           req.Title,
			"popup_type":      req.PopupType,
			"error":           err.Error(),
		}).Error("Failed to create popup")
		return nil, fmt.Errorf("failed to create popup: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"popup_id":        popup.ID,
		"proxy_config_id": req.ProxyConfigID,
		"title":           req.Title,
	}).Info("Popup created successfully")

	return popup, nil
}

// GetPopup 获取弹窗
func (s *PopupCRUDService) GetPopup(id uint) (*models.Popup, error) {
	var popup models.Popup
	if err := s.db.Preload("ProxyConfig").First(&popup, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("popup not found")
		}
		return nil, fmt.Errorf("failed to get popup: %w", err)
	}

	return &popup, nil
}

// UpdatePopup 更新弹窗
func (s *PopupCRUDService) UpdatePopup(id uint, req *UpdatePopupRequest) error {
	// 检查弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("popup not found")
		}
		return fmt.Errorf("failed to check popup: %w", err)
	}

	// 验证请求数据
	if err := s.validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	// 准备更新数据
	updateData := make(map[string]interface{})

	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Content != "" {
		updateData["content"] = req.Content
	}
	if req.PopupType != "" {
		updateData["popup_type"] = req.PopupType
	}
	if req.TriggerType != "" {
		updateData["trigger_type"] = req.TriggerType
	}
	if req.TriggerValue != "" {
		updateData["trigger_value"] = req.TriggerValue
	}
	if req.Position != "" {
		updateData["position"] = req.Position
	}
	if req.Style != "" {
		updateData["style"] = req.Style
	}
	if req.Enabled != nil {
		updateData["is_active"] = *req.Enabled
	}

	updateData["updated_at"] = time.Now()

	if err := s.db.Model(&popup).Updates(updateData).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to update popup")
		return fmt.Errorf("failed to update popup: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"popup_id": id,
	}).Info("Popup updated successfully")

	return nil
}

// DeletePopup 删除弹窗
func (s *PopupCRUDService) DeletePopup(id uint) error {
	// 检查弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("popup not found")
		}
		return fmt.Errorf("failed to check popup: %w", err)
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 软删除相关的提交记录
	if err := tx.Where("popup_id = ?", id).Delete(&models.Submission{}).Error; err != nil {
		tx.Rollback()
		s.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to delete popup submissions")
		return fmt.Errorf("failed to delete popup submissions: %w", err)
	}

	// 软删除弹窗
	if err := tx.Delete(&popup).Error; err != nil {
		tx.Rollback()
		s.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to delete popup")
		return fmt.Errorf("failed to delete popup: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to commit popup deletion")
		return fmt.Errorf("failed to commit popup deletion: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"popup_id": id,
	}).Info("Popup deleted successfully")

	return nil
}

// ListPopups 获取弹窗列表
func (s *PopupCRUDService) ListPopups(page, pageSize int, proxyConfigID *uint, popupType string, enabled *bool) ([]*models.Popup, int64, error) {
	var popups []*models.Popup
	var total int64

	// 构建查询
	query := s.db.Model(&models.Popup{})

	// 添加过滤条件
	if proxyConfigID != nil {
		query = query.Where("proxy_config_id = ?", *proxyConfigID)
	}
	if popupType != "" {
		query = query.Where("popup_type = ?", popupType)
	}
	if enabled != nil {
		query = query.Where("is_active = ?", *enabled)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count popups: %w", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("ProxyConfig").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&popups).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get popups: %w", err)
	}

	return popups, total, nil
}

// TogglePopupStatus 切换弹窗状态
func (s *PopupCRUDService) TogglePopupStatus(id uint) error {
	// 检查弹窗是否存在
	var popup models.Popup
	if err := s.db.First(&popup, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("popup not found")
		}
		return fmt.Errorf("failed to check popup: %w", err)
	}

	// 切换状态
	newStatus := !popup.IsActive
	if err := s.db.Model(&popup).Updates(map[string]interface{}{
		"is_active":  newStatus,
		"updated_at": time.Now(),
	}).Error; err != nil {
		s.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to toggle popup status")
		return fmt.Errorf("failed to toggle popup status: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"popup_id": id,
		"enabled":  newStatus,
	}).Info("Popup status toggled successfully")

	return nil
}

// GetPopupsByProxyConfig 根据代理配置获取弹窗
func (s *PopupCRUDService) GetPopupsByProxyConfig(proxyConfigID uint) ([]*models.Popup, error) {
	var popups []*models.Popup
	if err := s.db.Where("proxy_config_id = ? AND is_active = ?", proxyConfigID, true).Order("created_at ASC").Find(&popups).Error; err != nil {
		return nil, fmt.Errorf("failed to get popups by proxy config: %w", err)
	}

	return popups, nil
}
