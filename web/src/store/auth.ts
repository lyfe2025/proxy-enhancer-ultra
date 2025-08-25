import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { authApi } from '@/api'
import type { User, LoginRequest, RegisterRequest } from '@/types'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)
  const loading = ref(false)
  
  // 计算属性
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userPermissions = computed(() => {
    // 这里可以根据用户角色返回权限列表
    // 暂时返回空数组，后续可以扩展
    return []
  })
  
  // 设置token
  const setToken = (newToken: string | null) => {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }
  
  // 设置用户信息
  const setUser = (newUser: User | null) => {
    user.value = newUser
  }
  
  // 登录
  const login = async (loginData: LoginRequest) => {
    try {
      loading.value = true
      const response = await authApi.login(loginData)
      
      if (response.success && response.data) {
        const { token: newToken, user: userData } = response.data
        setToken(newToken)
        setUser(userData)
        
        ElMessage.success('登录成功')
        
        // 跳转到首页或之前访问的页面
        const redirect = router.currentRoute.value.query.redirect as string
        await router.push(redirect || '/')
        
        return true
      }
      return false
    } catch (error) {
      console.error('登录失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }
  
  // 注册
  const register = async (registerData: RegisterRequest) => {
    try {
      loading.value = true
      const response = await authApi.register(registerData)
      
      if (response.success) {
        ElMessage.success('注册成功，请登录')
        return true
      }
      return false
    } catch (error) {
      console.error('注册失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }
  
  // 登出
  const logout = async () => {
    try {
      // 调用后端登出接口
      await authApi.logout()
    } catch (error) {
      console.error('登出请求失败:', error)
    } finally {
      // 无论后端请求是否成功，都清除本地状态
      setToken(null)
      setUser(null)
      
      // 跳转到登录页
      await router.push('/login')
      
      ElMessage.success('已退出登录')
    }
  }
  
  // 获取当前用户信息
  const fetchUserInfo = async () => {
    try {
      const response = await authApi.getCurrentUser()
      
      if (response.success && response.data) {
        setUser(response.data)
        return true
      }
      return false
    } catch (error) {
      console.error('获取用户信息失败:', error)
      return false
    }
  }
  
  // 检查认证状态
  const checkAuth = async () => {
    if (!token.value) {
      return false
    }
    
    try {
      // 验证token是否有效
      const verifyResponse = await authApi.verifyToken()
      
      if (verifyResponse.success && verifyResponse.data?.valid) {
        // token有效，获取用户信息
        if (!user.value) {
          await fetchUserInfo()
        }
        return true
      } else {
        // token无效，清除状态
        setToken(null)
        setUser(null)
        return false
      }
    } catch (error) {
      console.error('检查认证状态失败:', error)
      // 网络错误或其他错误，清除状态
      setToken(null)
      setUser(null)
      return false
    }
  }
  
  // 刷新token
  const refreshToken = async () => {
    try {
      const response = await authApi.refreshToken()
      
      if (response.success && response.data) {
        setToken(response.data.token)
        return true
      }
      return false
    } catch (error) {
      console.error('刷新token失败:', error)
      return false
    }
  }
  
  // 修改密码
  const changePassword = async (oldPassword: string, newPassword: string) => {
    try {
      loading.value = true
      const response = await authApi.changePassword({
        old_password: oldPassword,
        new_password: newPassword
      })
      
      if (response.success) {
        ElMessage.success('密码修改成功')
        return true
      }
      return false
    } catch (error) {
      console.error('修改密码失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }
  
  // 检查权限
  const hasPermission = (permission: string) => {
    // 这里可以根据用户权限进行检查
    // 暂时返回true，后续可以扩展
    return true
  }
  
  return {
    // 状态
    token: readonly(token),
    user: readonly(user),
    loading: readonly(loading),
    
    // 计算属性
    isAuthenticated,
    userPermissions,
    
    // 方法
    login,
    register,
    logout,
    fetchUserInfo,
    checkAuth,
    refreshToken,
    changePassword,
    hasPermission
  }
})