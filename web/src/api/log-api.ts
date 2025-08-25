import api from './index'
import type { ApiResponse, PageResponse, PageParams } from './index'
import type { OperationLog } from './system-types'

// 操作日志API
export const getOperationLogs = (params: PageParams & {
  user_id?: string
  module?: string
  action?: string
  status?: string
  start_time?: string
  end_time?: string
}) => {
  return api.get<ApiResponse<PageResponse<OperationLog>>>('/system/logs', { params })
}

export const getOperationLogDetail = (id: string) => {
  return api.get<ApiResponse<OperationLog>>(`/system/logs/${id}`)
}

export const clearOperationLogs = (beforeDate?: string) => {
  return api.delete<ApiResponse<void>>('/system/logs/clear', {
    data: { beforeDate }
  })
}

// 访问日志
export const getAccessLogs = (params: PageParams & {
  user_id?: string
  ip_address?: string
  method?: string
  status?: number
  start_time?: string
  end_time?: string
}) => {
  return api.get<ApiResponse<PageResponse<any>>>('/system/logs/access', { params })
}

// 错误日志
export const getErrorLogs = (params: PageParams & {
  level?: string
  module?: string
  start_time?: string
  end_time?: string
}) => {
  return api.get<ApiResponse<PageResponse<any>>>('/system/logs/error', { params })
}

// 系统日志
export const getSystemLogs = (params: PageParams & {
  level?: string
  component?: string
  start_time?: string
  end_time?: string
}) => {
  return api.get<ApiResponse<PageResponse<any>>>('/system/logs/system', { params })
}

// 登录日志
export const getLoginLogs = (params: PageParams & {
  user_id?: string
  ip_address?: string
  status?: 'success' | 'failed'
  start_time?: string
  end_time?: string
}) => {
  return api.get<ApiResponse<PageResponse<any>>>('/system/logs/login', { params })
}

// 日志统计
export const getLogStats = (timeRange?: 'today' | 'week' | 'month' | 'year') => {
  return api.get<ApiResponse<{
    operation_logs: number
    access_logs: number
    error_logs: number
    login_logs: number
    failed_logins: number
    unique_visitors: number
    top_actions: Array<{
      action: string
      count: number
    }>
    top_users: Array<{
      username: string
      count: number
    }>
    hourly_distribution: Array<{
      hour: number
      count: number
    }>
    daily_trend: Array<{
      date: string
      count: number
    }>
  }>>('/system/logs/stats', {
    params: { timeRange }
  })
}

// 日志导出
export const exportOperationLogs = (params?: {
  user_id?: string
  module?: string
  action?: string
  status?: string
  start_time?: string
  end_time?: string
  format?: 'csv' | 'excel'
}) => {
  return api.get<Blob>('/system/logs/export', {
    params,
    responseType: 'blob'
  })
}

export const exportAccessLogs = (params?: {
  user_id?: string
  ip_address?: string
  start_time?: string
  end_time?: string
  format?: 'csv' | 'excel'
}) => {
  return api.get<Blob>('/system/logs/access/export', {
    params,
    responseType: 'blob'
  })
}

// 日志分析
export const getLogAnalysis = (type: 'operation' | 'access' | 'error', timeRange: string) => {
  return api.get<ApiResponse<{
    total_count: number
    period_comparison: {
      current: number
      previous: number
      change_percent: number
    }
    top_ips: Array<{
      ip: string
      count: number
      location?: string
    }>
    top_users: Array<{
      username: string
      count: number
    }>
    patterns: Array<{
      pattern: string
      count: number
      risk_level: 'low' | 'medium' | 'high'
    }>
  }>>(`/system/logs/analysis/${type}`, {
    params: { timeRange }
  })
}

// 实时日志
export const getRealtimeLogs = (type: 'operation' | 'access' | 'error', limit: number = 50) => {
  return api.get<ApiResponse<any[]>>(`/system/logs/realtime/${type}`, {
    params: { limit }
  })
}

// 日志搜索
export const searchLogs = (query: string, type?: string, params?: PageParams) => {
  return api.post<ApiResponse<PageResponse<any>>>('/system/logs/search', {
    query,
    type,
    ...params
  })
}

// 日志配置
export const getLogConfig = () => {
  return api.get<ApiResponse<{
    retention_days: number
    max_file_size: number
    compression_enabled: boolean
    remote_storage_enabled: boolean
    alert_thresholds: {
      error_count: number
      failed_login_count: number
    }
  }>>('/system/logs/config')
}

export const updateLogConfig = (config: any) => {
  return api.patch<ApiResponse<void>>('/system/logs/config', config)
}

// 日志清理
export const cleanupLogs = (type: string, beforeDate: string) => {
  return api.delete<ApiResponse<{
    deleted_count: number
    freed_space: number
  }>>(`/system/logs/cleanup/${type}`, {
    data: { beforeDate }
  })
}

// 日志API对象导出
export const logApi = {
  getOperationLogs,
  getOperationLogDetail,
  clearOperationLogs,
  getAccessLogs,
  getErrorLogs,
  getSystemLogs,
  getLoginLogs,
  getLogStats,
  exportOperationLogs,
  exportAccessLogs,
  getLogAnalysis,
  getRealtimeLogs,
  searchLogs,
  getLogConfig,
  updateLogConfig,
  cleanupLogs
}
