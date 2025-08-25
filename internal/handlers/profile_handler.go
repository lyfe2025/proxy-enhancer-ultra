package handlers

import (
	"encoding/json"
	"net/http"

	"proxy-enhancer-ultra/internal/middleware"
	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"
)

// ProfileHandler 用户资料处理器
type ProfileHandler struct {
	BaseHandler
	userService *services.UserService
	logger      logger.Logger
}

// NewProfileHandler 创建新的用户资料处理器
func NewProfileHandler(userService *services.UserService, logger logger.Logger) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetProfile 获取用户资料
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "Profile retrieved successfully", userInfo)
}

// UpdateProfile 更新用户资料
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "Profile updated successfully", nil)
}

// ChangePassword 修改密码
func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "Password changed successfully", nil)
}
