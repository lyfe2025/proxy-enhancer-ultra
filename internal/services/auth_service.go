package services

import (
	"errors"
	"fmt"
	"time"

	"proxy-enhancer-ultra/internal/auth"
	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
	logger     logger.Logger
}

// NewAuthService 创建新的认证服务
func NewAuthService(db *gorm.DB, jwtManager *auth.JWTManager, logger logger.Logger) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户（支持用户名或邮箱登录）
	var user models.User
	if err := s.db.Where("(username = ? OR email = ?) AND is_active = ?", req.Username, req.Username, true).First(&user).Error; err != nil {
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

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(tokenString string) (*LoginResponse, error) {
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
	userCrud := NewUserCRUDService(s.db, s.logger)
	userInfo, err := userCrud.GetUser(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     newToken,
		User:      userInfo,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}
