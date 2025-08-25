package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// SanitizeHTML 清理HTML内容（简单版本）
func SanitizeHTML(html string) string {
	// 移除script标签
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	html = scriptRegex.ReplaceAllString(html, "")

	// 移除on*事件属性
	eventRegex := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
	html = eventRegex.ReplaceAllString(html, "")

	// 移除javascript:协议
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	html = jsRegex.ReplaceAllString(html, "")

	return html
}

// RewriteURLsInHTML 重写HTML中的URL
func RewriteURLsInHTML(html, baseURL, proxyDomain string) string {
	// 重写href属性
	hrefRegex := regexp.MustCompile(`href\s*=\s*["']([^"']*)["']`)
	html = hrefRegex.ReplaceAllStringFunc(html, func(match string) string {
		return rewriteURLAttribute(match, "href", baseURL, proxyDomain)
	})

	// 重写src属性
	srcRegex := regexp.MustCompile(`src\s*=\s*["']([^"']*)["']`)
	html = srcRegex.ReplaceAllStringFunc(html, func(match string) string {
		return rewriteURLAttribute(match, "src", baseURL, proxyDomain)
	})

	// 重写action属性
	actionRegex := regexp.MustCompile(`action\s*=\s*["']([^"']*)["']`)
	html = actionRegex.ReplaceAllStringFunc(html, func(match string) string {
		return rewriteURLAttribute(match, "action", baseURL, proxyDomain)
	})

	return html
}

// rewriteURLAttribute 重写URL属性
func rewriteURLAttribute(match, attr, baseURL, proxyDomain string) string {
	// 提取URL
	regex := regexp.MustCompile(attr + `\s*=\s*["']([^"']*)["']`)
	matches := regex.FindStringSubmatch(match)
	if len(matches) < 2 {
		return match
	}

	originalURL := matches[1]
	rewrittenURL := rewriteURL(originalURL, baseURL, proxyDomain)

	return strings.Replace(match, originalURL, rewrittenURL, 1)
}

// rewriteURL 重写单个URL
func rewriteURL(originalURL, baseURL, proxyDomain string) string {
	// 如果是绝对URL且不是目标域名，需要重写
	if strings.HasPrefix(originalURL, "http://") || strings.HasPrefix(originalURL, "https://") {
		u, err := url.Parse(originalURL)
		if err != nil {
			return originalURL
		}

		// 如果是目标网站的URL，重写为代理URL
		baseU, err := url.Parse(baseURL)
		if err != nil {
			return originalURL
		}

		if u.Host == baseU.Host {
			return fmt.Sprintf("//%s%s", proxyDomain, u.Path)
		}
	}

	// 如果是相对URL，保持不变
	return originalURL
}

// InjectScript 在HTML中注入脚本
func InjectScript(html, script string) string {
	// 在</head>前注入
	headEndRegex := regexp.MustCompile(`(?i)</head>`)
	if headEndRegex.MatchString(html) {
		return headEndRegex.ReplaceAllString(html, fmt.Sprintf("<script>%s</script></head>", script))
	}

	// 如果没有</head>，在</body>前注入
	bodyEndRegex := regexp.MustCompile(`(?i)</body>`)
	if bodyEndRegex.MatchString(html) {
		return bodyEndRegex.ReplaceAllString(html, fmt.Sprintf("<script>%s</script></body>", script))
	}

	// 如果都没有，直接在末尾添加
	return html + fmt.Sprintf("<script>%s</script>", script)
}

// InjectHTML 在HTML中注入HTML内容
func InjectHTML(html, content string) string {
	// 在</body>前注入
	bodyEndRegex := regexp.MustCompile(`(?i)</body>`)
	if bodyEndRegex.MatchString(html) {
		return bodyEndRegex.ReplaceAllString(html, content+"</body>")
	}

	// 如果没有</body>，直接在末尾添加
	return html + content
}
