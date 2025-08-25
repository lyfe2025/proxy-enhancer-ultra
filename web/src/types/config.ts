// 系统配置相关类型定义

// 系统配置类型
export interface SystemConfig {
  // 基本设置
  basic: {
    siteName: string
    siteDescription?: string
    siteLogo?: string
    siteFavicon?: string
    siteKeywords?: string
    timezone: string
    language: string
    dateFormat: string
    timeFormat: string
  }
  
  // 安全设置
  security: {
    passwordPolicy: {
      minLength: number
      requireUppercase: boolean
      requireLowercase: boolean
      requireNumbers: boolean
      requireSpecialChars: boolean
      maxAge: number // 密码最大有效期（天）
    }
    sessionTimeout: number // 会话超时时间（分钟）
    maxLoginAttempts: number // 最大登录尝试次数
    lockoutDuration: number // 锁定持续时间（分钟）
    twoFactorAuth: boolean // 是否启用双因子认证
    ipWhitelist: string[] // IP白名单
    ipBlacklist: string[] // IP黑名单
    allowedFileTypes: string[] // 允许上传的文件类型
    maxFileSize: number // 最大文件大小（MB）
  }
  
  // 邮件设置
  email: {
    enabled: boolean
    smtp: {
      host: string
      port: number
      username: string
      password: string
      secure: boolean
      fromName: string
      fromEmail: string
    }
    templates: {
      welcome: string
      resetPassword: string
      notification: string
    }
  }
  
  // 短信设置
  sms: {
    enabled: boolean
    provider: string
    config: Record<string, any>
  }
  
  // Webhook设置
  webhook: {
    enabled: boolean
    urls: string[]
    events: string[]
    secret?: string
  }
  
  // 存储设置
  storage: {
    type: 'local' | 'oss' | 's3' | 'qiniu'
    config: Record<string, any>
    maxSize: number
    allowedTypes: string[]
  }
  
  // 缓存设置
  cache: {
    type: 'memory' | 'redis'
    config: Record<string, any>
    ttl: number
  }
  
  // 日志设置
  logging: {
    level: 'debug' | 'info' | 'warn' | 'error'
    maxSize: number // 单个日志文件最大大小（MB）
    maxFiles: number // 最大日志文件数量
    retention: number // 日志保留天数
    modules: string[] // 启用日志的模块
  }
  
  // 备份设置
  backup: {
    enabled: boolean
    schedule: string // cron表达式
    retention: number // 备份保留天数
    storage: {
      type: 'local' | 'oss' | 's3'
      config: Record<string, any>
    }
  }
  
  // 监控设置
  monitoring: {
    enabled: boolean
    metrics: {
      cpu: boolean
      memory: boolean
      disk: boolean
      network: boolean
      database: boolean
    }
    alerts: {
      enabled: boolean
      thresholds: {
        cpu: number
        memory: number
        disk: number
        errorRate: number
      }
      channels: string[]
    }
  }
}

// 系统许可证类型
export interface SystemLicense {
  id: string
  type: 'trial' | 'standard' | 'professional' | 'enterprise'
  status: 'active' | 'expired' | 'invalid'
  features: string[]
  limits: {
    users: number
    storage: number
    bandwidth: number
  }
  issueDate: string
  expireDate: string
  company: string
  contact: string
}

// 系统更新类型
export interface SystemUpdate {
  id: string
  version: string
  title: string
  description: string
  type: 'major' | 'minor' | 'patch' | 'hotfix'
  size: number
  downloadUrl: string
  checksum: string
  releaseDate: string
  changelog: string[]
  requirements: {
    minVersion: string
    dependencies: string[]
  }
  status: 'available' | 'downloading' | 'downloaded' | 'installing' | 'installed' | 'failed'
}
