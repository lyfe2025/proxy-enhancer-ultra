package handlers

import (
	"net/http"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// SubmissionStatsHandler 提交统计处理器
type SubmissionStatsHandler struct {
	BaseHandler
	submissionService *services.SubmissionService
	logger            logger.Logger
}

// NewSubmissionStatsHandler 创建新的提交统计处理器
func NewSubmissionStatsHandler(submissionService *services.SubmissionService, logger logger.Logger) *SubmissionStatsHandler {
	return &SubmissionStatsHandler{
		submissionService: submissionService,
		logger:            logger,
	}
}

// GetSubmissionStats 获取提交统计信息
func (h *SubmissionStatsHandler) GetSubmissionStats(w http.ResponseWriter, r *http.Request) {
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

	h.respondWithSuccess(w, http.StatusOK, "Submission statistics retrieved successfully", stats)
}

// GetOverallStats 获取总体统计信息
func (h *SubmissionStatsHandler) GetOverallStats(w http.ResponseWriter, r *http.Request) {
	// 这里可以实现获取全局统计信息的逻辑
	// 比如所有弹窗的提交总数、今日提交总数等

	// 暂时返回空的统计信息，实际实现需要在service层添加相应方法
	stats := map[string]interface{}{
		"total_submissions":         0,
		"today_submissions":         0,
		"active_popups":             0,
		"total_popups":              0,
		"average_daily_submissions": 0.0,
	}

	h.respondWithSuccess(w, http.StatusOK, "Overall statistics retrieved successfully", stats)
}

// GetSubmissionTrends 获取提交趋势数据
func (h *SubmissionStatsHandler) GetSubmissionTrends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popupID, err := uuid.Parse(vars["popup_id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	// 获取天数参数，默认30天
	days := 30
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		// 这里可以解析days参数
		// 暂时使用默认值
	}

	// 这里可以实现获取趋势数据的逻辑
	// 暂时返回模拟数据
	trends := []map[string]interface{}{
		{
			"date":  "2024-01-01",
			"count": 10,
		},
		{
			"date":  "2024-01-02",
			"count": 15,
		},
		// 更多趋势数据...
	}

	responseData := map[string]interface{}{
		"popup_id": popupID,
		"period":   days,
		"trends":   trends,
	}

	h.respondWithSuccess(w, http.StatusOK, "Submission trends retrieved successfully", responseData)
}
