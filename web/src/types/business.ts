// 业务相关类型定义

// 字典类型
export interface Dictionary {
  id: string
  type: string
  code: string
  label: string
  value: string
  description?: string
  sort: number
  status: 'active' | 'inactive'
  parentId?: string
  children?: Dictionary[]
  createdAt: string
  updatedAt: string
}

// 文件类型
export interface FileInfo {
  id: string
  name: string
  originalName: string
  path: string
  url: string
  size: number
  type: string
  extension: string
  mimeType: string
  hash: string
  uploaderId: string
  uploaderName: string
  uploadTime: string
  downloadCount: number
  status: 'active' | 'deleted'
}
