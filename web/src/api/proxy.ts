import { api } from './request'
import type { ProxyConfig, Domain, ApiResponse, PaginatedResponse } from '@/types'

// 代理配置相关API
export const proxyApi = {
  // 获取代理配置列表
  getProxyConfigs(params?: {
    page?: number
    page_size?: number
    search?: string
    is_active?: boolean
  }): Promise<ApiResponse<PaginatedResponse<ProxyConfig>>> {
    return api.get('/proxy/configs', { params })
  },
  
  // 获取单个代理配置
  getProxyConfig(id: number): Promise<ApiResponse<ProxyConfig>> {
    return api.get(`/proxy/configs/${id}`)
  },
  
  // 创建代理配置
  createProxyConfig(data: Omit<ProxyConfig, 'id' | 'created_at' | 'updated_at' | 'domains'>): Promise<ApiResponse<ProxyConfig>> {
    return api.post('/proxy/configs', data)
  },
  
  // 更新代理配置
  updateProxyConfig(id: number, data: Partial<Omit<ProxyConfig, 'id' | 'created_at' | 'updated_at' | 'domains'>>): Promise<ApiResponse<ProxyConfig>> {
    return api.put(`/proxy/configs/${id}`, data)
  },
  
  // 删除代理配置
  deleteProxyConfig(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/proxy/configs/${id}`)
  },
  
  // 启用/禁用代理配置
  toggleProxyConfig(id: number, is_active: boolean): Promise<ApiResponse<ProxyConfig>> {
    return api.patch(`/proxy/configs/${id}/toggle`, { is_active })
  },
  
  // 测试代理配置
  testProxyConfig(id: number): Promise<ApiResponse<{ success: boolean; message: string; response_time?: number }>> {
    return api.post(`/proxy/configs/${id}/test`)
  }
}

// 域名管理相关API
export const domainApi = {
  // 获取域名列表
  getDomains(proxy_config_id?: number): Promise<ApiResponse<Domain[]>> {
    const params = proxy_config_id ? { proxy_config_id } : {}
    return api.get('/proxy/domains', { params })
  },
  
  // 获取单个域名
  getDomain(id: number): Promise<ApiResponse<Domain>> {
    return api.get(`/proxy/domains/${id}`)
  },
  
  // 创建域名
  createDomain(data: Omit<Domain, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<Domain>> {
    return api.post('/proxy/domains', data)
  },
  
  // 更新域名
  updateDomain(id: number, data: Partial<Omit<Domain, 'id' | 'created_at' | 'updated_at'>>): Promise<ApiResponse<Domain>> {
    return api.put(`/proxy/domains/${id}`, data)
  },
  
  // 删除域名
  deleteDomain(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/proxy/domains/${id}`)
  },
  
  // 启用/禁用域名
  toggleDomain(id: number, is_active: boolean): Promise<ApiResponse<Domain>> {
    return api.patch(`/proxy/domains/${id}/toggle`, { is_active })
  }
}