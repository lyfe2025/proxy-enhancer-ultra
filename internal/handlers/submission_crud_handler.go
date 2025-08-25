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

// SubmissionCRUDHandler 提交CRUD操作处理器
type SubmissionCRUDHandler struct {
	BaseHandler
	submissionService *services.SubmissionService
	logger            logger.Logger
}

// NewSubmissionCRUDHandler 创建新的提交CRUD处理器
func NewSubmissionCRUDHandler(submissionService *services.SubmissionService, logger logger.Logger) *SubmissionCRUDHandler {
	return &SubmissionCRUDHandler{
		submissionService: submissionService,
		logger:            logger,
	}
}

// CreateSubmission 创建提交
func (h *SubmissionCRUDHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req services.CreateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 基本验证
	if req.PopupID == uuid.Nil {
		h.respondWithError(w, http.StatusBadRequest, "Popup ID is required")
		return
	}

	if req.FormData == nil || len(req.FormData) == 0 {
		h.respondWithError(w, http.StatusBadRequest, "Form data is required")
		return
	}

	submission, err := h.submissionService.CreateSubmission(&req)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id":   req.PopupID,
			"ip_address": req.IPAddress,
			"error":      err.Error(),
		}).Error("Failed to create submission")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithSuccess(w, http.StatusCreated, "Submission created successfully", submission)
}

// GetSubmission 获取提交
func (h *SubmissionCRUDHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid submission ID")
		return
	}

	submission, err := h.submissionService.GetSubmission(uint(id))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"submission_id": id,
			"error":         err.Error(),
		}).Error("Failed to get submission")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithSuccess(w, http.StatusOK, "Submission retrieved successfully", submission)
}

// UpdateSubmission 更新提交
func (h *SubmissionCRUDHandler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid submission ID")
		return
	}

	var req services.UpdateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.submissionService.UpdateSubmission(uint(id), &req); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"submission_id": id,
			"error":         err.Error(),
		}).Error("Failed to update submission")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithSuccess(w, http.StatusOK, "Submission updated successfully", nil)
}

// DeleteSubmission 删除提交
func (h *SubmissionCRUDHandler) DeleteSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid submission ID")
		return
	}

	if err := h.submissionService.DeleteSubmission(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"submission_id": id,
			"error":         err.Error(),
		}).Error("Failed to delete submission")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithSuccess(w, http.StatusOK, "Submission deleted successfully", nil)
}

// DeleteSubmissionsByPopup 根据弹窗删除所有提交
func (h *SubmissionCRUDHandler) DeleteSubmissionsByPopup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popupID, err := uuid.Parse(vars["popup_id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	if err := h.submissionService.DeleteSubmissionsByPopup(popupID); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": popupID,
			"error":    err.Error(),
		}).Error("Failed to delete submissions by popup")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithSuccess(w, http.StatusOK, "All submissions for popup deleted successfully", nil)
}
