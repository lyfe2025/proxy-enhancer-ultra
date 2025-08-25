// 用户相关类型
export interface User {
  id: string
  username: string
  email: string
  nickname?: string
  role_id: string
  role_name?: string
  status: 'active' | 'inactive' | 'locked'
  is_active: boolean
  last_login?: string
  created_at: string
  updated_at: string
  avatar?: string
  phone?: string
  department?: string
  position?: string
}

// 角色相关类型
export interface Role {
  id: string
  name: string
  code: string
  description?: string
  status: 'active' | 'inactive'
  permissions?: Permission[]
  permission_ids?: string[]
  user_count?: number
  created_at: string
  updated_at: string
  sort?: number
}

// 权限相关类型
export interface Permission {
  id: string
  name: string
  code: string
  type: 'menu' | 'button' | 'api'
  module: string
  description?: string
  status: 'active' | 'inactive'
  parent_id?: string
  path?: string
  icon?: string
  sort: number
  created_at: string
  updated_at: string
  children?: Permission[]
}

// 系统配置相关类型
export interface SystemConfig {
  system_name: string
  system_version: string
  system_logo?: string
  system_description?: string
  session_timeout: number
  max_login_attempts: number
  password_min_length: number
  password_complexity: boolean
  enable_2fa: boolean
  enable_captcha: boolean
  captcha_threshold?: number
  enable_email_verification?: boolean
  enable_sms_verification?: boolean
  maintenance_mode: boolean
  maintenance_message?: string
  allow_registration?: boolean
  default_role_id?: string
  timezone?: string
  date_format?: string
  language?: string
  theme?: 'light' | 'dark' | 'auto'
  log_level: 'debug' | 'info' | 'warn' | 'error'
  log_retention_days: number
  enable_access_log: boolean
  enable_error_log: boolean
  enable_audit_log: boolean
  backup_enabled?: boolean
  backup_frequency?: 'daily' | 'weekly' | 'monthly'
  backup_retention_days?: number
  email_host?: string
  email_port?: number
  email_username?: string
  email_password?: string
  email_from?: string
  email_from_name?: string
  email_secure?: boolean
  sms_enabled?: boolean
  sms_provider?: string
  sms_access_key?: string
  sms_secret_key?: string
  sms_signature?: string
  sms_api_key?: string
  sms_api_secret?: string
  storage_type?: string
  storage_bucket?: string
  storage_domain?: string
  storage_access_key?: string
  storage_secret_key?: string
  monitoring?: {
    enabled: boolean
    collection_interval: number
    retention_days: number
    cpu_threshold: number
    memory_threshold: number
    disk_threshold: number
  }
}

// 系统信息相关类型
export interface SystemInfo {
  name: string
  version: string
  go_version: string
  build_time: string
  git_commit: string
  uptime: string
  cpu_usage: number
  memory_usage: number
  memory_total: number
  memory_used: number
  disk_usage: number
  disk_total: number
  disk_used: number
  db_status: string
  db_version: string
  redis_status?: string
  redis_version?: string
  online_users: number
  total_users: number
  total_roles: number
  total_permissions: number
  system_load: number[]
  network_in: number
  network_out: number
  last_backup?: string
  next_backup?: string
}

// 用户表单类型
export interface UserForm {
  id?: string
  username: string
  email: string
  nickname?: string
  password?: string
  confirm_password?: string
  role_id: string
  status: 'active' | 'inactive'
  is_active: boolean
  phone?: string
  department?: string
  position?: string
  avatar?: string
}

// 角色表单类型
export interface RoleForm {
  id?: string
  name: string
  code?: string
  description?: string
  permission_ids: string[]
  status: 'active' | 'inactive'
  sort?: number
}

// 权限表单类型
export interface PermissionForm {
  id?: string
  name: string
  code: string
  type: 'menu' | 'button' | 'api'
  module: string
  description?: string
  parent_id?: string
  path?: string
  icon?: string
  sort: number
  status: 'active' | 'inactive'
}

// 搜索和筛选类型
export interface UserSearchParams {
  keyword?: string
  role_id?: string
  status?: string
  department?: string
  page?: number
  page_size?: number
  sort_field?: string
  sort_order?: 'asc' | 'desc'
}

export interface RoleSearchParams {
  keyword?: string
  status?: string
  page?: number
  page_size?: number
  sort_field?: string
  sort_order?: 'asc' | 'desc'
}

export interface PermissionSearchParams {
  keyword?: string
  module?: string
  type?: string
  status?: string
  parent_id?: string
  page?: number
  page_size?: number
  sort_field?: string
  sort_order?: 'asc' | 'desc'
}

// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  timestamp?: number
}

export interface PaginatedResponse<T = any> {
  data: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 操作日志类型
export interface OperationLog {
  id: string
  user_id: string
  username: string
  action: string
  resource: string
  resource_id?: string
  ip_address: string
  user_agent: string
  status: 'success' | 'failed'
  error_message?: string
  created_at: string
  duration?: number
}

// 登录日志类型
export interface LoginLog {
  id: string
  user_id: string
  username: string
  ip_address: string
  user_agent: string
  login_type: 'password' | '2fa' | 'sso'
  status: 'success' | 'failed'
  failure_reason?: string
  location?: string
  device?: string
  created_at: string
}

// 系统统计类型
export interface SystemStats {
  users: {
    total: number
    active: number
    inactive: number
    locked: number
    new_today: number
    new_this_week: number
    new_this_month: number
  }
  roles: {
    total: number
    active: number
    inactive: number
  }
  permissions: {
    total: number
    by_module: Record<string, number>
    by_type: Record<string, number>
  }
  system: {
    uptime: string
    cpu_usage: number
    memory_usage: number
    disk_usage: number
    online_users: number
    requests_today: number
    errors_today: number
  }
}

// 导出选项类型
export interface ExportOptions {
  format: 'excel' | 'csv' | 'pdf'
  fields: string[]
  filters?: Record<string, any>
  filename?: string
}

// 导入选项类型
export interface ImportOptions {
  file: File
  format: 'excel' | 'csv'
  mapping: Record<string, string>
  skip_header: boolean
  update_existing: boolean
}

// 批量操作类型
export interface BatchOperation {
  action: 'delete' | 'activate' | 'deactivate' | 'assign_role' | 'reset_password'
  ids: string[]
  params?: Record<string, any>
}

// 权限检查类型
export interface PermissionCheck {
  resource: string
  action: string
  user_id?: string
  role_ids?: string[]
}

// 模块定义
export const MODULES = {
  USER: 'user',
  ROLE: 'role', 
  PERMISSION: 'permission',
  PROXY: 'proxy',
  RULE: 'rule',
  DATA: 'data',
  SYSTEM: 'system',
  LOG: 'log'
} as const

// 权限类型定义
export const PERMISSION_TYPES = {
  MENU: 'menu',
  BUTTON: 'button',
  API: 'api'
} as const

// 用户状态定义
export const USER_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  LOCKED: 'locked'
} as const

// 角色状态定义
export const ROLE_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive'
} as const

// 权限状态定义
export const PERMISSION_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive'
} as const