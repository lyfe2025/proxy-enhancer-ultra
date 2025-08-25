// 系统管理相关类型定义

// 用户类型
export interface User {
  id: string
  username: string
  email: string
  nickname?: string
  avatar?: string
  status: 'active' | 'inactive' | 'locked'
  is_active: boolean
  role_id: string
  role?: {
    id: string
    name: string
    code: string
    description?: string
  }
  permissions?: string[]
  phone?: string
  department?: string
  position?: string
  last_login_at?: string
  last_login_ip?: string
  created_at: string
  updated_at: string
  statusLoading?: boolean // 用于状态切换时的loading状态
}

// 角色类型
export interface Role {
  id: string
  name: string
  code: string
  description?: string
  status: 'active' | 'inactive'
  permissions?: Permission[]
  created_at: string
  updated_at: string
  sort?: number
  user_count?: number
}

// 权限类型
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

// 系统配置类型
export interface SystemConfig {
  system_name: string
  system_version: string
  system_logo?: string
  system_description?: string
  session_timeout: number
  max_login_attempts: number
  password_min_length: number
  password_complexity: string[]
  enable_2fa: boolean
  enable_captcha: boolean
  captcha_threshold: number
  enable_email_verification: boolean
  enable_sms_verification: boolean
  maintenance_mode: boolean
  maintenance_message?: string
  allow_registration: boolean
  default_role_id?: string
  timezone: string
  date_format: string
  language: string
  theme: 'light' | 'dark' | 'auto'
  log_level: 'debug' | 'info' | 'warn' | 'error'
  log_retention_days: number
  enable_access_log: boolean
  enable_error_log: boolean
  enable_audit_log: boolean
  backup_enabled: boolean
  backup_frequency: 'daily' | 'weekly' | 'monthly'
  backup_retention_days: number
  email_host?: string
  email_port?: number
  email_username?: string
  email_password?: string
  email_from?: string
  sms_provider?: string
  sms_api_key?: string
  sms_api_secret?: string
  // 嵌套结构
  basic?: {
    system_name: string
    system_version: string
    system_description?: string
    maintenance_mode: boolean
    maintenance_message?: string
  }
  security?: {
    password_min_length: number
    password_complexity: string[]
    login_attempts: number
    lockout_duration: number
    session_timeout: number
    two_factor_auth: boolean
  }
  email?: {
    enabled: boolean
    smtp_host?: string
    smtp_port?: number
    smtp_username?: string
    smtp_password?: string
    from_email?: string
    from_name?: string
    use_ssl: boolean
  }
  sms?: {
    enabled: boolean
    provider?: string
    access_key?: string
    secret_key?: string
    signature?: string
  }
  storage?: {
    type: string
    max_file_size: number
    allowed_extensions: string[]
  }
  logging?: {
    level: string
    retention_days: number
    file_enabled: boolean
    database_enabled: boolean
  }
  monitoring?: {
    enabled: boolean
    collection_interval: number
    retention_days: number
    cpu_threshold: number
    memory_threshold: number
    disk_threshold: number
  }
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

// 系统统计信息
export interface SystemStats {
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

// 系统健康状态
export interface SystemHealth {
  status: 'healthy' | 'warning' | 'error'
  database: 'connected' | 'disconnected'
  redis: 'connected' | 'disconnected'
  services: Array<{
    name: string
    status: 'running' | 'stopped' | 'error'
    uptime?: string
  }>
}

// 文件上传响应
export interface FileUploadResponse {
  url: string
  filename: string
  size: number
  type: string
}

// 系统备份信息
export interface BackupInfo {
  filename: string
  size: number
  createdAt: string
}

// 表单数据类型
export interface UserFormData {
  id?: string
  username: string
  email: string
  password?: string
  confirm_password?: string
  nickname?: string
  avatar?: string
  role_id: string
  status: 'active' | 'inactive' | 'locked'
  is_active: boolean
  phone?: string
  department?: string
  position?: string
}

export interface RoleFormData {
  id?: string
  name: string
  code: string
  description?: string
  status: 'active' | 'inactive'
  permission_ids: string[]
  sort?: number
}

export interface PermissionFormData {
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
