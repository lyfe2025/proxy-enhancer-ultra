// 用户类型
export interface User {
  id: number
  username: string
  email: string
  nickname?: string
  avatar?: string
  phone?: string
  role: string
  status: 'active' | 'inactive' | 'banned'
  permissions: string[]
  createdAt: string
  updatedAt: string
  lastLoginAt?: string
}

// 登录响应类型
export interface LoginResponse {
  user: User
  token: string
  refreshToken?: string
  expiresIn: number
  permissions: string[]
}

// 登录请求类型
export interface LoginRequest {
  username: string
  password: string
  captcha?: string
  rememberMe?: boolean
}

// 注册请求类型
export interface RegisterRequest {
  username: string
  email: string
  password: string
  confirmPassword: string
  captcha?: string
}

// 修改密码请求类型
export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

// 用户信息更新请求类型
export interface UpdateProfileRequest {
  nickname?: string
  email?: string
  avatar?: string
  phone?: string
}
