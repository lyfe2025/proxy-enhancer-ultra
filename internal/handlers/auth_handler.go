package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"proxy-enhancer-ultra/internal/middleware"
	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

// AuthHandler 认证处理器
type AuthHandler struct {
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

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"data":    response,
	})
}

// Register 用户注册
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

	h.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Registration successful",
		"data":    userInfo,
	})
}

// GetProfile 获取用户资料
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _, _, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userInfo, err := h.userService.GetUser(userID)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to get user profile")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": userInfo,
	})
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _, _, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req services.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 普通用户不能修改角色
	req.Role = ""
	req.Enabled = nil

	if err := h.userService.UpdateUser(userID, &req); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to update user profile")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Profile updated successfully",
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, _, _, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		h.respondWithError(w, http.StatusBadRequest, "Old password and new password are required")
		return
	}

	if len(req.NewPassword) < 6 {
		h.respondWithError(w, http.StatusBadRequest, "New password must be at least 6 characters")
		return
	}

	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to change password")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password changed successfully",
	})
}

// RefreshToken 刷新token
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

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Token refreshed successfully",
		"data":    response,
	})
}

// Admin endpoints

// CreateUser 创建用户（管理员）
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	userInfo, err := h.userService.CreateUser(&req)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"username": req.Username,
			"email":    req.Email,
			"error":    err.Error(),
		}).Error("Failed to create user")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"data":    userInfo,
	})
}

// GetUser 获取用户信息（管理员）
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userInfo, err := h.userService.GetUser(id)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Failed to get user")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": userInfo,
	})
}

// UpdateUser 更新用户信息（管理员）
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req services.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.userService.UpdateUser(id, &req); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Failed to update user")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
	})
}

// DeleteUser 删除用户（管理员）
func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Failed to delete user")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})
}

// ListUsers 获取用户列表（管理员）
func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	role := r.URL.Query().Get("role")
	enabledStr := r.URL.Query().Get("enabled")

	// 设置默认值
	page := 1
	pageSize := 10
	var enabled *bool

	// 解析页码
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 解析页面大小
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// 解析启用状态
	if enabledStr != "" {
		if e, err := strconv.ParseBool(enabledStr); err == nil {
			enabled = &e
		}
	}

	users, total, err := h.userService.ListUsers(page, pageSize, role, enabled)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to list users")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": users,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// respondWithError 返回错误响应
func (h *AuthHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON 返回JSON响应
func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}