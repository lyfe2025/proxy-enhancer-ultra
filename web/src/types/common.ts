// 通用类型定义

// 分页响应类型
export interface PageResponse<T = any> {
  list: T[]
  total: number
  page: number
  limit: number
  pages: number
}

// 分页请求类型
export interface PageRequest {
  page?: number
  limit?: number
  search?: string
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

// 统计数据类型
export interface StatisticsData {
  total: number
  today: number
  yesterday: number
  thisWeek: number
  lastWeek: number
  thisMonth: number
  lastMonth: number
  growth: {
    today: number
    week: number
    month: number
  }
}

// 图表数据类型
export interface ChartData {
  labels: string[]
  datasets: {
    label: string
    data: number[]
    backgroundColor?: string | string[]
    borderColor?: string | string[]
    borderWidth?: number
  }[]
}

// 表单验证规则类型
export interface FormRule {
  required?: boolean
  message?: string
  trigger?: string | string[]
  min?: number
  max?: number
  pattern?: RegExp
  validator?: (rule: any, value: any, callback: any) => void
}

// 通知类型
export interface Notification {
  id: string
  type: 'info' | 'success' | 'warning' | 'error'
  title: string
  content: string
  data?: Record<string, any>
  userId?: string
  read: boolean
  readTime?: string
  createdAt: string
  expireAt?: string
}
