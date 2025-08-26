package handlers

import (
	"encoding/json"
	"net/http"

	"proxy-enhancer-ultra/internal/middleware"
	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"
)

// AuthHandler 认证处理器 - 专注于核心认证功能
type AuthHandler struct {
	BaseHandler
	userService *services.UserService
	logger      logger.Logger
}

// NewAuthHandler 创建新的认证处理器
func NewAuthHandler(userService *services.UserService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logger:      logger,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码进行登录认证
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "登录请求"
// @Success 200 {object} Response{data=services.LoginResponse} "登录成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 401 {object} Response "认证失败"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 基本验证
	if req.Username == "" || req.Password == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	response, err := h.userService.Login(&req)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"username": req.Username,
			"error":    err.Error(),
			"ip":       middleware.GetClientIP(r),
		}).Warn("Login failed")
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"username": req.Username,
		"user_id":  response.User.ID,
		"ip":       middleware.GetClientIP(r),
	}).Info("User logged in successfully")

	h.respondWithSuccess(w, http.StatusOK, "Login successful", response)
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新的用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body services.CreateUserRequest true "注册请求"
// @Success 201 {object} Response{data=services.UserInfo} "注册成功"
// @Failure 400 {object} Response "请求参数错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 基本验证
	if req.Username == "" || req.Email == "" || req.Password == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username, email and password are required")
		return
	}

	if len(req.Password) < 6 {
		h.respondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// 注册时默认为普通用户角色
	req.Role = "user"

	userInfo, err := h.userService.CreateUser(&req)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"username": req.Username,
			"email":    req.Email,
			"error":    err.Error(),
			"ip":       middleware.GetClientIP(r),
		}).Error("Registration failed")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"username": req.Username,
		"user_id":  userInfo.ID,
		"ip":       middleware.GetClientIP(r),
	}).Info("User registered successfully")

	h.respondWithSuccess(w, http.StatusCreated, "Registration successful", userInfo)
}

// RefreshToken 刷新token
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{token=string} true "刷新令牌请求"
// @Success 200 {object} Response{data=services.LoginResponse} "令牌刷新成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 401 {object} Response "令牌无效"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		h.respondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	response, err := h.userService.RefreshToken(req.Token)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to refresh token")
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.respondWithSuccess(w, http.StatusOK, "Token refreshed successfully", response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 登出当前用户，清除认证状态
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response "登出成功"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, _, _, ok := middleware.GetUserFromContext(r.Context())
	if ok {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"ip":      middleware.GetClientIP(r),
		}).Info("User logged out")
	}

	h.respondWithSuccess(w, http.StatusOK, "Logout successful", nil)
}
