import api from './index'
import type { ApiResponse } from './index'
import type { Permission, PermissionFormData } from './system-types'

// 权限管理API
export const getPermissionList = (params?: {
  type?: string
  status?: string
  module?: string
}) => {
  return api.get<ApiResponse<Permission[]>>('/system/permissions', { params })
}

export const getPermissionTree = () => {
  return api.get<ApiResponse<Permission[]>>('/system/permissions/tree')
}

export const getPermissionDetail = (id: string) => {
  return api.get<ApiResponse<Permission>>(`/system/permissions/${id}`)
}

export const createPermission = (data: Omit<Permission, 'id'>) => {
  return api.post<ApiResponse<Permission>>('/system/permissions', data)
}

export const updatePermission = (id: string, data: Partial<Permission>) => {
  return api.put<ApiResponse<Permission>>(`/system/permissions/${id}`, data)
}

export const deletePermission = (id: string) => {
  return api.delete<ApiResponse<void>>(`/system/permissions/${id}`)
}

export const togglePermissionStatus = (id: string, status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>(`/system/permissions/${id}/status`, { status })
}

// 获取权限模块列表
export const getPermissionModules = () => {
  return api.get<ApiResponse<Array<{
    name: string
    label: string
    description?: string
    permission_count: number
  }>>>('/system/permissions/modules')
}

// 根据模块获取权限
export const getPermissionsByModule = (module: string) => {
  return api.get<ApiResponse<Permission[]>>(`/system/permissions/modules/${module}`)
}

// 权限排序
export const updatePermissionSort = (permissions: Array<{ id: string; sort: number }>) => {
  return api.patch<ApiResponse<void>>('/system/permissions/sort', { permissions })
}

// 批量创建权限
export const batchCreatePermissions = (permissions: Omit<Permission, 'id'>[]) => {
  return api.post<ApiResponse<Permission[]>>('/system/permissions/batch', { permissions })
}

// 批量删除权限
export const batchDeletePermissions = (ids: string[]) => {
  return api.delete<ApiResponse<void>>('/system/permissions/batch', { data: { ids } })
}

// 权限验证
export const checkPermission = (permission: string) => {
  return api.post<ApiResponse<{ hasPermission: boolean }>>('/system/permissions/check', { permission })
}

// 获取用户权限
export const getUserPermissions = (userId?: string) => {
  const url = userId ? `/system/permissions/user/${userId}` : '/system/permissions/user'
  return api.get<ApiResponse<string[]>>(url)
}

// 权限统计
export const getPermissionStats = () => {
  return api.get<ApiResponse<{
    total_permissions: number
    menu_permissions: number
    button_permissions: number
    api_permissions: number
    active_permissions: number
    inactive_permissions: number
    modules_count: number
  }>>('/system/permissions/stats')
}

// 权限使用情况
export const getPermissionUsage = (id: string) => {
  return api.get<ApiResponse<{
    permission: Permission
    roles: Array<{
      id: string
      name: string
      user_count: number
    }>
    total_users: number
  }>>(`/system/permissions/${id}/usage`)
}

// 导出权限数据
export const exportPermissions = (params?: {
  type?: string
  status?: string
  module?: string
  format?: 'csv' | 'excel'
}) => {
  return api.get<Blob>('/system/permissions/export', {
    params,
    responseType: 'blob'
  })
}

// 权限依赖检查
export const checkPermissionDependencies = (id: string) => {
  return api.get<ApiResponse<{
    can_delete: boolean
    dependencies: Array<{
      type: 'role' | 'user'
      name: string
      count: number
    }>
  }>>(`/system/permissions/${id}/dependencies`)
}

// 权限API对象导出
export const permissionApi = {
  getPermissionList,
  getPermissionTree,
  getPermissionDetail,
  createPermission,
  updatePermission,
  deletePermission,
  togglePermissionStatus,
  getPermissionModules,
  getPermissionsByModule,
  updatePermissionSort,
  batchCreatePermissions,
  batchDeletePermissions,
  checkPermission,
  getUserPermissions,
  getPermissionStats,
  getPermissionUsage,
  exportPermissions,
  checkPermissionDependencies
}
