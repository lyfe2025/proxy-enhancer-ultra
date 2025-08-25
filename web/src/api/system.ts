import { api } from './request'
import type { User, Role, Permission, SystemMetric, ProxyLog, DashboardStats, ApiResponse, PaginatedResponse } from '@/types'

// 用户管理相关API
export const userApi = {
  // 获取用户列表
  getUsers(params?: {
    page?: number
    page_size?: number
    search?: string
    is_active?: boolean
  }): Promise<ApiResponse<PaginatedResponse<User>>> {
    return api.get('/users', { params })
  },
  
  // 获取单个用户
  getUser(id: number): Promise<ApiResponse<User>> {
    return api.get(`/users/${id}`)
  },
  
  // 创建用户
  createUser(data: {
    username: string
    email: string
    password: string
    is_active?: boolean
  }): Promise<ApiResponse<User>> {
    return api.post('/users', data)
  },
  
  // 更新用户
  updateUser(id: number, data: Partial<{
    username: string
    email: string
    is_active: boolean
  }>): Promise<ApiResponse<User>> {
    return api.put(`/users/${id}`, data)
  },
  
  // 删除用户
  deleteUser(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/users/${id}`)
  },
  
  // 重置用户密码
  resetUserPassword(id: number, password: string): Promise<ApiResponse<null>> {
    return api.post(`/users/${id}/reset-password`, { password })
  },
  
  // 启用/禁用用户
  toggleUser(id: number, is_active: boolean): Promise<ApiResponse<User>> {
    return api.patch(`/users/${id}/toggle`, { is_active })
  }
}

// 角色管理相关API
export const roleApi = {
  // 获取角色列表
  getRoles(): Promise<ApiResponse<Role[]>> {
    return api.get('/roles')
  },
  
  // 获取单个角色
  getRole(id: number): Promise<ApiResponse<Role>> {
    return api.get(`/roles/${id}`)
  },
  
  // 创建角色
  createRole(data: Omit<Role, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<Role>> {
    return api.post('/roles', data)
  },
  
  // 更新角色
  updateRole(id: number, data: Partial<Omit<Role, 'id' | 'created_at' | 'updated_at'>>): Promise<ApiResponse<Role>> {
    return api.put(`/roles/${id}`, data)
  },
  
  // 删除角色
  deleteRole(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/roles/${id}`)
  },
  
  // 获取角色权限
  getRolePermissions(id: number): Promise<ApiResponse<Permission[]>> {
    return api.get(`/roles/${id}/permissions`)
  },
  
  // 更新角色权限
  updateRolePermissions(id: number, permission_ids: number[]): Promise<ApiResponse<null>> {
    return api.put(`/roles/${id}/permissions`, { permission_ids })
  }
}

// 权限管理相关API
export const permissionApi = {
  // 获取权限列表
  getPermissions(): Promise<ApiResponse<Permission[]>> {
    return api.get('/permissions')
  },
  
  // 获取单个权限
  getPermission(id: number): Promise<ApiResponse<Permission>> {
    return api.get(`/permissions/${id}`)
  },
  
  // 创建权限
  createPermission(data: Omit<Permission, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<Permission>> {
    return api.post('/permissions', data)
  },
  
  // 更新权限
  updatePermission(id: number, data: Partial<Omit<Permission, 'id' | 'created_at' | 'updated_at'>>): Promise<ApiResponse<Permission>> {
    return api.put(`/permissions/${id}`, data)
  },
  
  // 删除权限
  deletePermission(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/permissions/${id}`)
  }
}

// 系统监控相关API
export const monitorApi = {
  // 获取仪表板统计数据
  getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
    return api.get('/monitor/dashboard')
  },
  
  // 获取系统指标
  getSystemMetrics(params?: {
    metric_name?: string
    start_date?: string
    end_date?: string
    limit?: number
  }): Promise<ApiResponse<SystemMetric[]>> {
    return api.get('/monitor/metrics', { params })
  },
  
  // 获取代理日志
  getProxyLogs(params?: {
    page?: number
    page_size?: number
    start_date?: string
    end_date?: string
    status_code?: number
    method?: string
    search?: string
  }): Promise<ApiResponse<PaginatedResponse<ProxyLog>>> {
    return api.get('/monitor/logs', { params })
  },
  
  // 获取系统健康状态
  getSystemHealth(): Promise<ApiResponse<{
    status: 'healthy' | 'warning' | 'critical'
    uptime: number
    memory_usage: number
    cpu_usage: number
    disk_usage: number
    database_status: 'connected' | 'disconnected'
    services: {
      name: string
      status: 'running' | 'stopped' | 'error'
      last_check: string
    }[]
  }>> {
    return api.get('/monitor/health')
  },
  
  // 清理日志
  cleanupLogs(params: {
    before_date: string
    log_types?: string[]
  }): Promise<ApiResponse<{ deleted_count: number }>> {
    return api.post('/monitor/cleanup', params)
  }
}