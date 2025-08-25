import api from './index'
import type { ApiResponse } from './index'

// 仪表盘统计数据
export interface DashboardStats {
  proxy: {
    total: number
    active: number
    inactive: number
    error: number
  }
  rules: {
    total: number
    enabled: number
    disabled: number
    todayExecutions: number
  }
  data: {
    totalPopups: number
    activePopups: number
    totalSubmissions: number
    conversionRate: number
  }
  system: {
    uptime: string
    memoryUsage: number
    diskUsage: number
    activeUsers: number
  }
}

// 流量统计数据
export interface TrafficStats {
  total: {
    requests: number
    bytes: number
    avgResponseTime: number
  }
  today: {
    requests: number
    bytes: number
    avgResponseTime: number
  }
  hourly: Array<{
    hour: string
    requests: number
    bytes: number
    responseTime: number
  }>
  daily: Array<{
    date: string
    requests: number
    bytes: number
    responseTime: number
  }>
}

// 代理状态分布
export interface ProxyStatusDistribution {
  active: number
  inactive: number
  error: number
  testing: number
}

// 规则执行统计
export interface RuleExecutionStats {
  total: number
  success: number
  failed: number
  successRate: number
  topRules: Array<{
    id: number
    name: string
    executions: number
    successRate: number
  }>
  executionTrend: Array<{
    time: string
    executions: number
    success: number
    failed: number
  }>
}

// 系统性能数据
export interface SystemPerformance {
  cpu: {
    usage: number
    cores: number
    loadAverage: number[]
  }
  memory: {
    used: number
    total: number
    percentage: number
    available: number
  }
  disk: {
    used: number
    total: number
    percentage: number
    available: number
  }
  network: {
    bytesIn: number
    bytesOut: number
    packetsIn: number
    packetsOut: number
  }
}

// 最近活动
export interface RecentActivity {
  id: string
  type: 'proxy' | 'rule' | 'data' | 'system' | 'user'
  action: string
  description: string
  user?: string
  timestamp: string
  status: 'success' | 'warning' | 'error'
  details?: Record<string, any>
}

// 系统状态
export interface SystemStatus {
  overall: 'healthy' | 'warning' | 'error'
  services: Array<{
    name: string
    status: 'running' | 'stopped' | 'error'
    uptime?: string
    lastCheck: string
    responseTime?: number
  }>
  alerts: Array<{
    id: string
    level: 'info' | 'warning' | 'error'
    message: string
    timestamp: string
    resolved: boolean
  }>
}

// 获取仪表盘统计数据
export const getDashboardStats = () => {
  return api.get<ApiResponse<DashboardStats>>('/dashboard/stats')
}

// 获取流量统计数据
export const getTrafficStats = (params: {
  period?: 'hour' | 'day' | 'week' | 'month'
  startTime?: string
  endTime?: string
}) => {
  return api.get<ApiResponse<TrafficStats>>('/dashboard/traffic', { params })
}

// 获取代理状态分布
export const getProxyStatusDistribution = () => {
  return api.get<ApiResponse<ProxyStatusDistribution>>('/dashboard/proxy-status')
}

// 获取规则执行统计
export const getRuleExecutionStats = (params: {
  period?: 'hour' | 'day' | 'week' | 'month'
  startTime?: string
  endTime?: string
}) => {
  return api.get<ApiResponse<RuleExecutionStats>>('/dashboard/rule-stats', { params })
}

// 获取系统性能数据
export const getSystemPerformance = () => {
  return api.get<ApiResponse<SystemPerformance>>('/dashboard/performance')
}

// 获取最近活动
export const getRecentActivities = (params: {
  limit?: number
  type?: string
  startTime?: string
  endTime?: string
}) => {
  return api.get<ApiResponse<RecentActivity[]>>('/dashboard/activities', { params })
}

// 获取系统状态
export const getSystemStatus = () => {
  return api.get<ApiResponse<SystemStatus>>('/dashboard/status')
}

// 获取实时数据（WebSocket连接信息）
export const getRealtimeConfig = () => {
  return api.get<ApiResponse<{
    wsUrl: string
    channels: string[]
    reconnectInterval: number
  }>>('/dashboard/realtime-config')
}

// 获取快速操作配置
export const getQuickActions = () => {
  return api.get<ApiResponse<Array<{
    id: string
    name: string
    icon: string
    action: string
    params?: Record<string, any>
    permission?: string
  }>>>('/dashboard/quick-actions')
}

// 执行快速操作
export const executeQuickAction = (actionId: string, params?: Record<string, any>) => {
  return api.post<ApiResponse<{ message: string; result?: any }>>('/dashboard/quick-actions/execute', {
    actionId,
    params
  })
}

// 获取图表配置
export const getChartConfigs = () => {
  return api.get<ApiResponse<Record<string, {
    type: string
    options: Record<string, any>
    refreshInterval?: number
  }>>>('/dashboard/chart-configs')
}

// 更新图表配置
export const updateChartConfig = (chartId: string, config: Record<string, any>) => {
  return api.patch<ApiResponse<void>>(`/dashboard/chart-configs/${chartId}`, config)
}

// 获取用户偏好设置
export const getUserPreferences = () => {
  return api.get<ApiResponse<{
    theme: 'light' | 'dark'
    language: string
    timezone: string
    dashboardLayout: Record<string, any>
    notifications: {
      email: boolean
      browser: boolean
      sound: boolean
    }
  }>>('/dashboard/preferences')
}

// 更新用户偏好设置
export const updateUserPreferences = (preferences: Record<string, any>) => {
  return api.patch<ApiResponse<void>>('/dashboard/preferences', preferences)
}