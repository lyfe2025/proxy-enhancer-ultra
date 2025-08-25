package services

import (
	"errors"
	"fmt"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PasswordService 密码管理服务
type PasswordService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewPasswordService 创建新的密码管理服务
func NewPasswordService(db *gorm.DB, logger logger.Logger) *PasswordService {
	return &PasswordService{
		db:     db,
		logger: logger,
	}
}

// ChangePassword 修改密码
func (s *PasswordService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	// 获取用户
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	user.PasswordHash = string(hashedPassword)

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"username": user.Username,
	}).Info("Password changed successfully")

	return nil
}

// ResetPassword 重置密码（管理员功能）
func (s *PasswordService) ResetPassword(userID uuid.UUID, newPassword string) error {
	// 获取用户
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	user.PasswordHash = string(hashedPassword)

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"username": user.Username,
	}).Info("Password reset successfully")

	return nil
}

// ValidatePassword 验证密码强度
func (s *PasswordService) ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// 可以添加更多密码强度验证规则
	// 例如：必须包含大小写字母、数字、特殊字符等

	return nil
}
