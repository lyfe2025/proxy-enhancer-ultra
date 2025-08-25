/**
 * 代理管理相关工具函数
 */

/**
 * 获取代理类型对应的标签类型
 * @param type 代理类型
 * @returns Element Plus 标签类型
 */
export const getTypeTagType = (type: string): 'success' | 'warning' | 'info' => {
  const typeMap: Record<string, 'success' | 'warning' | 'info'> = {
    http: 'info',
    https: 'success',
    socks5: 'warning',
    socks4: 'info'
  }
  return typeMap[type] || 'info'
}

/**
 * 获取代理状态对应的标签类型
 * @param status 代理状态
 * @returns Element Plus 标签类型
 */
export const getStatusTagType = (status: string): 'success' | 'info' | 'danger' => {
  const statusMap: Record<string, 'success' | 'info' | 'danger'> = {
    active: 'success',
    inactive: 'info',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

/**
 * 获取代理状态的中文文本
 * @param status 代理状态
 * @returns 中文状态文本
 */
export const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    active: '活跃',
    inactive: '停用',
    error: '异常'
  }
  return statusMap[status] || status
}

/**
 * 格式化字节数为可读的文件大小
 * @param bytes 字节数
 * @returns 格式化后的文件大小字符串
 */
export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 格式化日期为本地化字符串
 * @param date 日期字符串
 * @returns 格式化后的日期字符串
 */
export const formatDate = (date: string): string => {
  return new Date(date).toLocaleString('zh-CN')
}