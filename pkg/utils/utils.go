package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码哈希
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// MD5Hash 计算MD5哈希
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SHA256Hash 计算SHA256哈希
func SHA256Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GetClientIP 获取客户端IP地址
func GetClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// 检查X-Real-IP头
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 检查X-Forwarded头
	if xf := r.Header.Get("X-Forwarded"); xf != "" {
		return xf
	}

	// 使用RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidURL 验证URL格式
func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsValidDomain 验证域名格式
func IsValidDomain(domain string) bool {
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?))*$`)
	return domainRegex.MatchString(domain)
}

// NormalizeURL 标准化URL
func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// 如果没有协议，默认使用http
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// 标准化主机名（转为小写）
	u.Host = strings.ToLower(u.Host)

	// 移除默认端口
	if (u.Scheme == "http" && strings.HasSuffix(u.Host, ":80")) ||
		(u.Scheme == "https" && strings.HasSuffix(u.Host, ":443")) {
		u.Host = u.Host[:strings.LastIndex(u.Host, ":")]
	}

	return u.String(), nil
}

// ExtractDomain 从URL中提取域名
func ExtractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

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

// ParseUserAgent 解析User-Agent
func ParseUserAgent(userAgent string) map[string]string {
	result := make(map[string]string)
	result["raw"] = userAgent

	// 简单的User-Agent解析
	if strings.Contains(userAgent, "Chrome") {
		result["browser"] = "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		result["browser"] = "Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		result["browser"] = "Safari"
	} else if strings.Contains(userAgent, "Edge") {
		result["browser"] = "Edge"
	} else {
		result["browser"] = "Unknown"
	}

	if strings.Contains(userAgent, "Windows") {
		result["os"] = "Windows"
	} else if strings.Contains(userAgent, "Mac") {
		result["os"] = "macOS"
	} else if strings.Contains(userAgent, "Linux") {
		result["os"] = "Linux"
	} else if strings.Contains(userAgent, "Android") {
		result["os"] = "Android"
	} else if strings.Contains(userAgent, "iOS") {
		result["os"] = "iOS"
	} else {
		result["os"] = "Unknown"
	}

	if strings.Contains(userAgent, "Mobile") {
		result["device"] = "Mobile"
	} else if strings.Contains(userAgent, "Tablet") {
		result["device"] = "Tablet"
	} else {
		result["device"] = "Desktop"
	}

	return result
}

// FormatBytes 格式化字节数
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration 格式化时间间隔
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// StringToInt 字符串转整数
func StringToInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

// StringToUint 字符串转无符号整数
func StringToUint(s string, defaultValue uint) uint {
	if i, err := strconv.ParseUint(s, 10, 32); err == nil {
		return uint(i)
	}
	return defaultValue
}

// StringToBool 字符串转布尔值
func StringToBool(s string, defaultValue bool) bool {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return defaultValue
}

// Contains 检查切片是否包含元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates 移除切片中的重复元素
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// MergeStringMaps 合并字符串映射
func MergeStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// GetMapKeys 获取映射的所有键
func GetMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// IsLocalhost 检查是否为本地地址
func IsLocalhost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

// IsPrivateIP 检查是否为私有IP
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// 检查IPv4私有地址
	if parsedIP.To4() != nil {
		return parsedIP.IsLoopback() ||
			parsedIP.IsLinkLocalUnicast() ||
			parsedIP.IsPrivate()
	}

	// 检查IPv6私有地址
	return parsedIP.IsLoopback() ||
		parsedIP.IsLinkLocalUnicast() ||
		parsedIP.IsPrivate()
}