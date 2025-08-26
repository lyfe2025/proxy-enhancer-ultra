import api from './index'
import type { ApiResponse } from './index'
import type { SystemStats, SystemHealth } from './system-types'

// 获取系统统计数据
export const getSystemStats = (params?: {
  start_date?: string
  end_date?: string
  granularity?: 'hour' | 'day' | 'week' | 'month'
}) => api.get<SystemStats>('/monitoring/dashboard', { params })

// 获取健康检查状态
export const getHealthStatus = () => 
  api.get<HealthStatus>('/monitoring/health')

// 获取性能指标
export const getPerformanceMetrics = (params?: {
  start_date?: string
  end_date?: string
  metric_type?: 'cpu' | 'memory' | 'disk' | 'network'
}) => api.get<PerformanceMetrics>('/monitoring/metrics/system', { params })

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
  }>>('/monitoring/stats/business')
}

// 获取用户活动统计
export const getUserActivityStats = (params?: {
  start_date?: string
  end_date?: string
  user_id?: string
}) => api.get<UserActivityStats>('/monitoring/dashboard', { params })

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
  }>>('/monitoring/stats/load')
}

// 获取错误统计
export const getErrorStats = (params?: {
  start_date?: string
  end_date?: string
  error_type?: string
}) => api.get<ErrorStats>('/monitoring/dashboard', { params })

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
  }>>('/monitoring/stats/api', {
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
  }>>('/monitoring/stats/database')
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
  }>>('/monitoring/stats/cache')
}

// 获取实时统计
export const getRealtimeStats = () => 
  api.get<RealtimeStats>('/monitoring/dashboard')

// 导出统计报告
export const exportStatsReport = (params: {
  start_date: string
  end_date: string
  report_type: 'summary' | 'detailed' | 'custom'
  format: 'pdf' | 'excel' | 'csv'
  include_charts?: boolean
}) => api.get<Blob>('/monitoring/dashboard', { 
  params,
  responseType: 'blob'
})

// 获取代理统计
export const getProxyStats = (params?: {
  start_date?: string
  end_date?: string
  proxy_id?: string
}) => api.get<ProxyStats>('/monitoring/metrics/proxy', { params })

// 获取系统信息
export const getSystemInfo = () => 
  api.get<SystemInfo>('/monitoring/health')

// 获取服务状态
export const getServiceStatus = () => 
  api.get<ServiceStatus[]>('/monitoring/health')

// 获取趋势分析
export const getTrendAnalysis = (params: {
  metric: string
  start_date: string
  end_date: string
  prediction_days?: number
}) => api.get<TrendAnalysis>('/monitoring/dashboard', { params })

// 获取对比分析
export const getComparisonAnalysis = (params: {
  metric: string
  current_period: { start: string; end: string }
  comparison_period: { start: string; end: string }
}) => api.get<ComparisonAnalysis>('/monitoring/dashboard', { params })

// 统计API对象导出
export const statsApi = {
  getSystemStats,
  getHealthStatus,
  getPerformanceMetrics,
  getUserActivityStats,
  getErrorStats,
  getRealtimeStats,
  exportStatsReport,
  getProxyStats,
  getSystemInfo,
  getServiceStatus,
  getTrendAnalysis,
  getComparisonAnalysis
}
