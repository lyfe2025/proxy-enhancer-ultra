/**
 * 代理管理相关常量定义
 */

// 代理类型选项
export const PROXY_TYPE_OPTIONS = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' }
] as const

// 代理状态选项
export const PROXY_STATUS_OPTIONS = [
  { label: '全部', value: '' },
  { label: '活跃', value: 'active' },
  { label: '未激活', value: 'inactive' },
  { label: '错误', value: 'error' },
  { label: '测试中', value: 'testing' }
] as const

// 分页配置
export const PAGINATION_CONFIG = {
  DEFAULT_PAGE: 1,
  DEFAULT_SIZE: 10,
  PAGE_SIZES: [10, 20, 50, 100]
}

// 端口范围
export const PORT_RANGE = {
  MIN: 1,
  MAX: 65535
} as const

// 默认代理表单数据
export const DEFAULT_PROXY_FORM = {
  name: '',
  type: 'http' as const,
  host: '',
  port: 8080,
  username: '',
  password: '',
  description: ''
}

// 默认搜索表单数据
export const DEFAULT_SEARCH_FORM = {
  keyword: '',
  type: '',
  status: ''
}

// 表单验证规则
export const PROXY_FORM_RULES = {
  name: [
    { required: true, message: '请输入代理名称', trigger: 'blur' },
    { min: 2, max: 50, message: '代理名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择代理类型', trigger: 'change' }
  ],
  host: [
    { required: true, message: '请输入主机地址', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (!value) {
          callback(new Error('请输入主机地址'))
          return
        }
        
        // 验证IP地址格式
        const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
        // 验证域名格式
        const domainRegex = /^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)?$/
        
        if (!ipRegex.test(value) && !domainRegex.test(value)) {
          callback(new Error('请输入有效的IP地址或域名'))
          return
        }
        
        callback()
      },
      trigger: 'blur'
    }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    {
      validator: (rule: any, value: number, callback: Function) => {
        if (!value) {
          callback(new Error('请输入端口号'))
          return
        }
        
        if (!Number.isInteger(value) || value < PORT_RANGE.MIN || value > PORT_RANGE.MAX) {
          callback(new Error(`端口号必须在 ${PORT_RANGE.MIN} 到 ${PORT_RANGE.MAX} 之间`))
          return
        }
        
        callback()
      },
      trigger: 'blur'
    }
  ],
  username: [
    { max: 100, message: '用户名长度不能超过 100 个字符', trigger: 'blur' }
  ],
  password: [
    { max: 100, message: '密码长度不能超过 100 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '描述长度不能超过 500 个字符', trigger: 'blur' }
  ]
}

// 代理测试超时时间（毫秒）
export const PROXY_TEST_TIMEOUT = 10000

// 刷新间隔（毫秒）
export const REFRESH_INTERVAL = 30000

// 批量操作最大数量
export const MAX_BATCH_SIZE = 50