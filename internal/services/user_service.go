package services

import (
	"proxy-enhancer-ultra/internal/auth"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserService 用户服务 - 组合各个专门的服务以保持向后兼容性
type UserService struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
	logger     logger.Logger

	// 组合的专门服务
	authService     *AuthService
	crudService     *UserCRUDService
	passwordService *PasswordService
}

// NewUserService 创建新的用户服务
func NewUserService(db *gorm.DB, jwtManager *auth.JWTManager, logger logger.Logger) *UserService {
	return &UserService{
		db:         db,
		jwtManager: jwtManager,
		logger:     logger,

		// 初始化专门的服务
		authService:     NewAuthService(db, jwtManager, logger),
		crudService:     NewUserCRUDService(db, logger),
		passwordService: NewPasswordService(db, logger),
	}
}

// 为了向后兼容，这些类型现在从user_types.go重新导出
// 实际定义已移动到专门的类型文件中

// Login 用户登录 - 委托给认证服务
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	return s.authService.Login(req)
}

// CreateUser 创建用户 - 委托给CRUD服务
func (s *UserService) CreateUser(req *CreateUserRequest) (*UserInfo, error) {
	return s.crudService.CreateUser(req)
}

// GetUser 获取用户信息 - 委托给CRUD服务
func (s *UserService) GetUser(id uuid.UUID) (*UserInfo, error) {
	return s.crudService.GetUser(id)
}

// UpdateUser 更新用户信息 - 委托给CRUD服务
func (s *UserService) UpdateUser(id uuid.UUID, req *UpdateUserRequest) error {
	return s.crudService.UpdateUser(id, req)
}

// DeleteUser 删除用户 - 委托给CRUD服务
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.crudService.DeleteUser(id)
}

// ListUsers 获取用户列表 - 委托给CRUD服务
func (s *UserService) ListUsers(page, pageSize int, role string, enabled *bool) ([]*UserInfo, int64, error) {
	return s.crudService.ListUsers(page, pageSize, role, enabled)
}

// ChangePassword 修改密码 - 委托给密码服务
func (s *UserService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	return s.passwordService.ChangePassword(userID, oldPassword, newPassword)
}

// RefreshToken 刷新token - 委托给认证服务
func (s *UserService) RefreshToken(tokenString string) (*LoginResponse, error) {
	return s.authService.RefreshToken(tokenString)
}
