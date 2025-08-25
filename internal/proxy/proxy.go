package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"proxy-enhancer-ultra/internal/config"
	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

// ProxyServer 反向代理服务器
type ProxyServer struct {
	db     *gorm.DB
	logger logger.Logger
	config *config.Config
	client *http.Client
}

// NewProxyServer 创建新的代理服务器
func NewProxyServer(db *gorm.DB, logger logger.Logger, cfg *config.Config) *ProxyServer {
	client := &http.Client{
		Timeout: 30 * time.Second, // 默认30秒超时
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 限制重定向次数
			if len(via) >= 10 {
				return fmt.Errorf("stopped after %d redirects", 10)
			}
			return nil
		},
	}

	return &ProxyServer{
		db:     db,
		logger: logger,
		config: cfg,
		client: client,
	}
}

// ServeHTTP 处理HTTP请求的核心方法
func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 记录请求开始时间
	startTime := time.Now()
	
	// 获取代理配置
	proxyConfig, err := p.getProxyConfig(r.Host)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"host": r.Host,
			"error": err.Error(),
		}).Error("Failed to get proxy config")
		http.Error(w, "Proxy configuration not found", http.StatusNotFound)
		return
	}

	// 构建目标URL
	targetURL := p.buildTargetURL(proxyConfig.TargetURL, r)
	
	// 创建代理请求
	proxyReq, err := p.createProxyRequest(r, targetURL)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"target_url": targetURL,
		}).Error("Failed to create proxy request")
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// 发送代理请求
	resp, err := p.client.Do(proxyReq)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"target_url": targetURL,
		}).Error("Failed to send proxy request")
		http.Error(w, "Failed to connect to target server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to read response body")
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// 处理响应内容
	processedBody, err := p.processResponse(body, resp.Header.Get("Content-Type"), proxyConfig)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Warn("Failed to process response, serving original")
		processedBody = body
	}

	// 复制响应头
	p.copyResponseHeaders(w, resp, len(processedBody))
	
	// 写入响应
	w.WriteHeader(resp.StatusCode)
	w.Write(processedBody)

	// 记录代理日志
	p.logProxyRequest(r, resp, proxyConfig.BaseModel.ID, time.Since(startTime))
}

// getProxyConfig 获取代理配置
func (p *ProxyServer) getProxyConfig(host string) (*models.ProxyConfig, error) {
	// 如果没有数据库连接，返回默认配置
	if p.db == nil {
		// 生成一个默认的UUID
		defaultID := uuid.New()
		return &models.ProxyConfig{
			BaseModel: models.BaseModel{
				ID: defaultID,
			},
			ProxyDomain: host,
			TargetURL:   "https://www.baidu.com", // 默认代理到百度
			IsActive:    true,
		}, nil
	}
	
	var proxyConfig models.ProxyConfig
	err := p.db.Where("proxy_domain = ? AND is_active = ?", host, true).First(&proxyConfig).Error
	if err != nil {
		return nil, err
	}
	return &proxyConfig, nil
}

// buildTargetURL 构建目标URL
func (p *ProxyServer) buildTargetURL(targetURL string, r *http.Request) string {
	if !strings.HasSuffix(targetURL, "/") {
		targetURL += "/"
	}
	
	path := strings.TrimPrefix(r.URL.Path, "/")
	fullURL := targetURL + path
	
	if r.URL.RawQuery != "" {
		fullURL += "?" + r.URL.RawQuery
	}
	
	return fullURL
}

// createProxyRequest 创建代理请求
func (p *ProxyServer) createProxyRequest(originalReq *http.Request, targetURL string) (*http.Request, error) {
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
		if p.shouldSkipHeader(key) {
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
	clientIP := p.getClientIP(originalReq)
	if existingXFF := proxyReq.Header.Get("X-Forwarded-For"); existingXFF != "" {
		proxyReq.Header.Set("X-Forwarded-For", existingXFF+", "+clientIP)
	} else {
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}
	
	// 设置其他代理头
	proxyReq.Header.Set("X-Forwarded-Proto", p.getScheme(originalReq))
	proxyReq.Header.Set("X-Real-IP", clientIP)

	return proxyReq, nil
}

// shouldSkipHeader 判断是否应该跳过某个请求头
func (p *ProxyServer) shouldSkipHeader(key string) bool {
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

// getClientIP 获取客户端IP
func (p *ProxyServer) getClientIP(r *http.Request) string {
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
func (p *ProxyServer) getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}

// processResponse 处理响应内容
func (p *ProxyServer) processResponse(body []byte, contentType string, proxyConfig *models.ProxyConfig) ([]byte, error) {
	// 只处理HTML内容
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		return body, nil
	}

	// 解析HTML
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return body, err
	}

	// 注入自定义脚本和样式
	p.injectCustomContent(doc, proxyConfig)
	
	// 修改链接和资源路径
	p.rewriteURLs(doc, proxyConfig)

	// 将修改后的HTML转换回字节
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return body, err
	}

	return buf.Bytes(), nil
}

// injectCustomContent 注入自定义内容
func (p *ProxyServer) injectCustomContent(doc *html.Node, proxyConfig *models.ProxyConfig) {
	// 查找head标签
	headNode := p.findNode(doc, "head")
	if headNode == nil {
		return
	}

	// 注入基础样式和脚本
	p.injectBaseAssets(headNode, proxyConfig)
	
	// 查找body标签
	bodyNode := p.findNode(doc, "body")
	if bodyNode == nil {
		return
	}

	// 注入弹窗和交互功能
	p.injectInteractiveElements(bodyNode, proxyConfig)
}

// injectBaseAssets 注入基础资源
func (p *ProxyServer) injectBaseAssets(headNode *html.Node, proxyConfig *models.ProxyConfig) {
	// 注入CSS样式
	styleNode := &html.Node{
		Type: html.ElementNode,
		Data: "style",
		Attr: []html.Attribute{{Key: "type", Val: "text/css"}},
	}
	styleContent := &html.Node{
		Type: html.TextNode,
		Data: p.getInjectedCSS(),
	}
	styleNode.AppendChild(styleContent)
	headNode.AppendChild(styleNode)

	// 注入JavaScript
	scriptNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{{Key: "type", Val: "text/javascript"}},
	}
	scriptContent := &html.Node{
		Type: html.TextNode,
		Data: p.getInjectedJS(proxyConfig),
	}
	scriptNode.AppendChild(scriptContent)
	headNode.AppendChild(scriptNode)
}

// injectInteractiveElements 注入交互元素
func (p *ProxyServer) injectInteractiveElements(bodyNode *html.Node, proxyConfig *models.ProxyConfig) {
	// 如果没有数据库连接，跳过弹窗注入
	if p.db == nil {
		return
	}
	
	// 获取该代理配置的弹窗规则
	var popups []models.Popup
	p.db.Where("proxy_config_id = ? AND enabled = ?", proxyConfig.ID, true).Find(&popups)

	for _, popup := range popups {
		// 创建弹窗容器
		popupContainer := p.createPopupElement(&popup)
		bodyNode.AppendChild(popupContainer)
	}
}

// createPopupElement 创建弹窗元素
func (p *ProxyServer) createPopupElement(popup *models.Popup) *html.Node {
	// 创建弹窗容器
	container := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{Key: "id", Val: fmt.Sprintf("popup-%d", popup.ID)},
			{Key: "class", Val: "proxy-popup"},
			{Key: "style", Val: "display: none;"},
		},
	}

	// 添加弹窗内容
	content := &html.Node{
		Type: html.RawNode,
		Data: popup.Content,
	}
	container.AppendChild(content)

	return container
}

// rewriteURLs 重写URL
func (p *ProxyServer) rewriteURLs(doc *html.Node, proxyConfig *models.ProxyConfig) {
	// 重写所有链接和资源引用
	p.walkNodes(doc, func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "a", "link":
				p.rewriteAttribute(node, "href", proxyConfig)
			case "img", "script":
				p.rewriteAttribute(node, "src", proxyConfig)
			case "form":
				p.rewriteAttribute(node, "action", proxyConfig)
			}
		}
	})
}

// rewriteAttribute 重写属性
func (p *ProxyServer) rewriteAttribute(node *html.Node, attrName string, proxyConfig *models.ProxyConfig) {
	for i, attr := range node.Attr {
		if attr.Key == attrName {
			node.Attr[i].Val = p.rewriteURL(attr.Val, proxyConfig)
			break
		}
	}
}

// rewriteURL 重写单个URL
func (p *ProxyServer) rewriteURL(originalURL string, proxyConfig *models.ProxyConfig) string {
	// 如果是相对URL，转换为绝对URL
	if strings.HasPrefix(originalURL, "/") {
		return "https://" + proxyConfig.ProxyDomain + originalURL
	}
	
	// 如果是目标域名的URL，替换为代理域名
	targetHost := p.extractHost(proxyConfig.TargetURL)
	if strings.Contains(originalURL, targetHost) {
		return strings.Replace(originalURL, targetHost, proxyConfig.ProxyDomain, 1)
	}
	
	return originalURL
}

// extractHost 从URL中提取主机名
func (p *ProxyServer) extractHost(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

// findNode 查找指定标签的节点
func (p *ProxyServer) findNode(doc *html.Node, tagName string) *html.Node {
	var result *html.Node
	p.walkNodes(doc, func(node *html.Node) {
		if result == nil && node.Type == html.ElementNode && node.Data == tagName {
			result = node
		}
	})
	return result
}

// walkNodes 遍历所有节点
func (p *ProxyServer) walkNodes(node *html.Node, fn func(*html.Node)) {
	fn(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		p.walkNodes(child, fn)
	}
}

// copyResponseHeaders 复制响应头
func (p *ProxyServer) copyResponseHeaders(w http.ResponseWriter, resp *http.Response, contentLength int) {
	for key, values := range resp.Header {
		// 跳过一些不应该复制的头
		if p.shouldSkipResponseHeader(key) {
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

// shouldSkipResponseHeader 判断是否应该跳过某个响应头
func (p *ProxyServer) shouldSkipResponseHeader(key string) bool {
	skipHeaders := []string{
		"Connection",
		"Transfer-Encoding",
		"Content-Length", // 我们会重新设置
		"Content-Security-Policy", // 可能会阻止我们的脚本
	}
	
	for _, skipHeader := range skipHeaders {
		if strings.EqualFold(key, skipHeader) {
			return true
		}
	}
	return false
}

// logProxyRequest 记录代理请求日志
func (p *ProxyServer) logProxyRequest(r *http.Request, resp *http.Response, proxyConfigID uuid.UUID, duration time.Duration) {
	// 如果没有数据库连接，只记录到日志文件
	if p.db == nil {
		p.logger.WithFields(map[string]interface{}{
			"method":       r.Method,
			"url":          r.URL.String(),
			"user_agent":   r.Header.Get("User-Agent"),
			"user_ip":      p.getClientIP(r),
			"status_code":  resp.StatusCode,
			"response_time": duration.Milliseconds(),
		}).Info("Proxy request (no database)")
		return
	}
	
	proxyLog := &models.ProxyLog{
		ProxyConfigID: proxyConfigID,
		Method:        r.Method,
		URL:           r.URL.String(),
		UserAgent:     r.Header.Get("User-Agent"),
		UserIP:        p.getClientIP(r),
		StatusCode:    resp.StatusCode,
		ResponseTime:  duration.Milliseconds(),
	}
	
	// 异步保存日志，避免影响响应性能
	go func() {
		if err := p.db.Create(proxyLog).Error; err != nil {
			p.logger.WithFields(map[string]interface{}{
				"error": err.Error(),
			}).Error("Failed to save proxy log")
		}
	}()
}

// getInjectedCSS 获取注入的CSS样式
func (p *ProxyServer) getInjectedCSS() string {
	return `
/* 代理增强器样式 */
.proxy-popup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: white;
    border: 1px solid #ccc;
    border-radius: 8px;
    box-shadow: 0 4px 20px rgba(0,0,0,0.3);
    z-index: 10000;
    max-width: 500px;
    max-height: 80vh;
    overflow-y: auto;
    padding: 20px;
}

.proxy-popup-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.5);
    z-index: 9999;
}

.proxy-popup-close {
    position: absolute;
    top: 10px;
    right: 15px;
    font-size: 24px;
    cursor: pointer;
    color: #999;
}

.proxy-popup-close:hover {
    color: #333;
}
`
}

// getInjectedJS 获取注入的JavaScript代码
func (p *ProxyServer) getInjectedJS(proxyConfig *models.ProxyConfig) string {
	return fmt.Sprintf(`
// 代理增强器脚本
(function() {
    'use strict';
    
    // 代理配置
    window.PROXY_CONFIG = {
        domain: '%s',
        configId: '%s'
    };
    
    // 弹窗管理器
    window.ProxyPopupManager = {
        show: function(popupId) {
            var popup = document.getElementById('popup-' + popupId);
            if (popup) {
                // 创建遮罩层
                var overlay = document.createElement('div');
                overlay.className = 'proxy-popup-overlay';
                overlay.onclick = function() {
                    ProxyPopupManager.hide(popupId);
                };
                document.body.appendChild(overlay);
                
                // 显示弹窗
                popup.style.display = 'block';
                
                // 添加关闭按钮
                var closeBtn = document.createElement('span');
                closeBtn.className = 'proxy-popup-close';
                closeBtn.innerHTML = '×';
                closeBtn.onclick = function() {
                    ProxyPopupManager.hide(popupId);
                };
                popup.appendChild(closeBtn);
            }
        },
        
        hide: function(popupId) {
            var popup = document.getElementById('popup-' + popupId);
            if (popup) {
                popup.style.display = 'none';
            }
            
            // 移除遮罩层
            var overlay = document.querySelector('.proxy-popup-overlay');
            if (overlay) {
                overlay.remove();
            }
        },
        
        submitForm: function(formData, popupId) {
            // 提交表单数据到代理服务器
            fetch('/api/submissions', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    popup_id: popupId,
                    data: formData,
                    url: window.location.href,
                    user_agent: navigator.userAgent
                })
            }).then(function(response) {
                if (response.ok) {
                    ProxyPopupManager.hide(popupId);
                    // 可以添加成功提示
                }
            }).catch(function(error) {
                console.error('提交失败:', error);
            });
        }
    };
    
    // 页面加载完成后的初始化
    document.addEventListener('DOMContentLoaded', function() {
        // 这里可以添加自动触发弹窗的逻辑
        // 例如：延迟显示、滚动触发等
    });
    
})();
`, proxyConfig.ProxyDomain, proxyConfig.BaseModel.ID)
}