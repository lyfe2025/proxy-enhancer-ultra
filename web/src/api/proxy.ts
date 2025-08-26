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
  tags?: string | string[]
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
  status?: string
  type?: string
  country?: string
}) => {
  return api.get<ApiResponse<PageResponse<ProxyConfig>>>('/proxy/configs', { params })
}

// 获取代理详情
export const getProxyDetail = (id: number) => {
  return api.get<ApiResponse<ProxyConfig>>(`/proxy/configs/${id}`)
}

// 创建代理
export const createProxy = (data: Omit<ProxyConfig, 'id' | 'createdAt' | 'updatedAt'>) => {
  return api.post<ApiResponse<ProxyConfig>>('/proxy/configs', data)
}

// 更新代理
export const updateProxy = (id: number, data: Partial<ProxyConfig>) => {
  return api.put<ApiResponse<ProxyConfig>>(`/proxy/configs/${id}`, data)
}

// 删除代理
export const deleteProxy = (id: number) => {
  return api.delete<ApiResponse<void>>(`/proxy/configs/${id}`)
}

// 批量删除代理
export const batchDeleteProxy = (ids: number[]) => {
  return api.delete<ApiResponse<void>>('/proxy/configs/batch', { data: { ids } })
}

// 测试代理连接
export const testProxy = (id: number) => {
  return api.post<ApiResponse<ProxyTestResult>>(`/proxy/configs/${id}/test`)
}

// 批量测试代理
export const batchTestProxy = (ids: string[]) => {
  return api.post<ApiResponse<ProxyTestResult[]>>('/proxy/configs/batch/test', { ids })
}

// 获取代理统计信息
export const getProxyStats = () => {
  return api.get<ApiResponse<ProxyStats>>('/proxy/configs/stats')
}

// 启用/禁用代理
export const toggleProxyStatus = (id: number, status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>(`/proxy/configs/${id}/status`, { status })
}

// 批量启用/禁用代理
export const batchToggleProxyStatus = (ids: string[], status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>('/proxy/configs/batch/status', { ids, status })
}

// 导入代理配置
export const importProxyConfig = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post<ApiResponse<{ imported: number; failed: number }>>('/proxy/configs/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 导出代理配置
export const exportProxyConfig = (proxyIds?: string[]) => {
  const params = proxyIds ? { ids: proxyIds.join(',') } : {}
  return api.get('/proxy/configs/export', {
    params,
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