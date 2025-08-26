import api from './index'
import type { ApiResponse, PageResponse, PageParams } from './index'
import type { User, UserFormData } from './system-types'

// 用户管理API
export const getUserList = (params: PageParams & {
  status?: string
  role_id?: string
}) => {
  return api.get<ApiResponse<PageResponse<User>>>('/users/admin', { params })
}

export const getUserDetail = (id: string) => {
  return api.get<ApiResponse<User>>(`/users/admin/${id}`)
}

export const createUser = (data: Omit<User, 'id' | 'created_at' | 'updated_at'>) => {
  return api.post<ApiResponse<User>>('/users/admin', data)
}

export const updateUser = (id: string, data: Partial<User>) => {
  return api.put<ApiResponse<User>>(`/users/admin/${id}`, data)
}

export const deleteUser = (id: string) => {
  return api.delete<ApiResponse<void>>(`/users/admin/${id}`)
}

export const batchDeleteUser = (ids: string[]) => {
  return api.delete<ApiResponse<void>>('/users/admin/batch', { data: { ids } })
}

export const toggleUserStatus = (id: string, status: 'active' | 'inactive' | 'locked') => {
  return api.patch<ApiResponse<void>>(`/users/admin/${id}/status`, { status })
}

export const resetUserPassword = (id: string, newPassword: string) => {
  return api.patch<ApiResponse<void>>(`/users/admin/${id}/password`, { password: newPassword })
}

// 用户信息管理 (管理员API)
export const getCurrentUser = () => {
  return api.get<ApiResponse<User>>('/users/profile')
}

export const updateCurrentUser = (data: Partial<UserFormData>) => {
  return api.put<ApiResponse<User>>('/users/profile', data)
}

export const changePassword = (data: {
  old_password: string
  new_password: string
  confirm_password: string
}) => {
  return api.patch<ApiResponse<void>>('/users/change-password', data)
}

export const uploadAvatar = (file: File) => {
  const formData = new FormData()
  formData.append('avatar', file)
  return api.post<ApiResponse<{ url: string }>>('/users/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 用户偏好设置
export const getUserPreferences = () => {
  return api.get<ApiResponse<Record<string, any>>>('/users/preferences')
}

export const updateUserPreferences = (preferences: Record<string, any>) => {
  return api.patch<ApiResponse<void>>('/users/preferences', { preferences })
}

// 用户活动记录
export const getUserActivities = (userId?: string, params?: PageParams) => {
  const url = userId ? `/users/admin/${userId}/activities` : '/users/admin/activities'
  return api.get<ApiResponse<PageResponse<any>>>(url, { params })
}

// 用户统计信息
export const getUserStats = () => {
  return api.get<ApiResponse<{
    total_users: number
    active_users: number
    inactive_users: number
    locked_users: number
    online_users: number
    new_users_today: number
    new_users_week: number
    new_users_month: number
  }>>('/monitoring/dashboard')
}

// 导出用户数据
export const exportUsers = (params?: {
  status?: string
  role_id?: string
  format?: 'csv' | 'excel'
}) => {
  return api.get<Blob>('/users/admin/export', {
    params,
    responseType: 'blob'
  })
}

// 批量导入用户
export const importUsers = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post<ApiResponse<{
    success_count: number
    failed_count: number
    errors: string[]
  }>>('/users/admin/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 用户API对象导出
export const userApi = {
  getUserList,
  getUserDetail,
  createUser,
  updateUser,
  deleteUser,
  batchDeleteUser,
  toggleUserStatus,
  resetUserPassword,
  getCurrentUser,
  updateCurrentUser,
  changePassword,
  uploadAvatar,
  getUserPreferences,
  updateUserPreferences,
  getUserActivities,
  getUserStats,
  exportUsers,
  importUsers
}
