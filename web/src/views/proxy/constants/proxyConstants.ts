import type { FormItemRule } from 'element-plus'

/**
 * 代理管理相关常量配置
 */

/**
 * 代理类型选项
 */
export const PROXY_TYPE_OPTIONS = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' },
  { label: 'SOCKS4', value: 'socks4' }
] as const

/**
 * 代理状态选项
 */
export const PROXY_STATUS_OPTIONS = [
  { label: '活跃', value: 'active' },
  { label: '停用', value: 'inactive' },
  { label: '异常', value: 'error' }
] as const

/**
 * 分页配置
 */
export const PAGINATION_CONFIG = {
  DEFAULT_PAGE: 1,
  DEFAULT_SIZE: 20,
  PAGE_SIZES: [10, 20, 50, 100]
} as const

/**
 * 默认代理表单数据
 */
export const DEFAULT_PROXY_FORM = {
  name: '',
  type: 'http',
  host: '',
  port: 8080,
  username: '',
  password: '',
  description: ''
} as const

/**
 * 默认搜索表单数据
 */
export const DEFAULT_SEARCH_FORM = {
  keyword: '',
  type: '',
  status: ''
} as const

/**
 * 代理表单验证规则
 */
export const PROXY_FORM_RULES: Record<string, FormItemRule[]> = {
  name: [
    { required: true, message: '请输入代理名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择代理类型', trigger: 'change' }
  ],
  host: [
    { required: true, message: '请输入服务器地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: '端口号必须在1-65535之间', trigger: 'blur' }
  ]
}

/**
 * 端口号范围
 */
export const PORT_RANGE = {
  MIN: 1,
  MAX: 65535
} as const