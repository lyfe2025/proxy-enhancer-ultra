package handlers

import (
	"net/http"
	"strconv"
	"time"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// SubmissionQueryHandler 提交查询处理器
type SubmissionQueryHandler struct {
	BaseHandler
	submissionService *services.SubmissionService
	logger            logger.Logger
}

// NewSubmissionQueryHandler 创建新的提交查询处理器
func NewSubmissionQueryHandler(submissionService *services.SubmissionService, logger logger.Logger) *SubmissionQueryHandler {
	return &SubmissionQueryHandler{
		submissionService: submissionService,
		logger:            logger,
	}
}

// ListSubmissions 获取提交列表
func (h *SubmissionQueryHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
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

	responseData := map[string]interface{}{
		"data": submissions,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	h.respondWithSuccess(w, http.StatusOK, "Submissions retrieved successfully", responseData)
}

// GetSubmissionsByPopup 根据弹窗获取提交列表
func (h *SubmissionQueryHandler) GetSubmissionsByPopup(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "Submissions retrieved successfully", submissions)
}

// GetSubmissionsByDateRange 根据日期范围获取提交
func (h *SubmissionQueryHandler) GetSubmissionsByDateRange(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "Submissions retrieved successfully", submissions)
}
