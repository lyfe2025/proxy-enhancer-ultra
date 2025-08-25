import api from './index'
import type { ApiResponse } from './index'
import type { SystemStats, SystemHealth } from './system-types'

// 系统统计API
export const getSystemStats = () => {
  return api.get<ApiResponse<SystemStats>>('/system/stats')
}

// 系统监控API
export const getSystemHealth = () => {
  return api.get<ApiResponse<SystemHealth>>('/system/health')
}

// 性能监控
export const getPerformanceMetrics = (timeRange?: 'hour' | 'day' | 'week' | 'month') => {
  return api.get<ApiResponse<{
    cpu_usage: Array<{ timestamp: string; value: number }>
    memory_usage: Array<{ timestamp: string; value: number }>
    disk_usage: Array<{ timestamp: string; value: number }>
    network_io: Array<{ timestamp: string; in: number; out: number }>
    database_connections: Array<{ timestamp: string; active: number; idle: number }>
    response_times: Array<{ timestamp: string; avg: number; p95: number; p99: number }>
  }>>('/system/metrics/performance', {
    params: { timeRange }
  })
}

// 业务统计
export const getBusinessStats = () => {
  return api.get<ApiResponse<{
    users: {
      total: number
      active: number
      new_today: number
      new_week: number
      new_month: number
    }
    content: {
      total_posts: number
      total_comments: number
      today_posts: number
      today_comments: number
    }
    system: {
      total_requests: number
      api_calls: number
      page_views: number
      error_rate: number
    }
  }>>('/system/stats/business')
}

// 用户活跃度统计
export const getUserActivityStats = (timeRange?: string) => {
  return api.get<ApiResponse<{
    online_users: number
    daily_active_users: number
    weekly_active_users: number
    monthly_active_users: number
    user_growth: Array<{
      date: string
      new_users: number
      total_users: number
    }>
    activity_heatmap: Array<{
      hour: number
      day: number
      count: number
    }>
    top_pages: Array<{
      path: string
      views: number
      unique_visitors: number
    }>
  }>>('/system/stats/user-activity', {
    params: { timeRange }
  })
}

// 系统负载统计
export const getSystemLoadStats = () => {
  return api.get<ApiResponse<{
    current_load: {
      cpu: number
      memory: number
      disk: number
      network: number
    }
    peak_times: Array<{
      time: string
      cpu: number
      memory: number
      requests: number
    }>
    resource_usage_trend: Array<{
      timestamp: string
      cpu: number
      memory: number
      disk: number
    }>
    alerts: Array<{
      type: string
      level: 'warning' | 'critical'
      message: string
      timestamp: string
    }>
  }>>('/system/stats/load')
}

// 错误统计
export const getErrorStats = (timeRange?: string) => {
  return api.get<ApiResponse<{
    total_errors: number
    error_rate: number
    error_trend: Array<{
      timestamp: string
      count: number
    }>
    error_types: Array<{
      type: string
      count: number
      percentage: number
    }>
    top_errors: Array<{
      message: string
      count: number
      last_occurrence: string
    }>
    error_distribution: Array<{
      status_code: number
      count: number
      percentage: number
    }>
  }>>('/system/stats/errors', {
    params: { timeRange }
  })
}

// API使用统计
export const getApiStats = (timeRange?: string) => {
  return api.get<ApiResponse<{
    total_requests: number
    requests_per_second: number
    avg_response_time: number
    success_rate: number
    top_endpoints: Array<{
      path: string
      method: string
      count: number
      avg_response_time: number
    }>
    status_distribution: Array<{
      status: number
      count: number
      percentage: number
    }>
    hourly_distribution: Array<{
      hour: number
      count: number
    }>
  }>>('/system/stats/api', {
    params: { timeRange }
  })
}

// 数据库统计
export const getDatabaseStats = () => {
  return api.get<ApiResponse<{
    connection_pool: {
      active: number
      idle: number
      max: number
    }
    query_performance: {
      avg_query_time: number
      slow_queries: number
      total_queries: number
    }
    table_sizes: Array<{
      table: string
      rows: number
      size: string
    }>
    index_usage: Array<{
      table: string
      index: string
      usage_count: number
    }>
  }>>('/system/stats/database')
}

// 缓存统计
export const getCacheStats = () => {
  return api.get<ApiResponse<{
    hit_rate: number
    miss_rate: number
    total_requests: number
    memory_usage: number
    key_count: number
    expired_keys: number
    top_keys: Array<{
      key: string
      hits: number
      size: number
    }>
  }>>('/system/stats/cache')
}

// 实时统计
export const getRealtimeStats = () => {
  return api.get<ApiResponse<{
    online_users: number
    cpu_usage: number
    memory_usage: number
    active_requests: number
    recent_activities: Array<{
      user: string
      action: string
      timestamp: string
    }>
  }>>('/system/stats/realtime')
}

// 导出统计报表
export const exportStatsReport = (type: string, timeRange: string, format: 'pdf' | 'excel' | 'csv') => {
  return api.get<Blob>(`/system/stats/export/${type}`, {
    params: { timeRange, format },
    responseType: 'blob'
  })
}

// 统计API对象导出
export const statsApi = {
  getSystemStats,
  getSystemHealth,
  getPerformanceMetrics,
  getBusinessStats,
  getUserActivityStats,
  getSystemLoadStats,
  getErrorStats,
  getApiStats,
  getDatabaseStats,
  getCacheStats,
  getRealtimeStats,
  exportStatsReport
}
