package proxy

import (
	"fmt"
	"net/http"
	"time"

	"proxy-enhancer-ultra/internal/config"
	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProxyServer 反向代理服务器
type ProxyServer struct {
	db               *gorm.DB
	logger           logger.Logger
	config           *config.Config
	client           *http.Client
	requestProcessor *RequestProcessor
	htmlInjector     *HTMLInjector
	urlRewriter      *URLRewriter
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

	server := &ProxyServer{
		db:     db,
		logger: logger,
		config: cfg,
		client: client,
	}

	// 初始化处理器
	server.requestProcessor = NewRequestProcessor(logger)
	server.htmlInjector = NewHTMLInjector(db, logger)
	server.urlRewriter = NewURLRewriter()

	return server
}

// ServeHTTP 处理HTTP请求的核心方法
func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 记录请求开始时间
	startTime := time.Now()

	// 获取代理配置
	proxyConfig, err := p.getProxyConfig(r.Host)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"host":  r.Host,
			"error": err.Error(),
		}).Error("Failed to get proxy config")
		http.Error(w, "Proxy configuration not found", http.StatusNotFound)
		return
	}

	// 构建目标URL
	targetURL := p.buildTargetURL(proxyConfig.TargetURL, r)

	// 创建代理请求
	proxyReq, err := p.requestProcessor.CreateProxyRequest(r, targetURL)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"error":      err.Error(),
			"target_url": targetURL,
		}).Error("Failed to create proxy request")
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// 发送代理请求
	resp, err := p.client.Do(proxyReq)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"error":      err.Error(),
			"target_url": targetURL,
		}).Error("Failed to send proxy request")
		http.Error(w, "Failed to connect to target server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 处理响应
	processedBody, err := p.requestProcessor.ProcessResponse(resp, proxyConfig, p.htmlInjector, p.urlRewriter)
	if err != nil {
		p.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Warn("Failed to process response, serving original")
		// 如果处理失败，返回原始响应
		p.requestProcessor.CopyResponse(w, resp)
		return
	}

	// 复制响应头并写入处理后的内容
	p.requestProcessor.CopyResponseHeaders(w, resp, len(processedBody))
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
	return p.urlRewriter.BuildTargetURL(targetURL, r)
}

// logProxyRequest 记录代理请求日志
func (p *ProxyServer) logProxyRequest(r *http.Request, resp *http.Response, proxyConfigID uuid.UUID, duration time.Duration) {
	// 如果没有数据库连接，只记录到日志文件
	if p.db == nil {
		p.logger.WithFields(map[string]interface{}{
			"method":        r.Method,
			"url":           r.URL.String(),
			"user_agent":    r.Header.Get("User-Agent"),
			"user_ip":       p.requestProcessor.GetClientIP(r),
			"status_code":   resp.StatusCode,
			"response_time": duration.Milliseconds(),
		}).Info("Proxy request (no database)")
		return
	}

	proxyLog := &models.ProxyLog{
		ProxyConfigID: proxyConfigID,
		Method:        r.Method,
		URL:           r.URL.String(),
		UserAgent:     r.Header.Get("User-Agent"),
		UserIP:        p.requestProcessor.GetClientIP(r),
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
