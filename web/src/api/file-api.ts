import api from './index'
import type { ApiResponse } from './index'
import type { FileUploadResponse, BackupInfo } from './system-types'

// 文件上传API
export const uploadFile = (file: File, category?: string) => {
  const formData = new FormData()
  formData.append('file', file)
  if (category) {
    formData.append('category', category)
  }
  return api.post<ApiResponse<FileUploadResponse>>('/files/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 批量文件上传
export const batchUploadFiles = (files: File[], category?: string) => {
  const formData = new FormData()
  files.forEach(file => {
    formData.append('files', file)
  })
  if (category) {
    formData.append('category', category)
  }
  return api.post<ApiResponse<FileUploadResponse[]>>('/files/upload/batch', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 图片上传和处理
export const uploadImage = (file: File, options?: {
  width?: number
  height?: number
  quality?: number
  crop?: boolean
}) => {
  const formData = new FormData()
  formData.append('image', file)
  if (options) {
    formData.append('options', JSON.stringify(options))
  }
  return api.post<ApiResponse<FileUploadResponse & {
    thumbnails?: Array<{
      size: string
      url: string
    }>
  }>>('/files/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 文件管理
export const getFileList = (params?: {
  category?: string
  type?: string
  page?: number
  pageSize?: number
}) => {
  return api.get<ApiResponse<{
    files: Array<{
      id: string
      filename: string
      original_name: string
      url: string
      size: number
      type: string
      category: string
      created_at: string
    }>
    total: number
    page: number
    pageSize: number
  }>>('/files', { params })
}

export const deleteFile = (id: string) => {
  return api.delete<ApiResponse<void>>(`/files/${id}`)
}

export const batchDeleteFiles = (ids: string[]) => {
  return api.delete<ApiResponse<void>>('/files/batch', { data: { ids } })
}

// 文件信息
export const getFileInfo = (id: string) => {
  return api.get<ApiResponse<{
    id: string
    filename: string
    original_name: string
    url: string
    size: number
    type: string
    category: string
    metadata: Record<string, any>
    created_at: string
    updated_at: string
  }>>(`/files/${id}`)
}

// 文件统计
export const getFileStats = () => {
  return api.get<ApiResponse<{
    total_files: number
    total_size: number
    categories: Array<{
      category: string
      count: number
      size: number
    }>
    types: Array<{
      type: string
      count: number
      size: number
    }>
    recent_uploads: Array<{
      filename: string
      size: number
      created_at: string
    }>
  }>>('/files/stats')
}

// 系统备份API
export const createSystemBackup = () => {
  return api.post<ApiResponse<BackupInfo>>('/monitoring/backup')
}

export const getBackupList = () => {
  return api.get<ApiResponse<BackupInfo[]>>('/monitoring/backups')
}

export const downloadBackup = (filename: string) => {
  return api.get<Blob>(`/monitoring/backups/${filename}/download`, {
    responseType: 'blob'
  })
}

export const deleteBackup = (filename: string) => {
  return api.delete<ApiResponse<void>>(`/monitoring/backups/${filename}`)
}

// 备份恢复
export const restoreBackup = (filename: string) => {
  return api.post<ApiResponse<{
    success: boolean
    message: string
    restored_at: string
  }>>(`/monitoring/backups/${filename}/restore`)
}

// 备份配置
export const getBackupConfig = () => {
  return api.get<ApiResponse<{
    auto_backup: boolean
    frequency: 'daily' | 'weekly' | 'monthly'
    retention_days: number
    compress: boolean
    storage_location: string
    max_backups: number
  }>>('/monitoring/backup/config')
}

export const updateBackupConfig = (config: any) => {
  return api.patch<ApiResponse<void>>('/monitoring/backup/config', config)
}

// 文件清理
export const cleanupFiles = (params: {
  older_than_days: number
  categories?: string[]
  types?: string[]
}) => {
  return api.delete<ApiResponse<{
    deleted_count: number
    freed_space: number
  }>>('/files/cleanup', { data: params })
}

// 文件存储分析
export const getStorageAnalysis = () => {
  return api.get<ApiResponse<{
    total_used: number
    total_available: number
    usage_percentage: number
    growth_trend: Array<{
      date: string
      size: number
    }>
    largest_files: Array<{
      filename: string
      size: number
      created_at: string
    }>
    category_breakdown: Array<{
      category: string
      size: number
      percentage: number
    }>
  }>>('/files/storage-analysis')
}

// 文件搜索
export const searchFiles = (query: string, filters?: {
  category?: string
  type?: string
  size_min?: number
  size_max?: number
  date_from?: string
  date_to?: string
}) => {
  return api.post<ApiResponse<any[]>>('/files/search', {
    query,
    filters
  })
}

// 文件安全扫描
export const scanFiles = (fileIds?: string[]) => {
  return api.post<ApiResponse<{
    scanned_count: number
    threats_found: number
    scan_results: Array<{
      file_id: string
      filename: string
      status: 'safe' | 'threat' | 'suspicious'
      details?: string
    }>
  }>>('/files/scan', { fileIds })
}

// 文件API对象导出
export const fileApi = {
  uploadFile,
  batchUploadFiles,
  uploadImage,
  getFileList,
  deleteFile,
  batchDeleteFiles,
  getFileInfo,
  getFileStats,
  createSystemBackup,
  getBackupList,
  downloadBackup,
  deleteBackup,
  restoreBackup,
  getBackupConfig,
  updateBackupConfig,
  cleanupFiles,
  getStorageAnalysis,
  searchFiles,
  scanFiles
}
