import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginResponse } from '@/types/auth'
import * as authApi from '@/api/auth'
import type { ApiLoginResponse } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const isLoading = ref(false)
  const permissions = ref<string[]>([])

  // 计算属性
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const userDisplayName = computed(() => {
    if (!user.value) return ''
    return user.value.nickname || user.value.username || user.value.email
  })

  // 登录
  const login = async (loginForm: authApi.LoginRequest) => {
    try {
      isLoading.value = true
      const response = await authApi.login(loginForm)
      
      console.log('登录响应:', response) // 调试日志
      
      // 直接访问 response 的属性，响应拦截器已经返回了 ApiResponse 格式
      if ((response as any).success) {
        const apiResponse = response as any
        const loginData = apiResponse.data as ApiLoginResponse
        token.value = loginData.token
        
        // 转换用户数据以匹配我们的 User 类型
        const userData = loginData.user
        user.value = {
          id: userData.id, // 现在是string类型的UUID
          username: userData.username,
          email: userData.email,
          nickname: undefined, // 后端没有返回nickname字段
          avatar: undefined, // 后端没有返回avatar字段
          phone: undefined, // 后端没有返回phone字段
          role: userData.role, // 使用role字段而不是roleName
          status: userData.enabled ? 'active' : 'inactive', // 根据enabled字段设置状态
          permissions: [], // 后端没有返回permissions字段，设为空数组
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          lastLoginAt: new Date().toISOString() // 设置为当前登录时间
        }
        permissions.value = [] // 后端没有返回permissions，设为空数组
        
        // 存储token和用户信息
        localStorage.setItem('token', token.value!)
        localStorage.setItem('user', JSON.stringify(user.value))
        localStorage.setItem('permissions', JSON.stringify([]))
        
        console.log('登录成功，用户信息:', user.value) // 调试日志
        console.log('用户角色:', user.value.role, '是否为admin:', user.value.role === 'admin') // 调试日志
        return { success: true }
      } else {
        const apiResponse = response as any
        console.error('登录失败:', apiResponse.message) // 调试日志
        return { success: false, message: apiResponse.message }
      }
    } catch (error: any) {
      console.error('登录错误:', error)
      const message = error.response?.data?.message || error.message || '登录失败，请稍后重试'
      return { success: false, message }
    } finally {
      isLoading.value = false
    }
  }

  // 注册
  const register = async (registerForm: authApi.RegisterRequest) => {
    try {
      isLoading.value = true
      const response = await authApi.register(registerForm)
      
      if (response.data.success) {
        return { success: true }
      } else {
        return { success: false, message: response.data.message }
      }
    } catch (error: any) {
      console.error('注册错误:', error)
      const message = error.response?.data?.message || error.message || '注册失败，请稍后重试'
      return { success: false, message }
    } finally {
      isLoading.value = false
    }
  }

  // 登出
  const logout = async () => {
    try {
      if (token.value) {
        await authApi.logout()
      }
    } catch (error) {
      console.error('登出API调用失败:', error)
    } finally {
      // 清除本地状态
      token.value = null
      user.value = null
      permissions.value = []
      
      // 清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('permissions')
      
      // 注意：组件层需要处理退出成功提示和页面跳转
    }
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      if (!token.value) return false
      
      const response = await authApi.getUserInfo()
      
      if (response.data.success) {
        // 转换用户数据以匹配我们的 User 类型
        const userData = response.data.data.user
        user.value = {
          id: userData.id, // 现在是string类型的UUID
          username: userData.username,
          email: userData.email,
          nickname: undefined, // 后端没有返回nickname字段
          avatar: undefined, // 后端没有返回avatar字段
          phone: undefined, // 后端没有返回phone字段
          role: userData.role, // 使用role字段而不是roleName
          status: userData.enabled ? 'active' : 'inactive', // 根据enabled字段设置状态
          permissions: [], // 后端没有返回permissions字段，设为空数组
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          lastLoginAt: new Date().toISOString() // 设置为当前时间
        }
        permissions.value = [] // 后端没有返回permissions，设为空数组
        
        // 更新本地存储
        localStorage.setItem('user', JSON.stringify(user.value))
        localStorage.setItem('permissions', JSON.stringify([]))
        
        return true
      } else {
        // token可能已过期，清除登录状态
        await logout()
        return false
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      await logout()
      return false
    }
  }

  // 更新用户信息
  const updateUserInfo = async (userInfo: Partial<User>) => {
    try {
      isLoading.value = true
      const response = await authApi.updateUserInfo(userInfo)
      
      if (response.data.success) {
        user.value = { ...user.value!, ...response.data.data }
        localStorage.setItem('user', JSON.stringify(user.value))
        return { success: true }
      } else {
        return { success: false, message: response.data.message }
      }
    } catch (error: any) {
      console.error('更新用户信息失败:', error)
      const message = error.response?.data?.message || error.message || '更新失败，请稍后重试'
      return { success: false, message }
    } finally {
      isLoading.value = false
    }
  }

  // 修改密码
  const changePassword = async (oldPassword: string, newPassword: string) => {
    try {
      isLoading.value = true
      const response = await authApi.changePassword({ oldPassword, newPassword, confirmPassword: newPassword })
      
      if (response.data.success) {
        await logout()
        return { success: true, message: '密码修改成功，请重新登录' }
      } else {
        return { success: false, message: response.data.message }
      }
    } catch (error: any) {
      console.error('修改密码失败:', error)
      const message = error.response?.data?.message || error.message || '密码修改失败，请稍后重试'
      return { success: false, message }
    } finally {
      isLoading.value = false
    }
  }

  // 忘记密码
  const forgotPassword = async (email: string) => {
    try {
      isLoading.value = true
      const response = await authApi.forgotPassword(email)
      
      if (response.data.success) {
        return { success: true, message: '重置邮件已发送，请查收邮箱' }
      } else {
        return { success: false, message: response.data.message }
      }
    } catch (error: any) {
      console.error('发送重置邮件失败:', error)
      const message = error.response?.data?.message || error.message || '发送失败，请稍后重试'
      return { success: false, message }
    } finally {
      isLoading.value = false
    }
  }

  // 检查权限
  const hasPermission = (permission: string) => {
    // admin 角色拥有所有权限
    if (user.value?.role === 'admin') {
      return true
    }
    if (!permissions.value.length) return false
    return permissions.value.includes(permission) || permissions.value.includes('*')
  }

  // 检查多个权限（任一满足）
  const hasAnyPermission = (permissionList: string[]) => {
    return permissionList.some(permission => hasPermission(permission))
  }

  // 检查多个权限（全部满足）
  const hasAllPermissions = (permissionList: string[]) => {
    return permissionList.every(permission => hasPermission(permission))
  }

  // 检查角色
  const hasRole = (role: string) => {
    if (!user.value) return false
    return user.value.role === role || user.value.role === 'admin'
  }

  // 检查多个角色（任一满足）
  const hasAnyRole = (roleList: string[]) => {
    return roleList.some(role => hasRole(role))
  }

  // 检查多个角色（全部满足）
  const hasAllRoles = (roleList: string[]) => {
    return roleList.every(role => hasRole(role))
  }

  // 初始化认证状态
  const initAuth = async () => {
    const savedToken = localStorage.getItem('token')
    const savedUser = localStorage.getItem('user')
    const savedPermissions = localStorage.getItem('permissions')
    
    if (savedToken && savedUser) {
      token.value = savedToken
      try {
        user.value = JSON.parse(savedUser)
        permissions.value = savedPermissions ? JSON.parse(savedPermissions) : []
        
        // 验证token是否仍然有效
        const isValid = await fetchUserInfo()
        if (!isValid) {
          // token无效，清除状态
          token.value = null
          user.value = null
          permissions.value = []
        }
      } catch (error) {
        console.error('解析本地用户数据失败:', error)
        // 清除无效数据
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        localStorage.removeItem('permissions')
      }
    }
  }

  return {
    // 状态
    user,
    token,
    isLoading,
    permissions,
    
    // 计算属性
    isAuthenticated,
    isAdmin,
    userDisplayName,
    
    // 方法
    login,
    register,
    logout,
    fetchUserInfo,
    updateUserInfo,
    changePassword,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    hasRole,
    hasAnyRole,
    hasAllRoles,
    initAuth,
    forgotPassword
  }
})