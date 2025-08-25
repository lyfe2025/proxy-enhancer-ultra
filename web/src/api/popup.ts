import { api } from './request'
import type { Popup, Submission, ApiResponse, PaginatedResponse } from '@/types'

// 弹窗管理相关API
export const popupApi = {
  // 获取弹窗列表
  getPopups(params?: {
    page?: number
    page_size?: number
    search?: string
    is_active?: boolean
  }): Promise<ApiResponse<PaginatedResponse<Popup>>> {
    return api.get('/popups', { params })
  },
  
  // 获取单个弹窗
  getPopup(id: number): Promise<ApiResponse<Popup>> {
    return api.get(`/popups/${id}`)
  },
  
  // 创建弹窗
  createPopup(data: Omit<Popup, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<Popup>> {
    return api.post('/popups', data)
  },
  
  // 更新弹窗
  updatePopup(id: number, data: Partial<Omit<Popup, 'id' | 'created_at' | 'updated_at'>>): Promise<ApiResponse<Popup>> {
    return api.put(`/popups/${id}`, data)
  },
  
  // 删除弹窗
  deletePopup(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/popups/${id}`)
  },
  
  // 启用/禁用弹窗
  togglePopup(id: number, is_active: boolean): Promise<ApiResponse<Popup>> {
    return api.patch(`/popups/${id}/toggle`, { is_active })
  },
  
  // 预览弹窗
  previewPopup(id: number): Promise<ApiResponse<{ html: string; css: string; js: string }>> {
    return api.get(`/popups/${id}/preview`)
  },
  
  // 复制弹窗
  duplicatePopup(id: number, name: string): Promise<ApiResponse<Popup>> {
    return api.post(`/popups/${id}/duplicate`, { name })
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
  }): Promise<ApiResponse<PaginatedResponse<Submission>>> {
    return api.get('/submissions', { params })
  },
  
  // 获取单个提交数据
  getSubmission(id: number): Promise<ApiResponse<Submission>> {
    return api.get(`/submissions/${id}`)
  },
  
  // 删除提交数据
  deleteSubmission(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/submissions/${id}`)
  },
  
  // 批量删除提交数据
  batchDeleteSubmissions(ids: number[]): Promise<ApiResponse<null>> {
    return api.post('/submissions/batch-delete', { ids })
  },
  
  // 导出提交数据
  exportSubmissions(params?: {
    popup_id?: number
    start_date?: string
    end_date?: string
    format?: 'csv' | 'excel' | 'json'
  }): Promise<ApiResponse<{ download_url: string }>> {
    return api.get('/submissions/export', { params })
  },
  
  // 获取提交数据统计
  getSubmissionStats(params?: {
    popup_id?: number
    start_date?: string
    end_date?: string
    group_by?: 'day' | 'week' | 'month'
  }): Promise<ApiResponse<{
    total_submissions: number
    unique_visitors: number
    conversion_rate: number
    chart_data: { date: string; count: number }[]
  }>> {
    return api.get('/submissions/stats', { params })
  }
}