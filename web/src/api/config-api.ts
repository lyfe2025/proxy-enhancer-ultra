import api from './index'
import type { ApiResponse } from './index'
import type { SystemConfig } from './system-types'

// 系统配置API
export const getSystemConfigs = (category?: string) => {
  return api.get<ApiResponse<SystemConfig[]>>('/system/configs', {
    params: { category }
  })
}

export const getSystemConfig = (key: string) => {
  return api.get<ApiResponse<{ key: string; value: any; description?: string }>>(`/system/configs/${key}`)
}

export const updateSystemConfig = (key: string, value: string) => {
  return api.patch<ApiResponse<void>>('/system/configs', { key, value })
}

export const batchUpdateSystemConfigs = (configs: { key: string; value: string }[]) => {
  return api.patch<ApiResponse<void>>('/system/configs/batch', { configs })
}

// 按分类获取配置
export const getBasicConfig = () => {
  return api.get<ApiResponse<SystemConfig['basic']>>('/system/configs/basic')
}

export const updateBasicConfig = (config: SystemConfig['basic']) => {
  return api.patch<ApiResponse<void>>('/system/configs/basic', config)
}

export const getSecurityConfig = () => {
  return api.get<ApiResponse<SystemConfig['security']>>('/system/configs/security')
}

export const updateSecurityConfig = (config: SystemConfig['security']) => {
  return api.patch<ApiResponse<void>>('/system/configs/security', config)
}

export const getEmailConfig = () => {
  return api.get<ApiResponse<SystemConfig['email']>>('/system/configs/email')
}

export const updateEmailConfig = (config: SystemConfig['email']) => {
  return api.patch<ApiResponse<void>>('/system/configs/email', config)
}

export const getSmsConfig = () => {
  return api.get<ApiResponse<SystemConfig['sms']>>('/system/configs/sms')
}

export const updateSmsConfig = (config: SystemConfig['sms']) => {
  return api.patch<ApiResponse<void>>('/system/configs/sms', config)
}

export const getStorageConfig = () => {
  return api.get<ApiResponse<SystemConfig['storage']>>('/system/configs/storage')
}

export const updateStorageConfig = (config: SystemConfig['storage']) => {
  return api.patch<ApiResponse<void>>('/system/configs/storage', config)
}

export const getLoggingConfig = () => {
  return api.get<ApiResponse<SystemConfig['logging']>>('/system/configs/logging')
}

export const updateLoggingConfig = (config: SystemConfig['logging']) => {
  return api.patch<ApiResponse<void>>('/system/configs/logging', config)
}

export const getMonitoringConfig = () => {
  return api.get<ApiResponse<SystemConfig['monitoring']>>('/system/configs/monitoring')
}

export const updateMonitoringConfig = (config: SystemConfig['monitoring']) => {
  return api.patch<ApiResponse<void>>('/system/configs/monitoring', config)
}

// 配置验证
export const validateConfig = (key: string, value: any) => {
  return api.post<ApiResponse<{
    valid: boolean
    errors?: string[]
  }>>('/system/configs/validate', { key, value })
}

// 测试配置
export const testEmailConfig = (config: SystemConfig['email']) => {
  return api.post<ApiResponse<{
    success: boolean
    message: string
  }>>('/system/configs/test/email', config)
}

export const testSmsConfig = (config: SystemConfig['sms']) => {
  return api.post<ApiResponse<{
    success: boolean
    message: string
  }>>('/system/configs/test/sms', config)
}

// 重置配置
export const resetConfig = (key: string) => {
  return api.delete<ApiResponse<void>>(`/system/configs/${key}`)
}

export const resetAllConfigs = () => {
  return api.delete<ApiResponse<void>>('/system/configs/all')
}

// 导出配置
export const exportConfigs = (category?: string) => {
  return api.get<Blob>('/system/configs/export', {
    params: { category },
    responseType: 'blob'
  })
}

// 导入配置
export const importConfigs = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post<ApiResponse<{
    success_count: number
    failed_count: number
    errors: string[]
  }>>('/system/configs/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 配置历史记录
export const getConfigHistory = (key: string) => {
  return api.get<ApiResponse<Array<{
    id: string
    key: string
    old_value: any
    new_value: any
    operator: string
    created_at: string
  }>>>(`/system/configs/${key}/history`)
}

// 系统配置API对象导出
export const configApi = {
  getSystemConfigs,
  getSystemConfig,
  updateSystemConfig,
  batchUpdateSystemConfigs,
  getBasicConfig,
  updateBasicConfig,
  getSecurityConfig,
  updateSecurityConfig,
  getEmailConfig,
  updateEmailConfig,
  getSmsConfig,
  updateSmsConfig,
  getStorageConfig,
  updateStorageConfig,
  getLoggingConfig,
  updateLoggingConfig,
  getMonitoringConfig,
  updateMonitoringConfig,
  validateConfig,
  testEmailConfig,
  testSmsConfig,
  resetConfig,
  resetAllConfigs,
  exportConfigs,
  importConfigs,
  getConfigHistory
}
