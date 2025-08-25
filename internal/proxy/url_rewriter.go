package proxy

import (
	"net/http"
	"net/url"
	"strings"

	"proxy-enhancer-ultra/internal/models"

	"golang.org/x/net/html"
)

// URLRewriter URL重写器
type URLRewriter struct {
	htmlParser *HTMLParser
}

// NewURLRewriter 创建新的URL重写器
func NewURLRewriter() *URLRewriter {
	return &URLRewriter{
		htmlParser: NewHTMLParser(),
	}
}

// BuildTargetURL 构建目标URL
func (ur *URLRewriter) BuildTargetURL(targetURL string, r *http.Request) string {
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

// RewriteURLs 重写HTML中的URL
func (ur *URLRewriter) RewriteURLs(doc *html.Node, proxyConfig *models.ProxyConfig) {
	// 重写所有链接和资源引用
	ur.htmlParser.WalkNodes(doc, func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "a", "link":
				ur.rewriteAttribute(node, "href", proxyConfig)
			case "img", "script":
				ur.rewriteAttribute(node, "src", proxyConfig)
			case "form":
				ur.rewriteAttribute(node, "action", proxyConfig)
			}
		}
	})
}

// rewriteAttribute 重写属性
func (ur *URLRewriter) rewriteAttribute(node *html.Node, attrName string, proxyConfig *models.ProxyConfig) {
	for i, attr := range node.Attr {
		if attr.Key == attrName {
			node.Attr[i].Val = ur.rewriteURL(attr.Val, proxyConfig)
			break
		}
	}
}

// rewriteURL 重写单个URL
func (ur *URLRewriter) rewriteURL(originalURL string, proxyConfig *models.ProxyConfig) string {
	// 如果是相对URL，转换为绝对URL
	if strings.HasPrefix(originalURL, "/") {
		return "https://" + proxyConfig.ProxyDomain + originalURL
	}

	// 如果是目标域名的URL，替换为代理域名
	targetHost := ur.extractHost(proxyConfig.TargetURL)
	if strings.Contains(originalURL, targetHost) {
		return strings.Replace(originalURL, targetHost, proxyConfig.ProxyDomain, 1)
	}

	return originalURL
}

// extractHost 从URL中提取主机名
func (ur *URLRewriter) extractHost(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}
