package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"proxy-enhancer-ultra/internal/services"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/gorilla/mux"
)

// RuleHandler 规则处理器
type RuleHandler struct {
	ruleService *services.RuleService
	logger      logger.Logger
}

// NewRuleHandler 创建新的规则处理器
func NewRuleHandler(ruleService *services.RuleService, logger logger.Logger) *RuleHandler {
	return &RuleHandler{
		ruleService: ruleService,
		logger:      logger,
	}
}

// CreateRule 创建规则
func (h *RuleHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req services.CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 基本验证
	if req.ProxyConfigID == 0 {
		h.respondWithError(w, http.StatusBadRequest, "Proxy config ID is required")
		return
	}

	if req.RuleType == "" {
		h.respondWithError(w, http.StatusBadRequest, "Rule type is required")
		return
	}

	if req.Selector == "" {
		h.respondWithError(w, http.StatusBadRequest, "Selector is required")
		return
	}

	if req.Action == "" {
		h.respondWithError(w, http.StatusBadRequest, "Action is required")
		return
	}

	rule, err := h.ruleService.CreateRule(&req)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"proxy_config_id": req.ProxyConfigID,
			"rule_type":       req.RuleType,
			"error":           err.Error(),
		}).Error("Failed to create rule")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Rule created successfully",
		"data":    rule,
	})
}

// GetRule 获取规则
func (h *RuleHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	rule, err := h.ruleService.GetRule(uint(id))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"rule_id": id,
			"error":   err.Error(),
		}).Error("Failed to get rule")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": rule,
	})
}

// UpdateRule 更新规则
func (h *RuleHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	var req services.UpdateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.ruleService.UpdateRule(uint(id), &req); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"rule_id": id,
			"error":   err.Error(),
		}).Error("Failed to update rule")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Rule updated successfully",
	})
}

// DeleteRule 删除规则
func (h *RuleHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	if err := h.ruleService.DeleteRule(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"rule_id": id,
			"error":   err.Error(),
		}).Error("Failed to delete rule")
		h.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Rule deleted successfully",
	})
}

// ListRules 获取规则列表
func (h *RuleHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	proxyConfigIDStr := r.URL.Query().Get("proxy_config_id")
	ruleType := r.URL.Query().Get("rule_type")
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

	rules, total, err := h.ruleService.ListRules(page, pageSize, proxyConfigID, ruleType, enabled)
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to list rules")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve rules")
		return
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": rules,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// ToggleRuleStatus 切换规则状态
func (h *RuleHandler) ToggleRuleStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	if err := h.ruleService.ToggleRuleStatus(uint(id)); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"rule_id": id,
			"error":   err.Error(),
		}).Error("Failed to toggle rule status")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Rule status toggled successfully",
	})
}

// GetRulesByProxyConfig 根据代理配置获取规则
func (h *RuleHandler) GetRulesByProxyConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proxyConfigID, err := strconv.ParseUint(vars["proxy_config_id"], 10, 32)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid proxy config ID")
		return
	}

	rules, err := h.ruleService.GetRulesByProxyConfig(uint(proxyConfigID))
	if err != nil {
		h.logger.WithFields(map[string]interface{}{
			"proxy_config_id": proxyConfigID,
			"error":           err.Error(),
		}).Error("Failed to get rules by proxy config")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve rules")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": rules,
	})
}

// UpdateRulePriorities 批量更新规则优先级
func (h *RuleHandler) UpdateRulePriorities(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Updates []struct {
			ID       uint `json:"id"`
			Priority int  `json:"priority"`
		} `json:"updates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Updates) == 0 {
		h.respondWithError(w, http.StatusBadRequest, "No updates provided")
		return
	}

	// 验证更新数据
	for _, update := range req.Updates {
		if update.ID == 0 {
			h.respondWithError(w, http.StatusBadRequest, "Invalid rule ID in updates")
			return
		}
	}

	if err := h.ruleService.UpdateRulePriorities(req.Updates); err != nil {
		h.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to update rule priorities")
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Rule priorities updated successfully",
	})
}

// respondWithError 返回错误响应
func (h *RuleHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON 返回JSON响应
func (h *RuleHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}