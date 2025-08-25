package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/gorilla/mux"
)

// MonitoringHandler 监控处理器
type MonitoringHandler struct {
	monitoringService *services.MonitoringService
	logger            logger.Logger
}

// NewMonitoringHandler 创建新的监控处理器
func NewMonitoringHandler(monitoringService *services.MonitoringService, logger logger.Logger) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringService: monitoringService,
		logger:            logger,
	}
}

// GetSystemMetrics 获取系统指标
func (h *MonitoringHandler) GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.monitoringService.CollectSystemMetrics()
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to collect system metrics")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to collect system metrics")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": metrics,
	})
}

// GetProxyMetrics 获取代理指标
func (h *MonitoringHandler) GetProxyMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.monitoringService.CollectProxyMetrics()
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to collect proxy metrics")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to collect proxy metrics")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": metrics,
	})
}

// GetOverallStats 获取总体统计
func (h *MonitoringHandler) GetOverallStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.monitoringService.GetOverallStats()
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get overall stats")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get overall stats")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": stats,
	})
}

// GetSystemMetricsHistory 获取系统指标历史
func (h *MonitoringHandler) GetSystemMetricsHistory(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	hoursStr := r.URL.Query().Get("hours")
	hours := 24 // 默认24小时

	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil && h > 0 && h <= 168 { // 最多7天
			hours = h
		}
	}

	metrics, err := h.monitoringService.GetSystemMetricsHistory(hours)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"hours": hours,
			"error": err.Error(),
		}).Error("Failed to get system metrics history")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get system metrics history")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": metrics,
		"meta": map[string]interface{}{
			"hours": hours,
			"count": len(metrics),
		},
	})
}

// GetProxyMetricsHistory 获取代理指标历史
func (h *MonitoringHandler) GetProxyMetricsHistory(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	hoursStr := r.URL.Query().Get("hours")
	hours := 24 // 默认24小时

	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil && h > 0 && h <= 168 { // 最多7天
			hours = h
		}
	}

	metrics, err := h.monitoringService.GetProxyMetricsHistory(hours)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"hours": hours,
			"error": err.Error(),
		}).Error("Failed to get proxy metrics history")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get proxy metrics history")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": metrics,
		"meta": map[string]interface{}{
			"hours": hours,
			"count": len(metrics),
		},
	})
}

// GetProxyStats 获取特定代理的统计信息
func (h *MonitoringHandler) GetProxyStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proxyConfigID, err := strconv.ParseUint(vars["proxy_config_id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	stats, err := h.monitoringService.GetProxyStats(uint(proxyConfigID))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"proxy_config_id": proxyConfigID,
			"error":           err.Error(),
		}).Error("Failed to get proxy stats")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": stats,
	})
}

// CleanupOldMetrics 清理旧的指标数据
func (h *MonitoringHandler) CleanupOldMetrics(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	daysStr := r.URL.Query().Get("days")
	days := 30 // 默认30天

	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 { // 最多1年
			days = d
		}
	}

	if err := h.monitoringService.CleanupOldMetrics(days); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"days":  days,
			"error": err.Error(),
		}).Error("Failed to cleanup old metrics")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to cleanup old metrics")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Old metrics cleaned up successfully",
		"days":    days,
	})
}

// GetHealthCheck 健康检查
func (h *MonitoringHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	// 获取系统基本信息
	systemMetrics, err := h.monitoringService.CollectSystemMetrics()
	if err != nil {
		h.respondWithError(w, http.StatusServiceUnavailable, "System unhealthy")
		return
	}

	// 获取总体统计
	stats, err := h.monitoringService.GetOverallStats()
	if err != nil {
		h.respondWithError(w, http.StatusServiceUnavailable, "Database unhealthy")
		return
	}

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": systemMetrics.Timestamp,
		"system": map[string]interface{}{
			"memory_usage": systemMetrics.MemoryUsage,
			"goroutines":   systemMetrics.Goroutines,
		},
		"database": map[string]interface{}{
			"connected": true,
		},
		"services": map[string]interface{}{
			"total_proxies":  stats.TotalProxies,
			"active_proxies": stats.ActiveProxies,
			"total_users":    stats.TotalUsers,
		},
	}

	h.respondWithJSON(w, http.StatusOK, health)
}

// GetDashboardData 获取仪表板数据
func (h *MonitoringHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	// 获取总体统计
	stats, err := h.monitoringService.GetOverallStats()
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get overall stats for dashboard")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get dashboard data")
		return
	}

	// 获取系统指标
	systemMetrics, err := h.monitoringService.CollectSystemMetrics()
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get system metrics for dashboard")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get dashboard data")
		return
	}

	// 获取代理指标
	proxyMetrics, err := h.monitoringService.CollectProxyMetrics()
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get proxy metrics for dashboard")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get dashboard data")
		return
	}

	// 获取最近24小时的系统指标历史
	systemHistory, err := h.monitoringService.GetSystemMetricsHistory(24)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get system metrics history for dashboard")
		// 不返回错误，继续处理
		systemHistory = []*models.SystemMetric{}
	}

	// 获取最近24小时的代理指标历史
	proxyHistory, err := h.monitoringService.GetProxyMetricsHistory(24)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to get proxy metrics history for dashboard")
		// 不返回错误，继续处理
		proxyHistory = []*models.SystemMetric{}
	}

	dashboard := map[string]interface{}{
		"overview": stats,
		"current_metrics": map[string]interface{}{
			"system": systemMetrics,
			"proxy":  proxyMetrics,
		},
		"history": map[string]interface{}{
			"system": systemHistory,
			"proxy":  proxyHistory,
		},
		"timestamp": systemMetrics.Timestamp,
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": dashboard,
	})
}

// respondWithError 返回错误响应
func (h *MonitoringHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON 返回JSON响应
func (h *MonitoringHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}