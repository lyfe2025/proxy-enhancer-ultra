// 代理配置类型
export interface Proxy {
  id?: number
  name: string
  type: 'http' | 'https' | 'socks5'
  host: string
  port: number
  username?: string
  password?: string
  status: 'active' | 'inactive' | 'error'
  description?: string
  createdAt?: string
  updatedAt?: string
  lastUsed?: string
  responseTime?: number
  successRate?: number
  totalRequests?: number
  failedRequests?: number
}

// 代理创建请求
export interface ProxyCreateRequest {
  name: string
  type: 'http' | 'https' | 'socks5'
  host: string
  port: number
  username?: string
  password?: string
  description?: string
}

// 代理更新请求
export interface ProxyUpdateRequest {
  name?: string
  type?: 'http' | 'https' | 'socks5'
  host?: string
  port?: number
  username?: string
  password?: string
  description?: string
}
