import api from './index'
import type { ApiResponse, PageResponse, PageParams } from './index'

// 代理配置类型
export interface ProxyConfig {
  id?: number
  name: string
  type: 'http' | 'https' | 'socks5'
  host: string
  port: number
  username?: string
  password?: string
  status: 'active' | 'inactive' | 'error'
  description?: string
  createdAt?: string
  updatedAt?: string
  lastUsed?: string
  responseTime?: number
  successRate?: number
  totalRequests?: number
  failedRequests?: number
}

// 代理统计信息
export interface ProxyStats {
  total: number
  active: number
  inactive: number
  error: number
  totalTraffic: string
  avgResponseTime: number
  successRate: number
}

// 代理测试结果
export interface ProxyTestResult {
  id: number
  status: 'success' | 'failed'
  responseTime?: number
  error?: string
  testTime: string
}

// 获取代理列表
export const getProxyList = (params: PageParams & {
  type?: string
  status?: string
}) => {
  return api.get<ApiResponse<PageResponse<ProxyConfig>>>('/proxy/list', { params })
}

// 获取代理详情
export const getProxyDetail = (id: number) => {
  return api.get<ApiResponse<ProxyConfig>>(`/proxy/${id}`)
}

// 创建代理
export const createProxy = (data: Omit<ProxyConfig, 'id' | 'createdAt' | 'updatedAt'>) => {
  return api.post<ApiResponse<ProxyConfig>>('/proxy', data)
}

// 更新代理
export const updateProxy = (id: number, data: Partial<ProxyConfig>) => {
  return api.put<ApiResponse<ProxyConfig>>(`/proxy/${id}`, data)
}

// 删除代理
export const deleteProxy = (id: number) => {
  return api.delete<ApiResponse<void>>(`/proxy/${id}`)
}

// 批量删除代理
export const batchDeleteProxy = (ids: number[]) => {
  return api.delete<ApiResponse<void>>('/proxy/batch', { data: { ids } })
}

// 测试代理连接
export const testProxy = (id: number) => {
  return api.post<ApiResponse<ProxyTestResult>>(`/proxy/${id}/test`)
}

// 批量测试代理
export const batchTestProxy = (ids: number[]) => {
  return api.post<ApiResponse<ProxyTestResult[]>>('/proxy/batch-test', { ids })
}

// 获取代理统计信息
export const getProxyStats = () => {
  return api.get<ApiResponse<ProxyStats>>('/proxy/stats')
}

// 启用/禁用代理
export const toggleProxyStatus = (id: number, status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>(`/proxy/${id}/status`, { status })
}

// 批量启用/禁用代理
export const batchToggleProxyStatus = (ids: number[], status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>('/proxy/batch-status', { ids, status })
}

// 导入代理配置
export const importProxyConfig = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post<ApiResponse<{ success: number; failed: number }>>('/proxy/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 导出代理配置
export const exportProxyConfig = (ids?: number[]) => {
  return api.post<Blob>('/proxy/export', { ids }, {
    responseType: 'blob'
  })
}

// 创建 proxyApi 对象，包含所有相关函数
export const proxyApi = {
  getProxies: getProxyList,
  getProxyStats,
  testProxy,
  batchTestProxies: batchTestProxy,
  deleteProxy,
  batchDeleteProxies: batchDeleteProxy,
  updateProxy,
  createProxy
}