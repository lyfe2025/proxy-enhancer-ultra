import api from './index'
import type { ApiResponse, PageResponse, PageParams } from './index'

// 弹窗配置类型
export interface PopupConfig {
  id?: number
  name: string
  title: string
  description?: string
  type: 'form' | 'survey' | 'feedback' | 'subscription'
  status: 'active' | 'inactive' | 'draft'
  triggerRules: PopupTriggerRule[]
  formFields: FormField[]
  fields: FormField[] // 兼容性字段
  triggerType: string // 兼容性字段
  styling: PopupStyling
  settings: PopupSettings
  createdAt?: string
  updatedAt?: string
  viewCount?: number
  submitCount?: number
  conversionRate?: number
}

// 弹窗触发规则
export interface PopupTriggerRule {
  id?: number
  type: 'page_load' | 'time_delay' | 'scroll_percentage' | 'exit_intent' | 'click_element'
  value: string
  operator?: 'equals' | 'contains' | 'greater_than' | 'less_than'
}

// 表单字段
export interface FormField {
  id?: string
  name: string
  label: string
  type: 'text' | 'email' | 'phone' | 'number' | 'textarea' | 'select' | 'radio' | 'checkbox' | 'date'
  required: boolean
  placeholder?: string
  options?: string[]
  validation?: FieldValidation
  order: number
}

// 字段验证规则
export interface FieldValidation {
  minLength?: number
  maxLength?: number
  pattern?: string
  min?: number
  max?: number
  customMessage?: string
}

// 弹窗样式
export interface PopupStyling {
  width?: string
  height?: string
  position: 'center' | 'top' | 'bottom' | 'left' | 'right'
  backgroundColor?: string
  textColor?: string
  borderRadius?: string
  shadow?: boolean
  overlay?: boolean
  overlayColor?: string
}

// 弹窗设置
export interface PopupSettings {
  showCloseButton: boolean
  closeOnOverlayClick: boolean
  autoClose?: number
  showOnce?: boolean
  cookieExpiry?: number
  mobileOptimized: boolean
}

// 提交数据
export interface SubmissionData {
  id?: number
  popupId: number
  popupName: string
  data: Record<string, any>
  userAgent?: string
  ipAddress?: string
  referrer?: string
  submittedAt: string
  processed?: boolean
  tags?: string[]
}

// 数据统计
export interface DataStats {
  totalPopups: number
  activePopups: number
  totalSubmissions: number
  conversionRate: number
  avgSubmissionsPerDay: number
  topPerformingPopup?: {
    id: number
    name: string
    conversionRate: number
  }
}

// 获取弹窗列表
export const getPopupList = (params: PageParams & {
  type?: string
  status?: string
}) => {
  return api.get<ApiResponse<PageResponse<PopupConfig>>>('/data/popups', { params })
}

// 获取弹窗详情
export const getPopupDetail = (id: number) => {
  return api.get<ApiResponse<PopupConfig>>(`/data/popups/${id}`)
}

// 创建弹窗
export const createPopup = (data: Omit<PopupConfig, 'id' | 'createdAt' | 'updatedAt'>) => {
  return api.post<ApiResponse<PopupConfig>>('/data/popups', data)
}

// 更新弹窗
export const updatePopup = (id: number, data: Partial<PopupConfig>) => {
  return api.put<ApiResponse<PopupConfig>>(`/data/popups/${id}`, data)
}

// 删除弹窗
export const deletePopup = (id: number) => {
  return api.delete<ApiResponse<void>>(`/data/popups/${id}`)
}

// 批量删除弹窗
export const batchDeletePopup = (ids: number[]) => {
  return api.delete<ApiResponse<void>>('/data/popups/batch', { data: { ids } })
}

// 启用/禁用弹窗
export const togglePopupStatus = (id: number, status: 'active' | 'inactive') => {
  return api.patch<ApiResponse<void>>(`/data/popups/${id}/status`, { status })
}

// 复制弹窗
export const duplicatePopup = (id: number, name?: string) => {
  return api.post<ApiResponse<PopupConfig>>(`/data/popups/${id}/duplicate`, { name })
}

// 预览弹窗
export const previewPopup = (id: number) => {
  return api.get<ApiResponse<{ html: string; css: string; js: string }>>(`/data/popups/${id}/preview`)
}

// 获取提交数据列表
export const getSubmissionList = (params: PageParams & {
  popupId?: number
  startTime?: string
  endTime?: string
  processed?: boolean
}) => {
  return api.get<ApiResponse<PageResponse<SubmissionData>>>('/data/submissions', { params })
}

// 获取提交数据详情
export const getSubmissionDetail = (id: number) => {
  return api.get<ApiResponse<SubmissionData>>(`/data/submissions/${id}`)
}

// 标记提交数据为已处理
export const markSubmissionProcessed = (id: number) => {
  return api.patch<ApiResponse<void>>(`/data/submissions/${id}/processed`)
}

// 批量标记提交数据为已处理
export const batchMarkSubmissionProcessed = (ids: number[]) => {
  return api.patch<ApiResponse<void>>('/data/submissions/batch-processed', { ids })
}

// 删除提交数据
export const deleteSubmission = (id: number) => {
  return api.delete<ApiResponse<void>>(`/data/submissions/${id}`)
}

// 批量删除提交数据
export const batchDeleteSubmission = (ids: number[]) => {
  return api.delete<ApiResponse<void>>('/data/submissions/batch', { data: { ids } })
}

// 添加标签到提交数据
export const addSubmissionTags = (id: number, tags: string[]) => {
  return api.patch<ApiResponse<void>>(`/data/submissions/${id}/tags`, { tags })
}

// 获取数据统计
export const getDataStats = () => {
  return api.get<ApiResponse<DataStats>>('/data/stats')
}

// 导出提交数据
export const exportSubmissions = (params: {
  popupId?: number
  startTime?: string
  endTime?: string
  format?: 'csv' | 'excel' | 'json'
}) => {
  return api.post<Blob>('/data/submissions/export', params, {
    responseType: 'blob'
  })
}

// 获取弹窗分析数据
export const getPopupAnalytics = (id: number, params: {
  startTime?: string
  endTime?: string
  granularity?: 'hour' | 'day' | 'week' | 'month'
}) => {
  return api.get<ApiResponse<{
    views: Array<{ time: string; count: number }>
    submissions: Array<{ time: string; count: number }>
    conversionRate: Array<{ time: string; rate: number }>
    topSources: Array<{ source: string; count: number }>
    deviceStats: Array<{ device: string; count: number }>
  }>>(`/data/popups/${id}/analytics`, { params })
}

// 创建 dataApi 对象，包含所有相关函数
export const dataApi = {
  getPopups: getPopupList,
  getSubmissions: getSubmissionList,
  getDataStats,
  deletePopup,
  batchUpdatePopups: togglePopupStatus,
  batchDeletePopups: batchDeletePopup,
  updatePopup,
  createPopup,
  updateSubmission: markSubmissionProcessed,
  deleteSubmission,
  batchUpdateSubmissions: batchMarkSubmissionProcessed,
  batchDeleteSubmissions: batchDeleteSubmission,
  exportSubmissions
}