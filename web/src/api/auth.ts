import { api } from './request'
import type { LoginRequest, LoginResponse, RegisterRequest, User, ApiResponse } from '@/types'

// 认证相关API
export const authApi = {
  // 用户登录
  login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return api.post('/auth/login', data)
  },
  
  // 用户注册
  register(data: RegisterRequest): Promise<ApiResponse<User>> {
    return api.post('/auth/register', data)
  },
  
  // 获取当前用户信息
  getCurrentUser(): Promise<ApiResponse<User>> {
    return api.get('/auth/me')
  },
  
  // 刷新token
  refreshToken(): Promise<ApiResponse<{ token: string; expires_at: string }>> {
    return api.post('/auth/refresh')
  },
  
  // 用户登出
  logout(): Promise<ApiResponse<null>> {
    return api.post('/auth/logout')
  },
  
  // 修改密码
  changePassword(data: { old_password: string; new_password: string }): Promise<ApiResponse<null>> {
    return api.post('/auth/change-password', data)
  },
  
  // 重置密码
  resetPassword(data: { email: string }): Promise<ApiResponse<null>> {
    return api.post('/auth/reset-password', data)
  },
  
  // 验证token
  verifyToken(): Promise<ApiResponse<{ valid: boolean }>> {
    return api.get('/auth/verify')
  }
}