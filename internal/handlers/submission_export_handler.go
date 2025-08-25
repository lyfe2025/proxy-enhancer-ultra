package handlers

import (
	"net/http"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// SubmissionExportHandler 提交导出处理器
type SubmissionExportHandler struct {
	BaseHandler
	submissionService *services.SubmissionService
	logger            logger.Logger
}

// NewSubmissionExportHandler 创建新的提交导出处理器
func NewSubmissionExportHandler(submissionService *services.SubmissionService, logger logger.Logger) *SubmissionExportHandler {
	return &SubmissionExportHandler{
		submissionService: submissionService,
		logger:            logger,
	}
}

// ExportSubmissions 导出提交数据
func (h *SubmissionExportHandler) ExportSubmissions(w http.ResponseWriter, r *http.Request) {
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

	// 验证导出格式
	if !h.isValidExportFormat(format) {
		h.respondWithError(w, http.StatusBadRequest, "Invalid export format. Supported formats: json, csv")
		return
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
	h.setExportHeaders(w, format)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// ExportSubmissionsWithDateRange 按日期范围导出提交数据
func (h *SubmissionExportHandler) ExportSubmissionsWithDateRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popupID, err := uuid.Parse(vars["popup_id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	// 获取导出格式和日期范围
	format := r.URL.Query().Get("format")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if format == "" {
		format = "json"
	}

	if !h.isValidExportFormat(format) {
		h.respondWithError(w, http.StatusBadRequest, "Invalid export format. Supported formats: json, csv")
		return
	}

	// 验证日期参数
	if startDate == "" || endDate == "" {
		h.respondWithError(w, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	// 这里可以调用服务层的按日期范围导出方法
	// 目前使用现有的导出方法
	data, err := h.submissionService.ExportSubmissions(popupID, format)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id":   popupID,
			"format":     format,
			"start_date": startDate,
			"end_date":   endDate,
			"error":      err.Error(),
		}).Error("Failed to export submissions with date range")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 设置响应头
	h.setExportHeaders(w, format)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// isValidExportFormat 检查导出格式是否有效
func (h *SubmissionExportHandler) isValidExportFormat(format string) bool {
	validFormats := []string{"json", "csv", "xlsx"}
	for _, validFormat := range validFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

// setExportHeaders 设置导出文件的响应头
func (h *SubmissionExportHandler) setExportHeaders(w http.ResponseWriter, format string) {
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=submissions.json")
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=submissions.csv")
	case "xlsx":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=submissions.xlsx")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}
