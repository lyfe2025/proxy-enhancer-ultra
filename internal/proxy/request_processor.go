package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"
)

// RequestProcessor 请求处理器
type RequestProcessor struct {
	logger logger.Logger
}

// NewRequestProcessor 创建新的请求处理器
func NewRequestProcessor(logger logger.Logger) *RequestProcessor {
	return &RequestProcessor{
		logger: logger,
	}
}

// CreateProxyRequest 创建代理请求
func (rp *RequestProcessor) CreateProxyRequest(originalReq *http.Request, targetURL string) (*http.Request, error) {
	// 解析目标URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	// 创建新请求
	proxyReq, err := http.NewRequest(originalReq.Method, targetURL, originalReq.Body)
	if err != nil {
		return nil, err
	}

	// 复制请求头
	for key, values := range originalReq.Header {
		// 跳过一些不应该转发的头
		if rp.shouldSkipHeader(key) {
			continue
		}
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// 设置Host头
	proxyReq.Host = parsedURL.Host
	proxyReq.Header.Set("Host", parsedURL.Host)

	// 设置X-Forwarded-For头
	clientIP := rp.GetClientIP(originalReq)
	if existingXFF := proxyReq.Header.Get("X-Forwarded-For"); existingXFF != "" {
		proxyReq.Header.Set("X-Forwarded-For", existingXFF+", "+clientIP)
	} else {
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// 设置其他代理头
	proxyReq.Header.Set("X-Forwarded-Proto", rp.getScheme(originalReq))
	proxyReq.Header.Set("X-Real-IP", clientIP)

	return proxyReq, nil
}

// ProcessResponse 处理响应内容
func (rp *RequestProcessor) ProcessResponse(resp *http.Response, proxyConfig *models.ProxyConfig, htmlInjector *HTMLInjector, urlRewriter *URLRewriter) ([]byte, error) {
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 只处理HTML内容
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		return body, nil
	}

	// 处理HTML内容
	return htmlInjector.ProcessHTML(body, proxyConfig, urlRewriter)
}

// CopyResponse 复制完整响应
func (rp *RequestProcessor) CopyResponse(w http.ResponseWriter, resp *http.Response) error {
	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	_, err := io.Copy(w, resp.Body)
	return err
}

// CopyResponseHeaders 复制响应头
func (rp *RequestProcessor) CopyResponseHeaders(w http.ResponseWriter, resp *http.Response, contentLength int) {
	for key, values := range resp.Header {
		// 跳过一些不应该复制的头
		if rp.shouldSkipResponseHeader(key) {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 更新Content-Length
	if contentLength > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	}
}

// shouldSkipHeader 判断是否应该跳过某个请求头
func (rp *RequestProcessor) shouldSkipHeader(key string) bool {
	skipHeaders := []string{
		"Connection",
		"Proxy-Connection",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	for _, skipHeader := range skipHeaders {
		if strings.EqualFold(key, skipHeader) {
			return true
		}
	}
	return false
}

// shouldSkipResponseHeader 判断是否应该跳过某个响应头
func (rp *RequestProcessor) shouldSkipResponseHeader(key string) bool {
	skipHeaders := []string{
		"Connection",
		"Transfer-Encoding",
		"Content-Length",          // 我们会重新设置
		"Content-Security-Policy", // 可能会阻止我们的脚本
	}

	for _, skipHeader := range skipHeaders {
		if strings.EqualFold(key, skipHeader) {
			return true
		}
	}
	return false
}

// GetClientIP 获取客户端IP
func (rp *RequestProcessor) GetClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For头获取
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// 尝试从X-Real-IP头获取
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 从RemoteAddr获取
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}

	return r.RemoteAddr
}

// getScheme 获取请求协议
func (rp *RequestProcessor) getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}
