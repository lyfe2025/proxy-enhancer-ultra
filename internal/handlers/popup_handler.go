package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/gorilla/mux"
)

// PopupHandler 弹窗处理器
type PopupHandler struct {
	popupService *services.PopupService
	logger       logger.Logger
}

// NewPopupHandler 创建新的弹窗处理器
func NewPopupHandler(popupService *services.PopupService, logger logger.Logger) *PopupHandler {
	return &PopupHandler{
		popupService: popupService,
		logger:       logger,
	}
}

// CreatePopup 创建弹窗
func (h *PopupHandler) CreatePopup(w http.ResponseWriter, r *http.Request) {
	var req services.CreatePopupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 基本验证
	if req.ProxyConfigID == 0 {
		h.respondWithError(w, http.StatusBadRequest, "Proxy config ID is required")
		return
	}

	if req.Title == "" {
		h.respondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.Content == "" {
		h.respondWithError(w, http.StatusBadRequest, "Content is required")
		return
	}

	if req.PopupType == "" {
		h.respondWithError(w, http.StatusBadRequest, "Popup type is required")
		return
	}

	if req.TriggerType == "" {
		h.respondWithError(w, http.StatusBadRequest, "Trigger type is required")
		return
	}

	popup, err := h.popupService.CreatePopup(&req)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"proxy_config_id": req.ProxyConfigID,
			"title":           req.Title,
			"popup_type":      req.PopupType,
			"error":           err.Error(),
		}).Error("Failed to create popup")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Popup created successfully",
		"data":    popup,
	})
}

// GetPopup 获取弹窗
func (h *PopupHandler) GetPopup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	popup, err := h.popupService.GetPopup(uint(id))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to get popup")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": popup,
	})
}

// UpdatePopup 更新弹窗
func (h *PopupHandler) UpdatePopup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	var req services.UpdatePopupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.popupService.UpdatePopup(uint(id), &req); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to update popup")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Popup updated successfully",
	})
}

// DeletePopup 删除弹窗
func (h *PopupHandler) DeletePopup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	if err := h.popupService.DeletePopup(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to delete popup")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Popup deleted successfully",
	})
}

// ListPopups 获取弹窗列表
func (h *PopupHandler) ListPopups(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	proxyConfigIDStr := r.URL.Query().Get("proxy_config_id")
	popupType := r.URL.Query().Get("popup_type")
	enabledStr := r.URL.Query().Get("enabled")

	// 设置默认值
	page := 1
	pageSize := 10
	var proxyConfigID *uint
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

	// 解析代理配置ID
	if proxyConfigIDStr != "" {
		if pcid, err := strconv.ParseUint(proxyConfigIDStr, 10, 32); err == nil {
			id := uint(pcid)
			proxyConfigID = &id
		}
	}

	// 解析启用状态
	if enabledStr != "" {
		if e, err := strconv.ParseBool(enabledStr); err == nil {
			enabled = &e
		}
	}

	popups, total, err := h.popupService.ListPopups(page, pageSize, proxyConfigID, popupType, enabled)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to list popups")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve popups")
		return
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": popups,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// TogglePopupStatus 切换弹窗状态
func (h *PopupHandler) TogglePopupStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	if err := h.popupService.TogglePopupStatus(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to toggle popup status")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Popup status toggled successfully",
	})
}

// GetPopupsByProxyConfig 根据代理配置获取弹窗
func (h *PopupHandler) GetPopupsByProxyConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proxyConfigID, err := strconv.ParseUint(vars["proxy_config_id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	popups, err := h.popupService.GetPopupsByProxyConfig(uint(proxyConfigID))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"proxy_config_id": proxyConfigID,
			"error":           err.Error(),
		}).Error("Failed to get popups by proxy config")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve popups")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": popups,
	})
}

// GetPopupStats 获取弹窗统计信息
func (h *PopupHandler) GetPopupStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid popup ID")
		return
	}

	stats, err := h.popupService.GetPopupStats(uint(id))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"popup_id": id,
			"error":    err.Error(),
		}).Error("Failed to get popup stats")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": stats,
	})
}

// respondWithError 返回错误响应
func (h *PopupHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON 返回JSON响应
func (h *PopupHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}