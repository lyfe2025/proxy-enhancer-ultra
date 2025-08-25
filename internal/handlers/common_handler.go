package handlers

import (
	"encoding/json"
	"net/http"
)

// ApiResponse 标准API响应格式
type ApiResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// BaseHandler 基础处理器，提供通用响应方法
type BaseHandler struct{}

// respondWithError 返回错误响应
func (h *BaseHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := ApiResponse{
		Success: false,
		Code:    code,
		Message: message,
	}
	h.respondWithJSON(w, code, response)
}

// respondWithSuccess 返回成功响应
func (h *BaseHandler) respondWithSuccess(w http.ResponseWriter, code int, message string, data interface{}) {
	response := ApiResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	}
	h.respondWithJSON(w, code, response)
}

// respondWithJSON 返回JSON响应
func (h *BaseHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
