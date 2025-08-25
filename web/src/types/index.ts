// 导出拆分后的类型模块
export * from './user'
export * from './config'
export * from './monitoring'
export * from './logging'
export * from './common'
export * from './ui'
export * from './business'

// 基础类型定义
export interface BaseModel {
  id: number
  created_at: string
  updated_at: string
}

// 用户相关类型
export interface User extends BaseModel {
  username: string
  email: string
  is_active: boolean
  last_login?: string
}

export interface Role extends BaseModel {
  name: string
  description: string
}

export interface Permission extends BaseModel {
  name: string
  description: string
  resource: string
  action: string
}

// 认证相关类型
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
  expires_at: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

// 代理配置类型
export interface ProxyConfig extends BaseModel {
  name: string
  target_url: string
  is_active: boolean
  description?: string
  domains: Domain[]
}

export interface Domain extends BaseModel {
  domain: string
  proxy_config_id: number
  is_active: boolean
}

// 规则引擎类型
export interface Rule extends BaseModel {
  name: string
  description?: string
  conditions: RuleCondition[]
  actions: RuleAction[]
  is_active: boolean
  priority: number
}

export interface RuleCondition {
  field: string
  operator: string
  value: string
}

export interface RuleAction {
  type: string
  config: Record<string, any>
}

// 弹窗表单类型
export interface Popup extends BaseModel {
  title: string
  content: string
  form_config: FormConfig
  trigger_rules: TriggerRule[]
  is_active: boolean
}

export interface FormConfig {
  fields: FormField[]
  submit_url: string
  success_message?: string
}

export interface FormField {
  name: string
  label: string
  type: 'text' | 'email' | 'number' | 'select' | 'textarea' | 'checkbox'
  required: boolean
  options?: string[]
  placeholder?: string
  validation?: FieldValidation
}

export interface FieldValidation {
  min_length?: number
  max_length?: number
  pattern?: string
  message?: string
}

export interface TriggerRule {
  type: 'url' | 'time' | 'scroll' | 'click'
  config: Record<string, any>
}

// 提交数据类型
export interface Submission extends BaseModel {
  popup_id: number
  form_data: Record<string, any>
  user_agent?: string
  ip_address?: string
  referrer?: string
}

// 系统监控类型
export interface SystemMetric extends BaseModel {
  metric_name: string
  metric_value: number
  metric_type: string
  tags?: Record<string, string>
}

export interface ProxyLog extends BaseModel {
  method: string
  url: string
  status_code: number
  response_time: number
  user_agent?: string
  ip_address?: string
  error_message?: string
}



export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 表格相关类型
export interface TableColumn {
  prop: string
  label: string
  width?: string | number
  minWidth?: string | number
  sortable?: boolean
  formatter?: (row: any, column: any, cellValue: any) => string
}

export interface TablePagination {
  page: number
  pageSize: number
  total: number
}

// 表单相关类型
export interface FormRule {
  required?: boolean
  message?: string
  trigger?: string | string[]
  min?: number
  max?: number
  pattern?: RegExp
  validator?: (rule: any, value: any, callback: any) => void
}

export type FormRules = Record<string, FormRule[]>

// 菜单类型
export interface MenuItem {
  id: string
  title: string
  icon?: string
  path?: string
  children?: MenuItem[]
  permission?: string
}

// 路由元信息
export interface RouteMeta {
  title?: string
  icon?: string
  requiresAuth?: boolean
  permissions?: string[]
  keepAlive?: boolean
}

// 通用选项类型
export interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
}

// 统计数据类型
export interface DashboardStats {
  total_users: number
  active_proxies: number
  total_requests: number
  success_rate: number
}

export interface ChartData {
  labels: string[]
  datasets: {
    label: string
    data: number[]
    backgroundColor?: string
    borderColor?: string
  }[]
}