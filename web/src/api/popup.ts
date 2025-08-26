import api from './index'
import type { Popup, Submission, PaginatedResponse } from '@/types'
import type { ApiResponse } from './index'

// 弹窗管理相关API
export const popupApi = {
  // 获取弹窗列表
  getPopups(params?: {
    page?: number
    page_size?: number
    search?: string
    is_active?: boolean
  }) {
    return api.get('/popups', { params })
  },
  
  // 获取单个弹窗
  getPopup(id: number) {
    return api.get<Popup>(`/popups/${id}`)
  },
  
  // 创建弹窗
  createPopup(data: Omit<Popup, 'id' | 'created_at' | 'updated_at'>) {
    return api.post<Popup>('/popups', data)
  },
  
  // 更新弹窗
  updatePopup(id: number, data: Partial<Omit<Popup, 'id' | 'created_at' | 'updated_at'>>) {
    return api.put<Popup>(`/popups/${id}`, data)
  },
  
  // 删除弹窗
  deletePopup(id: number) {
    return api.delete<null>(`/popups/${id}`)
  },
  
  // 启用/禁用弹窗
  togglePopup(id: number, is_active: boolean) {
    return api.post<Popup>(`/popups/${id}/toggle`, { is_active })
  },
  
  // 预览弹窗
  previewPopup(id: number) {
    return api.get<{ html: string; css: string; js: string }>(`/popups/${id}/preview`)
  },
  
  // 复制弹窗
  duplicatePopup(id: number, name: string) {
    return api.post<Popup>(`/popups/${id}/duplicate`, { name })
  }
}

// 提交数据管理相关API
export const submissionApi = {
  // 获取提交数据列表
  getSubmissions(params?: {
    page?: number
    page_size?: number
    popup_id?: number
    start_date?: string
    end_date?: string
    search?: string
  }) {
    return api.get<PaginatedResponse<Submission>>('/submissions', { params })
  },
  
  // 获取单个提交数据
  getSubmission(id: number) {
    return api.get<Submission>(`/submissions/${id}`)
  },
  
  // 删除提交数据
  deleteSubmission(id: number) {
    return api.delete<null>(`/submissions/${id}`)
  },
  
  // 批量删除提交数据
  batchDeleteSubmissions(ids: number[]) {
    return api.post<null>('/submissions/batch-delete', { ids })
  },
  
  // 导出提交数据
  exportSubmissions(params?: {
    popup_id?: number
    start_date?: string
    end_date?: string
    format?: 'csv' | 'excel' | 'json'
  }) {
    return api.get<{ download_url: string }>('/submissions/export', { params })
  },
  
  // 获取提交数据统计
  getSubmissionStats(params?: {
    popup_id?: number
    start_date?: string
    end_date?: string
    group_by?: 'day' | 'week' | 'month'
  }) {
    return api.get<{
      total_submissions: number
      unique_visitors: number
      conversion_rate: number
      chart_data: { date: string; count: number }[]
    }>('/submissions/stats', { params })
  }
}