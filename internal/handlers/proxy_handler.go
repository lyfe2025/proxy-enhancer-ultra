package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/gorilla/mux"
)

// ProxyHandler 代理处理器
type ProxyHandler struct {
	proxyService *services.ProxyService
	logger       logger.Logger
}

// NewProxyHandler 创建新的代理处理器
func NewProxyHandler(proxyService *services.ProxyService, logger logger.Logger) *ProxyHandler {
	return &ProxyHandler{
		proxyService: proxyService,
		logger:       logger,
	}
}

// CreateProxyConfig 创建代理配置
// @Summary 创建代理配置
// @Description 创建新的代理配置，包括代理域名、目标URL等设置
// @Tags 代理管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param config body models.ProxyConfig true "代理配置"
// @Success 201 {object} Response{data=models.ProxyConfig} "创建成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 401 {object} Response "未授权"
// @Router /proxy/configs [post]
func (h *ProxyHandler) CreateProxyConfig(w http.ResponseWriter, r *http.Request) {
	var config models.ProxyConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.proxyService.CreateProxyConfig(&config); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error":  err.Error(),
			"domain": config.ProxyDomain,
		}).Error("Failed to create proxy config")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Proxy config created successfully",
		"data":    config,
	})
}

// GetProxyConfig 获取代理配置
func (h *ProxyHandler) GetProxyConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	config, err := h.proxyService.GetProxyConfig(uint(id))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to get proxy config")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": config,
	})
}

// UpdateProxyConfig 更新代理配置
func (h *ProxyHandler) UpdateProxyConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	var updates models.ProxyConfig
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.proxyService.UpdateProxyConfig(uint(id), &updates); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to update proxy config")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Proxy config updated successfully",
	})
}

// DeleteProxyConfig 删除代理配置
func (h *ProxyHandler) DeleteProxyConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	if err := h.proxyService.DeleteProxyConfig(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to delete proxy config")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Proxy config deleted successfully",
	})
}

// ListProxyConfigs 获取代理配置列表
func (h *ProxyHandler) ListProxyConfigs(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
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

	configs, total, err := h.proxyService.ListProxyConfigs(page, pageSize, enabled)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to list proxy configs")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve proxy configs")
		return
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": configs,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// ToggleProxyConfig 切换代理配置状态
func (h *ProxyHandler) ToggleProxyConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	if err := h.proxyService.ToggleProxyConfig(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to toggle proxy config")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Proxy config status toggled successfully",
	})
}

// GetProxyStats 获取代理统计信息
func (h *ProxyHandler) GetProxyStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	// 解析时间范围参数
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	var startTime, endTime time.Time

	if startTimeStr != "" {
		if st, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			startTime = st
		}
	}

	if endTimeStr != "" {
		if et, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			endTime = et
		}
	}

	stats, err := h.proxyService.GetProxyStats(uint(id), startTime, endTime)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to get proxy stats")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve proxy stats")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": stats,
	})
}

// respondWithError 返回错误响应
func (h *ProxyHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON 返回JSON响应
func (h *ProxyHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
