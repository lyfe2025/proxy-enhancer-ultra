package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// UserAdminHandler 用户管理处理器（管理员功能）
type UserAdminHandler struct {
	BaseHandler
	userService *services.UserService
	logger      logger.Logger
}

// NewUserAdminHandler 创建新的用户管理处理器
func NewUserAdminHandler(userService *services.UserService, logger logger.Logger) *UserAdminHandler {
	return &UserAdminHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser 创建用户（管理员）
func (h *UserAdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusCreated, "User created successfully", userInfo)
}

// GetUser 获取用户信息（管理员）
func (h *UserAdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "User retrieved successfully", userInfo)
}

// UpdateUser 更新用户信息（管理员）
func (h *UserAdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "User updated successfully", nil)
}

// DeleteUser 删除用户（管理员）
func (h *UserAdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "User deleted successfully", nil)
}

// ListUsers 获取用户列表（管理员）
func (h *UserAdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
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

	responseData := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	}
	h.respondWithSuccess(w, http.StatusOK, "Users retrieved successfully", responseData)
}
