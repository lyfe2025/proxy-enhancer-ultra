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

// 规则创建请求
export interface RuleCreateRequest {
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
}

// 规则更新请求
export interface RuleUpdateRequest {
  name?: string
  type?: 'url_match' | 'domain_match' | 'header_match' | 'custom'
  pattern?: string
  action?: 'proxy' | 'direct' | 'block' | 'modify'
  proxyId?: number
  priority?: number
  description?: string
  conditions?: RuleCondition[]
  actions?: RuleAction[]
}
