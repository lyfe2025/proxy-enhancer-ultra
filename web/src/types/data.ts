// 弹窗配置类型
export interface Popup {
  id?: number
  name: string
  title: string
  description?: string
  type: 'form' | 'survey' | 'feedback' | 'subscription'
  status: 'active' | 'inactive' | 'draft'
  triggerRules: PopupTriggerRule[]
  formFields: PopupField[]
  fields: PopupField[] // 兼容性字段
  triggerType: string // 兼容性字段
  styling: PopupStyling
  settings: PopupSettings
  createdAt?: string
  updatedAt?: string
  viewCount?: number
  submitCount?: number
  conversionRate?: number
}

// 弹窗触发规则
export interface PopupTriggerRule {
  id?: number
  type: 'page_load' | 'time_delay' | 'scroll_percentage' | 'exit_intent' | 'click_element'
  value: string
  operator?: 'equals' | 'contains' | 'greater_than' | 'less_than'
}

// 表单字段
export interface PopupField {
  id?: string
  name: string
  label: string
  type: 'text' | 'email' | 'phone' | 'number' | 'textarea' | 'select' | 'radio' | 'checkbox' | 'date'
  required: boolean
  placeholder?: string
  options?: string[]
  validation?: FieldValidation
  order: number
}

// 字段验证规则
export interface FieldValidation {
  minLength?: number
  maxLength?: number
  pattern?: string
  min?: number
  max?: number
  customMessage?: string
}

// 弹窗样式
export interface PopupStyling {
  width?: string
  height?: string
  position: 'center' | 'top' | 'bottom' | 'left' | 'right'
  backgroundColor?: string
  textColor?: string
  borderRadius?: string
  shadow?: boolean
  overlay?: boolean
  overlayColor?: string
}

// 弹窗设置
export interface PopupSettings {
  showCloseButton: boolean
  closeOnOverlayClick: boolean
  autoClose?: number
  showOnce?: boolean
  cookieExpiry?: number
  mobileOptimized: boolean
}

// 提交数据
export interface Submission {
  id?: number
  popupId: number
  popupName: string
  data: Record<string, any>
  userAgent?: string
  ipAddress?: string
  referrer?: string
  submittedAt: string
  processed?: boolean
  tags?: string[]
}

// 弹窗创建请求
export interface PopupCreateRequest {
  name: string
  title: string
  description?: string
  type: 'form' | 'survey' | 'feedback' | 'subscription'
  status: 'active' | 'inactive' | 'draft'
  triggerRules: PopupTriggerRule[]
  formFields: PopupField[]
  fields: PopupField[] // 兼容性字段
  triggerType: string // 兼容性字段
  styling: PopupStyling
  settings: PopupSettings
}

// 弹窗更新请求
export interface PopupUpdateRequest {
  name?: string
  title?: string
  description?: string
  type?: 'form' | 'survey' | 'feedback' | 'subscription'
  status?: 'active' | 'inactive' | 'draft'
  triggerRules?: PopupTriggerRule[]
  formFields?: PopupField[]
  fields?: PopupField[] // 兼容性字段
  triggerType?: string // 兼容性字段
  styling?: PopupStyling
  settings?: PopupSettings
}
