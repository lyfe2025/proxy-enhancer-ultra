import { api } from './request'
import type { Rule, ApiResponse, PaginatedResponse } from '@/types'

// 规则引擎相关API
export const rulesApi = {
  // 获取规则列表
  getRules(params?: {
    page?: number
    page_size?: number
    search?: string
    is_active?: boolean
  }): Promise<ApiResponse<PaginatedResponse<Rule>>> {
    return api.get('/rules', { params })
  },
  
  // 获取单个规则
  getRule(id: number): Promise<ApiResponse<Rule>> {
    return api.get(`/rules/${id}`)
  },
  
  // 创建规则
  createRule(data: Omit<Rule, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<Rule>> {
    return api.post('/rules', data)
  },
  
  // 更新规则
  updateRule(id: number, data: Partial<Omit<Rule, 'id' | 'created_at' | 'updated_at'>>): Promise<ApiResponse<Rule>> {
    return api.put(`/rules/${id}`, data)
  },
  
  // 删除规则
  deleteRule(id: number): Promise<ApiResponse<null>> {
    return api.delete(`/rules/${id}`)
  },
  
  // 启用/禁用规则
  toggleRule(id: number, is_active: boolean): Promise<ApiResponse<Rule>> {
    return api.patch(`/rules/${id}/toggle`, { is_active })
  },
  
  // 测试规则
  testRule(id: number, testData: any): Promise<ApiResponse<{ matched: boolean; result: any }>> {
    return api.post(`/rules/${id}/test`, testData)
  },
  
  // 批量更新规则优先级
  updateRulePriorities(rules: { id: number; priority: number }[]): Promise<ApiResponse<null>> {
    return api.post('/rules/priorities', { rules })
  },
  
  // 获取规则执行统计
  getRuleStats(id: number, params?: {
    start_date?: string
    end_date?: string
  }): Promise<ApiResponse<{
    total_executions: number
    successful_executions: number
    failed_executions: number
    avg_execution_time: number
  }>> {
    return api.get(`/rules/${id}/stats`, { params })
  }
}