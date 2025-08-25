package handlers

import (
	"net/http"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"
)

// SubmissionHandler 提交处理器 - 组合各个专门的处理器以保持向后兼容性
type SubmissionHandler struct {
	BaseHandler
	submissionService *services.SubmissionService
	logger            logger.Logger

	// 组合的专门处理器
	crudHandler   *SubmissionCRUDHandler
	queryHandler  *SubmissionQueryHandler
	exportHandler *SubmissionExportHandler
	statsHandler  *SubmissionStatsHandler
}

// NewSubmissionHandler 创建新的提交处理器
func NewSubmissionHandler(submissionService *services.SubmissionService, logger logger.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
		logger:            logger,

		// 初始化专门的处理器
		crudHandler:   NewSubmissionCRUDHandler(submissionService, logger),
		queryHandler:  NewSubmissionQueryHandler(submissionService, logger),
		exportHandler: NewSubmissionExportHandler(submissionService, logger),
		statsHandler:  NewSubmissionStatsHandler(submissionService, logger),
	}
}

// CreateSubmission 创建提交 - 委托给CRUD处理器
func (h *SubmissionHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	h.crudHandler.CreateSubmission(w, r)
}

// GetSubmission 获取提交 - 委托给CRUD处理器
func (h *SubmissionHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	h.crudHandler.GetSubmission(w, r)
}

// UpdateSubmission 更新提交 - 委托给CRUD处理器
func (h *SubmissionHandler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	h.crudHandler.UpdateSubmission(w, r)
}

// DeleteSubmission 删除提交 - 委托给CRUD处理器
func (h *SubmissionHandler) DeleteSubmission(w http.ResponseWriter, r *http.Request) {
	h.crudHandler.DeleteSubmission(w, r)
}

// ListSubmissions 获取提交列表 - 委托给查询处理器
func (h *SubmissionHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	h.queryHandler.ListSubmissions(w, r)
}

// GetSubmissionsByPopup 根据弹窗获取提交列表 - 委托给查询处理器
func (h *SubmissionHandler) GetSubmissionsByPopup(w http.ResponseWriter, r *http.Request) {
	h.queryHandler.GetSubmissionsByPopup(w, r)
}

// GetSubmissionStats 获取提交统计信息 - 委托给统计处理器
func (h *SubmissionHandler) GetSubmissionStats(w http.ResponseWriter, r *http.Request) {
	h.statsHandler.GetSubmissionStats(w, r)
}

// GetSubmissionsByDateRange 根据日期范围获取提交 - 委托给查询处理器
func (h *SubmissionHandler) GetSubmissionsByDateRange(w http.ResponseWriter, r *http.Request) {
	h.queryHandler.GetSubmissionsByDateRange(w, r)
}

// ExportSubmissions 导出提交数据 - 委托给导出处理器
func (h *SubmissionHandler) ExportSubmissions(w http.ResponseWriter, r *http.Request) {
	h.exportHandler.ExportSubmissions(w, r)
}

// DeleteSubmissionsByPopup 根据弹窗删除所有提交 - 委托给CRUD处理器
func (h *SubmissionHandler) DeleteSubmissionsByPopup(w http.ResponseWriter, r *http.Request) {
	h.crudHandler.DeleteSubmissionsByPopup(w, r)
}

// 为了向后兼容，这些方法现在从 BaseHandler 重新导出
// 实际实现已移动到 common_handler.go 中
