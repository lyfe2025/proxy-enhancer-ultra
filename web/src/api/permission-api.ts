import api from './index'
import type { ApiResponse } from './index'
import type { Permission, PermissionFormData } from './system-types'

// 权限管理API
export const getPermissions = (params?: {
  page?: number
  limit?: number
  search?: string
  module?: string
  status?: string
}) => {
  return api.get<ApiResponse<PageResponse<Permission>>>('/users/admin/permissions', { params })
}

export const getPermissionDetail = (id: string) => {
  return api.get<ApiResponse<Permission>>(`/users/admin/permissions/${id}`)
}

export const createPermission = (data: Omit<Permission, 'id' | 'created_at' | 'updated_at'>) => {
  return api.post<ApiResponse<Permission>>('/users/admin/permissions', data)
}

export const updatePermission = (id: string, data: Partial<Permission>) => {
  return api.put<ApiResponse<Permission>>(`/users/admin/permissions/${id}`, data)
}

export const deletePermission = (id: string) => {
  return api.delete<ApiResponse<void>>(`/users/admin/permissions/${id}`)
}

export const togglePermissionStatus = (id: string, status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>(`/users/admin/permissions/${id}/status`, { status })
}

// 获取所有权限（不分页）
export const getAllPermissions = () => {
  return api.get<ApiResponse<Permission[]>>('/users/admin/permissions/all')
}

// 获取权限模块列表
export const getPermissionModules = () => {
  return api.get<ApiResponse<string[]>>('/users/admin/permissions/modules')
}

// 批量更新权限状态
export const batchUpdatePermissionStatus = (ids: string[], status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>('/users/admin/permissions/batch-status', { ids, status })
}

// 权限统计
export const getPermissionStats = () => {
  return api.get<ApiResponse<{
    total_permissions: number
    active_permissions: number
    inactive_permissions: number
    modules_count: number
    permissions_by_module: Record<string, number>
  }>>('/users/admin/permissions/stats')
}

// 导出权限数据
export const exportPermissions = (params?: {
  module?: string
  status?: string
  format?: 'csv' | 'excel'
}) => {
  return api.get<Blob>('/users/admin/permissions/export', {
    params,
    responseType: 'blob'
  })
}

// 权限API对象导出
export const permissionApi = {
  getPermissions,
  getAllPermissions,
  getPermissionDetail,
  createPermission,
  updatePermission,
  deletePermission,
  togglePermissionStatus,
  getPermissionModules,
  batchUpdatePermissionStatus,
  getPermissionStats,
  exportPermissions
}
