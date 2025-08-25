// 日志相关类型定义

// 系统日志类型
export interface SystemLog {
  id: string
  level: 'debug' | 'info' | 'warn' | 'error'
  module: string
  message: string
  data?: Record<string, any>
  timestamp: string
  ip?: string
  userAgent?: string
  userId?: string
  requestId?: string
}

// 操作日志类型
export interface OperationLog {
  id: string
  userId: string
  username: string
  action: string
  module: string
  description: string
  ip: string
  userAgent: string
  requestId: string
  params?: Record<string, any>
  result?: Record<string, any>
  duration: number
  status: 'success' | 'failed'
  timestamp: string
}
