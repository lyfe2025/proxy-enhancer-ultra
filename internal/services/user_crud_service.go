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

// UserCRUDService 用户CRUD操作服务
type UserCRUDService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewUserCRUDService 创建新的用户CRUD服务
func NewUserCRUDService(db *gorm.DB, logger logger.Logger) *UserCRUDService {
	return &UserCRUDService{
		db:     db,
		logger: logger,
	}
}

// CreateUser 创建用户
func (s *UserCRUDService) CreateUser(req *CreateUserRequest) (*UserInfo, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	err := s.db.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("username already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 检查邮箱是否已存在
	err = s.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 设置默认角色
	role := req.Role
	if role == "" {
		role = "user"
	}

	// 创建用户
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsActive:     true, // 使用IsActive字段
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("User created successfully")

	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     "user",        // 默认角色
		Enabled:  user.IsActive, // 使用IsActive字段
	}, nil
}

// GetUser 获取用户信息
func (s *UserCRUDService) GetUser(id uuid.UUID) (*UserInfo, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     "user",        // 默认角色
		Enabled:  user.IsActive, // 使用IsActive字段
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserCRUDService) UpdateUser(id uuid.UUID, req *UpdateUserRequest) error {
	// 检查用户是否存在
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// 更新字段
	updates := make(map[string]interface{})

	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var existingUser models.User
		err := s.db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error
		if err == nil {
			return errors.New("email already exists")
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("database error: %w", err)
		}
		updates["email"] = req.Email
	}

	if req.Password != "" {
		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		updates["password_hash"] = string(hashedPassword)
	}

	if req.Enabled != nil {
		updates["is_active"] = *req.Enabled
	}

	// 执行更新
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  id,
		"username": user.Username,
	}).Info("User updated successfully")

	return nil
}

// DeleteUser 删除用户
func (s *UserCRUDService) DeleteUser(id uuid.UUID) error {
	// 检查用户是否存在
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// 软删除用户
	if err := s.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  id,
		"username": user.Username,
	}).Info("User deleted successfully")

	return nil
}

// ListUsers 获取用户列表
func (s *UserCRUDService) ListUsers(page, pageSize int, role string, enabled *bool) ([]*UserInfo, int64, error) {
	var users []*models.User
	var total int64

	query := s.db.Model(&models.User{})

	// 添加过滤条件
	if enabled != nil {
		query = query.Where("is_active = ?", *enabled)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换为UserInfo
	userInfos := make([]*UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     "user",        // 默认角色
			Enabled:  user.IsActive, // 使用IsActive字段
		}
	}

	return userInfos, total, nil
}
