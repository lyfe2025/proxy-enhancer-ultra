package services

import (
	"errors"
	"fmt"
	"time"

	"proxy-enhancer-ultra/internal/auth"
	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
	logger     logger.Logger
}

// NewUserService 创建新的用户服务
func NewUserService(db *gorm.DB, jwtManager *auth.JWTManager, logger logger.Logger) *UserService {
	return &UserService{
		db:         db,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string      `json:"token"`
	User      *UserInfo   `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Enabled  bool      `json:"enabled"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	var user models.User
	err := s.db.Where("username = ? AND enabled = ?", req.Username, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.WithFields(map[string]interface{}{
			"username": req.Username,
			"ip":       "unknown", // 可以从上下文获取
		}).Warn("Failed login attempt")
		return nil, errors.New("invalid username or password")
	}

	// 生成JWT token (暂时使用默认角色)
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 更新最后登录时间
	user.LastLoginAt = &time.Time{}
	*user.LastLoginAt = time.Now()
	s.db.Save(&user)

	s.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("User logged in successfully")

	return &LoginResponse{
		Token: token,
		User: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     "user",        // 默认角色
			Enabled:  user.IsActive, // 使用IsActive字段
		},
		ExpiresAt: time.Now().Add(24 * time.Hour), // 假设token有效期为24小时
	}, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *CreateUserRequest) (*UserInfo, error) {
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
func (s *UserService) GetUser(id uuid.UUID) (*UserInfo, error) {
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
func (s *UserService) UpdateUser(id uuid.UUID, req *UpdateUserRequest) error {
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

	// Role字段已移除，不再支持更新
	// if req.Role != "" {
	//	updates["role"] = req.Role
	// }

	if req.Enabled != nil {
		updates["is_active"] = *req.Enabled
	}

	// updated_at字段由BaseModel自动管理
	// updates["updated_at"] = time.Now()

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
func (s *UserService) DeleteUser(id uuid.UUID) error {
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
func (s *UserService) ListUsers(page, pageSize int, role string, enabled *bool) ([]*UserInfo, int64, error) {
	var users []*models.User
	var total int64

	query := s.db.Model(&models.User{})

	// 添加过滤条件
	// Role字段已移除，不再支持按角色过滤
	// if role != "" {
	//	query = query.Where("role = ?", role)
	// }
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

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
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
	// UpdatedAt字段由BaseModel自动管理

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  userID,
		"username": user.Username,
	}).Info("Password changed successfully")

	return nil
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(tokenString string) (*LoginResponse, error) {
	// 验证并刷新token
	newToken, err := s.jwtManager.RefreshToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// 从旧token中提取用户信息
	claims, err := s.jwtManager.ExtractClaims(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	// 获取最新的用户信息
	userInfo, err := s.GetUser(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     newToken,
		User:      userInfo,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}