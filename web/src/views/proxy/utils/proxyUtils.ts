/**
 * 代理管理相关工具函数
 */

/**
 * 获取代理类型对应的标签类型
 * @param type 代理类型
 * @returns Element Plus 标签类型
 */
export const getTypeTagType = (type: string): "success" | "info" | "warning" | "danger" | "primary" => {
  const typeMap: Record<string, "success" | "info" | "warning" | "danger" | "primary"> = {
    http: 'primary',
    https: 'success',
    socks5: 'warning'
  }
  return typeMap[type] || 'info'
}

/**
 * 获取代理状态对应的标签类型
 * @param status 代理状态
 * @returns Element Plus 标签类型
 */
export const getStatusTagType = (status: string): "success" | "info" | "warning" | "danger" | "primary" => {
  const statusMap: Record<string, "success" | "info" | "warning" | "danger" | "primary"> = {
    active: 'success',
    inactive: 'info',
    error: 'danger',
    testing: 'warning'
  }
  return statusMap[status] || 'info'
}

/**
 * 获取代理状态对应的文本
 * @param status 代理状态
 * @returns 状态文本
 */
export const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    active: '活跃',
    inactive: '未激活',
    error: '错误',
    testing: '测试中'
  }
  return statusMap[status] || '未知'
}

/**
 * 格式化字节数
 * @param bytes 字节数
 * @returns 格式化后的字符串
 */
export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 格式化日期
 * @param date 日期字符串或Date对象
 * @returns 格式化后的日期字符串
 */
export const formatDate = (date: string | Date): string => {
  if (!date) return '-'
  
  const d = new Date(date)
  if (isNaN(d.getTime())) return '-'
  
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

/**
 * 验证IP地址格式
 * @param ip IP地址字符串
 * @returns 是否为有效IP地址
 */
export const isValidIP = (ip: string): boolean => {
  const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  return ipRegex.test(ip)
}

/**
 * 验证域名格式
 * @param domain 域名字符串
 * @returns 是否为有效域名
 */
export const isValidDomain = (domain: string): boolean => {
  const domainRegex = /^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)?$/
  return domainRegex.test(domain)
}

/**
 * 验证端口号
 * @param port 端口号
 * @returns 是否为有效端口号
 */
export const isValidPort = (port: number): boolean => {
  return Number.isInteger(port) && port >= 1 && port <= 65535
}