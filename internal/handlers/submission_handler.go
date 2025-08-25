package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

// SubmissionHandler 提交处理器
type SubmissionHandler struct {
	submissionService *services.SubmissionService
	logger            logger.Logger
}

// NewSubmissionHandler 创建新的提交处理器
func NewSubmissionHandler(submissionService *services.SubmissionService, logger logger.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
		logger:            logger,
	}
}

// CreateSubmission 创建提交
func (h *SubmissionHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Submission created successfully",
		"data":    submission,
	})
}

// GetSubmission 获取提交
func (h *SubmissionHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
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
			"error":        err.Error(),
		}).Error("Failed to get submission")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": submission,
	})
}

// UpdateSubmission 更新提交
func (h *SubmissionHandler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
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
			"error":        err.Error(),
		}).Error("Failed to update submission")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Submission updated successfully",
	})
}

// DeleteSubmission 删除提交
func (h *SubmissionHandler) DeleteSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid submission ID")
		return
	}

	if err := h.submissionService.DeleteSubmission(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"submission_id": id,
			"error":        err.Error(),
		}).Error("Failed to delete submission")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Submission deleted successfully",
	})
}

// ListSubmissions 获取提交列表
func (h *SubmissionHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	popupIDStr := r.URL.Query().Get("popup_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// 设置默认值
	page := 1
	pageSize := 10
	var popupID *uuid.UUID
	var startDate, endDate *time.Time

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

	// 解析弹窗ID
	if popupIDStr != "" {
		if pid, err := uuid.Parse(popupIDStr); err == nil {
			popupID = &pid
		}
	}

	// 解析开始日期
	if startDateStr != "" {
		if sd, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &sd
		}
	}

	// 解析结束日期
	if endDateStr != "" {
		if ed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// 设置为当天的最后一秒
			endOfDay := ed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &endOfDay
		}
	}

	submissions, total, err := h.submissionService.ListSubmissions(page, pageSize, popupID, startDate, endDate)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to list submissions")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve submissions")
		return
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": submissions,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetSubmissionsByPopup 根据弹窗获取提交列表
func (h *SubmissionHandler) GetSubmissionsByPopup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popupID, err := uuid.Parse(vars["popup_id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	submissions, err := h.submissionService.GetSubmissionsByPopup(popupID)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": popupID,
			"error":    err.Error(),
		}).Error("Failed to get submissions by popup")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve submissions")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": submissions,
	})
}

// GetSubmissionStats 获取提交统计信息
func (h *SubmissionHandler) GetSubmissionStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popupID, err := uuid.Parse(vars["popup_id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	stats, err := h.submissionService.GetSubmissionStats(popupID)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": popupID,
			"error":    err.Error(),
		}).Error("Failed to get submission stats")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": stats,
	})
}

// GetSubmissionsByDateRange 根据日期范围获取提交
func (h *SubmissionHandler) GetSubmissionsByDateRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popupID, err := uuid.Parse(vars["popup_id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	// 解析日期参数
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		h.respondWithError(w, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid start date format (YYYY-MM-DD)")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid end date format (YYYY-MM-DD)")
		return
	}

	// 设置结束日期为当天的最后一秒
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	submissions, err := h.submissionService.GetSubmissionsByDateRange(popupID, startDate, endDate)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id":   popupID,
			"start_date": startDate,
			"end_date":   endDate,
			"error":      err.Error(),
		}).Error("Failed to get submissions by date range")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve submissions")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": submissions,
	})
}

// ExportSubmissions 导出提交数据
func (h *SubmissionHandler) ExportSubmissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popupID, err := uuid.Parse(vars["popup_id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	// 获取导出格式
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json" // 默认格式
	}

	data, err := h.submissionService.ExportSubmissions(popupID, format)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": popupID,
			"format":   format,
			"error":    err.Error(),
		}).Error("Failed to export submissions")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 设置响应头
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=submissions.json")
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=submissions.csv")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DeleteSubmissionsByPopup 根据弹窗删除所有提交
func (h *SubmissionHandler) DeleteSubmissionsByPopup(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "All submissions for popup deleted successfully",
	})
}

// respondWithError 返回错误响应
func (h *SubmissionHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON 返回JSON响应
func (h *SubmissionHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}