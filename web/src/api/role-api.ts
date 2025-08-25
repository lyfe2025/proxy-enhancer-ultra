import api from './index'
import type { ApiResponse, PageResponse, PageParams } from './index'
import type { Role, RoleFormData } from './system-types'

// 角色管理API
export const getRoleList = (params: PageParams & {
  status?: string
}) => {
  return api.get<ApiResponse<PageResponse<Role>>>('/system/roles', { params })
}

export const getRoleDetail = (id: string) => {
  return api.get<ApiResponse<Role>>(`/system/roles/${id}`)
}

export const createRole = (data: Omit<Role, 'id' | 'created_at' | 'updated_at'>) => {
  return api.post<ApiResponse<Role>>('/system/roles', data)
}

export const updateRole = (id: string, data: Partial<Role>) => {
  return api.put<ApiResponse<Role>>(`/system/roles/${id}`, data)
}

export const deleteRole = (id: string) => {
  return api.delete<ApiResponse<void>>(`/system/roles/${id}`)
}

export const toggleRoleStatus = (id: string, status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>(`/system/roles/${id}/status`, { status })
}

export const assignRolePermissions = (id: string, permissionIds: string[]) => {
  return api.patch<ApiResponse<void>>(`/system/roles/${id}/permissions`, { permissionIds })
}

// 获取所有角色（不分页）
export const getAllRoles = () => {
  return api.get<ApiResponse<Role[]>>('/system/roles/all')
}

// 获取角色权限
export const getRolePermissions = (id: string) => {
  return api.get<ApiResponse<string[]>>(`/system/roles/${id}/permissions`)
}

// 复制角色
export const copyRole = (id: string, newName: string) => {
  return api.post<ApiResponse<Role>>(`/system/roles/${id}/copy`, { name: newName })
}

// 角色排序
export const updateRoleSort = (roles: Array<{ id: string; sort: number }>) => {
  return api.patch<ApiResponse<void>>('/system/roles/sort', { roles })
}

// 获取角色统计信息
export const getRoleStats = () => {
  return api.get<ApiResponse<{
    total_roles: number
    active_roles: number
    inactive_roles: number
    roles_with_users: number
    permissions_count: number
  }>>('/system/roles/stats')
}

// 获取角色用户列表
export const getRoleUsers = (id: string, params?: PageParams) => {
  return api.get<ApiResponse<PageResponse<any>>>(`/system/roles/${id}/users`, { params })
}

// 批量分配用户到角色
export const assignUsersToRole = (roleId: string, userIds: string[]) => {
  return api.post<ApiResponse<void>>(`/system/roles/${roleId}/users`, { userIds })
}

// 移除角色中的用户
export const removeUsersFromRole = (roleId: string, userIds: string[]) => {
  return api.delete<ApiResponse<void>>(`/system/roles/${roleId}/users`, { data: { userIds } })
}

// 导出角色数据
export const exportRoles = (params?: {
  status?: string
  format?: 'csv' | 'excel'
}) => {
  return api.get<Blob>('/system/roles/export', {
    params,
    responseType: 'blob'
  })
}

// 角色权限矩阵
export const getRolePermissionMatrix = () => {
  return api.get<ApiResponse<{
    roles: Role[]
    permissions: any[]
    matrix: Record<string, string[]>
  }>>('/system/roles/permission-matrix')
}

// 角色API对象导出
export const roleApi = {
  getRoleList,
  getRoleDetail,
  createRole,
  updateRole,
  deleteRole,
  toggleRoleStatus,
  assignRolePermissions,
  getAllRoles,
  getRolePermissions,
  copyRole,
  updateRoleSort,
  getRoleStats,
  getRoleUsers,
  assignUsersToRole,
  removeUsersFromRole,
  exportRoles,
  getRolePermissionMatrix
}
