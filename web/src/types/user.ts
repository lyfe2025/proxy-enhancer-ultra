// 用户相关类型定义

// 用户类型
export interface User {
  id: string
  username: string
  email: string
  name: string
  phone?: string
  avatar?: string
  department?: string
  status: 'active' | 'inactive' | 'locked'
  roles: Role[]
  permissions: Permission[]
  lastLoginTime?: string
  lastLoginIp?: string
  loginCount: number
  createdAt: string
  updatedAt: string
  createdBy?: string
  updatedBy?: string
}

// 角色类型
export interface Role {
  id: string
  name: string
  code: string
  description?: string
  type: 'system' | 'custom'
  status: 'active' | 'inactive'
  permissions: Permission[]
  userCount: number
  createdAt: string
  updatedAt: string
  createdBy?: string
  updatedBy?: string
}

// 权限类型
export interface Permission {
  id: string
  name: string
  code: string
  type: 'menu' | 'button' | 'api' | 'data'
  module: string
  description?: string
  parentId?: string
  path?: string
  method?: string
  icon?: string
  sort: number
  children?: Permission[]
  createdAt: string
  updatedAt: string
}

// 登录历史类型
export interface LoginHistory {
  id: string
  userId: string
  username: string
  ip: string
  userAgent: string
  location?: string
  device?: string
  browser?: string
  os?: string
  status: 'success' | 'failed'
  failReason?: string
  timestamp: string
}

// 在线用户类型
export interface OnlineUser {
  id: string
  userId: string
  username: string
  name: string
  avatar?: string
  ip: string
  location?: string
  device?: string
  browser?: string
  os?: string
  loginTime: string
  lastActiveTime: string
  sessionId: string
}

// 部门类型
export interface Department {
  id: string
  name: string
  code: string
  description?: string
  parentId?: string
  level: number
  sort: number
  manager?: string
  phone?: string
  email?: string
  address?: string
  status: 'active' | 'inactive'
  children?: Department[]
  userCount: number
  createdAt: string
  updatedAt: string
}
