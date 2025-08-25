import api from './index'
import type { ApiResponse } from './index'

// 登录请求参数
export interface LoginRequest {
  username: string
  password: string
  captcha?: string
  captchaId?: string
  rememberMe?: boolean
}

// 登录响应数据
export interface LoginResponse {
  token: string
  refreshToken: string
  expiresIn: number
  user: {
    id: number
    username: string
    email: string
    nickname?: string
    avatar?: string
    roleId: number
    roleName: string
    permissions: string[]
  }
}

// 注册请求参数
export interface RegisterRequest {
  username: string
  email: string
  password: string
  confirmPassword: string
  nickname?: string
  captcha?: string
  captchaId?: string
  inviteCode?: string
}

// 修改密码请求参数
export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

// 重置密码请求参数
export interface ResetPasswordRequest {
  email: string
  code: string
  newPassword: string
  confirmPassword: string
}

// 用户信息更新请求参数
export interface UpdateProfileRequest {
  nickname?: string
  email?: string
  avatar?: string
  phone?: string
}

// 验证码响应
export interface CaptchaResponse {
  captchaId: string
  captchaImage: string // base64 图片
}

// 登录
export const login = (data: LoginRequest) => {
  return api.post<ApiResponse<LoginResponse>>('/auth/login', data)
}

// 注册
export const register = (data: RegisterRequest) => {
  return api.post<ApiResponse<{ message: string }>>('/auth/register', data)
}

// 登出
export const logout = () => {
  return api.post<ApiResponse<void>>('/auth/logout')
}

// 刷新token
export const refreshToken = (refreshToken: string) => {
  return api.post<ApiResponse<{ token: string; expiresIn: number }>>('/auth/refresh', {
    refreshToken
  })
}

// 获取当前用户信息
export const getCurrentUser = () => {
  return api.get<ApiResponse<LoginResponse['user']>>('/auth/me')
}

// 更新用户信息
export const updateProfile = (data: UpdateProfileRequest) => {
  return api.put<ApiResponse<LoginResponse['user']>>('/auth/profile', data)
}

// 修改密码
export const changePassword = (data: ChangePasswordRequest) => {
  return api.patch<ApiResponse<void>>('/auth/password', data)
}

// 发送重置密码邮件
export const sendResetPasswordEmail = (email: string) => {
  return api.post<ApiResponse<{ message: string }>>('/auth/reset-password/send', { email })
}

// 重置密码
export const resetPassword = (data: ResetPasswordRequest) => {
  return api.post<ApiResponse<{ message: string }>>('/auth/reset-password', data)
}

// 获取验证码
export const getCaptcha = () => {
  return api.get<ApiResponse<CaptchaResponse>>('/auth/captcha')
}

// 验证邮箱
export const verifyEmail = (token: string) => {
  return api.post<ApiResponse<{ message: string }>>('/auth/verify-email', { token })
}

// 重新发送验证邮件
export const resendVerificationEmail = () => {
  return api.post<ApiResponse<{ message: string }>>('/auth/resend-verification')
}

// 检查用户名是否可用
export const checkUsernameAvailable = (username: string) => {
  return api.get<ApiResponse<{ available: boolean }>>('/auth/check-username', {
    params: { username }
  })
}

// 检查邮箱是否可用
export const checkEmailAvailable = (email: string) => {
  return api.get<ApiResponse<{ available: boolean }>>('/auth/check-email', {
    params: { email }
  })
}

// 绑定第三方账号
export const bindThirdPartyAccount = (provider: string, code: string) => {
  return api.post<ApiResponse<{ message: string }>>('/auth/bind', {
    provider,
    code
  })
}

// 解绑第三方账号
export const unbindThirdPartyAccount = (provider: string) => {
  return api.delete<ApiResponse<{ message: string }>>(`/auth/bind/${provider}`)
}

// 获取用户的第三方账号绑定状态
export const getThirdPartyBindings = () => {
  return api.get<ApiResponse<Array<{
    provider: string
    bound: boolean
    nickname?: string
    avatar?: string
    boundAt?: string
  }>>>('/auth/bindings')
}

// 启用/禁用两步验证
export const toggleTwoFactorAuth = (enabled: boolean, code?: string) => {
  return api.patch<ApiResponse<{
    enabled: boolean
    qrCode?: string // 启用时返回二维码
    backupCodes?: string[] // 启用时返回备用码
  }>>('/auth/2fa', { enabled, code })
}

// 验证两步验证码
export const verifyTwoFactorCode = (code: string) => {
  return api.post<ApiResponse<{ valid: boolean }>>('/auth/2fa/verify', { code })
}

// 获取用户信息
export const getUserInfo = () => {
  return api.get<ApiResponse<LoginResponse>>('/auth/user-info')
}

// 更新用户信息
export const updateUserInfo = (data: UpdateProfileRequest) => {
  return api.put<ApiResponse<LoginResponse['user']>>('/auth/user-info', data)
}

// 忘记密码
export const forgotPassword = (email: string) => {
  return api.post<ApiResponse<{ message: string }>>('/auth/forgot-password', { email })
}

// 创建 authApi 对象，包含所有相关函数
export const authApi = {
  login,
  register,
  logout,
  refreshToken,
  getCurrentUser,
  updateProfile,
  changePassword,
  sendResetPasswordEmail,
  resetPassword,
  getCaptcha,
  verifyEmail,
  resendVerificationEmail,
  checkUsernameAvailable,
  getUserInfo,
  updateUserInfo,
  forgotPassword
}