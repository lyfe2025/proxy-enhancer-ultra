import api from './index'
import type { ApiResponse, PageResponse, PageParams } from './index'

// 规则类型
export interface Rule {
  id?: number
  name: string
  type: 'url_match' | 'domain_match' | 'header_match' | 'custom'
  pattern: string
  action: 'proxy' | 'direct' | 'block' | 'modify'
  proxyId?: number
  priority: number
  status: 'enabled' | 'disabled'
  description?: string
  conditions: RuleCondition[]
  actions: RuleAction[]
  createdAt?: string
  updatedAt?: string
  lastExecuted?: string
  executionCount?: number
  successCount?: number
  failureCount?: number
}

// 规则条件
export interface RuleCondition {
  id?: number
  type: 'url' | 'domain' | 'header' | 'method' | 'ip' | 'user_agent'
  operator: 'equals' | 'contains' | 'starts_with' | 'ends_with' | 'regex' | 'in' | 'not_in'
  value: string
  caseSensitive?: boolean
}

// 规则动作
export interface RuleAction {
  id?: number
  type: 'set_proxy' | 'set_header' | 'remove_header' | 'redirect' | 'block' | 'log'
  value: string
  params?: Record<string, any>
}

// 规则统计信息
export interface RuleStats {
  total: number
  enabled: number
  disabled: number
  todayExecutions: number
  successRate: number
  avgResponseTime: number
}

// 规则执行日志
export interface RuleExecutionLog {
  id: number
  ruleId: number
  ruleName: string
  url: string
  method: string
  status: 'success' | 'failed'
  responseTime: number
  error?: string
  executedAt: string
  userAgent?: string
  clientIp?: string
}

// 获取规则列表
export const getRuleList = (params: PageParams & {
  type?: string
  status?: string
  priority?: string
}) => {
  return api.get<ApiResponse<PageResponse<Rule>>>('/api/admin/rules', { params })
}

// 获取规则详情
export const getRuleDetail = (id: number) => {
  return api.get<ApiResponse<Rule>>(`/rules/${id}`)
}

// 创建规则
export const createRule = (data: Omit<Rule, 'id' | 'createdAt' | 'updatedAt'>) => {
  return api.post<ApiResponse<Rule>>('/api/admin/rules', data)
}

// 更新规则
export const updateRule = (id: number, data: Partial<Rule>) => {
  return api.put<ApiResponse<Rule>>(`/rules/${id}`, data)
}

// 删除规则
export const deleteRule = (id: number) => {
  return api.delete<ApiResponse<void>>(`/rules/${id}`)
}

// 批量删除规则
export const batchDeleteRule = (ids: number[]) => {
  return api.delete<ApiResponse<void>>('/rules/batch', { data: { ids } })
}

// 测试规则
export const testRule = (id: number, testUrl: string) => {
  return api.post<ApiResponse<{
    matched: boolean
    action: string
    proxyId?: number
    responseTime: number
  }>>(`/rules/${id}/test`, { url: testUrl })
}

// 启用/禁用规则
export const toggleRuleStatus = (id: number, status: 'enabled' | 'disabled') => {
  return api.patch<ApiResponse<void>>(`/rules/${id}/status`, { status })
}

// 批量启用/禁用规则
export const batchToggleRuleStatus = (ids: number[], status: 'enabled' | 'disabled') => {
  return api.patch<ApiResponse<void>>('/api/admin/rules/batch-status', { ids, status })
}

// 调整规则优先级
export const updateRulePriority = (id: number, priority: number) => {
  return api.patch<ApiResponse<void>>(`/rules/${id}/priority`, { priority })
}

// 批量调整规则优先级
export const batchUpdateRulePriority = (rules: { id: number; priority: number }[]) => {
  return api.patch<ApiResponse<void>>('/api/admin/rules/batch-priority', { rules })
}

// 获取规则统计信息
export const getRuleStats = () => {
  return api.get<ApiResponse<RuleStats>>('/api/admin/rules/stats')
}

// 获取规则执行日志
export const getRuleExecutionLogs = (params: PageParams & {
  ruleId?: number
  status?: string
  startTime?: string
  endTime?: string
}) => {
  return api.get<ApiResponse<PageResponse<RuleExecutionLog>>>('/api/admin/rules/logs', { params })
}

// 导入规则配置
export const importRuleConfig = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post<ApiResponse<{ success: number; failed: number }>>('/api/admin/rules/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 导出规则配置
export const exportRuleConfig = (ids?: number[]) => {
  return api.post<Blob>('/api/admin/rules/export', { ids }, {
    responseType: 'blob'
  })
}

// 复制规则
export const duplicateRule = (id: number, name?: string) => {
  return api.post<ApiResponse<Rule>>(`/rules/${id}/duplicate`, { name })
}

// 创建 rulesApi 对象，包含所有相关函数
export const rulesApi = {
  getRules: getRuleList,
  getRuleStats,
  toggleRuleStatus,
  testRule,
  batchUpdateRules: batchToggleRuleStatus,
  deleteRule,
  batchDeleteRules: batchDeleteRule,
  updateRule,
  createRule
}