package utils

import (
	"strconv"
	"strings"
)

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
