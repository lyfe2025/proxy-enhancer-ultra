package utils

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

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
