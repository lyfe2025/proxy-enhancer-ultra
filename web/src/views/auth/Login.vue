<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <LoginBackground />

    <!-- 登录表单 -->
    <div class="login-form-container">
      <LoginForm
        ref="loginFormRef"
        :loading="loading"
        :show-captcha="showCaptcha"
        :captcha-url="captchaUrl"
        @submit="handleLogin"
        @refresh-captcha="refreshCaptcha"
        @forgot-password="showForgotPassword = true"
      />

      <!-- 社交登录 -->
      <SocialLogin @social-login="handleSocialLogin" />
    </div>

    <!-- 忘记密码对话框 -->
    <ForgotPasswordDialog
      v-model="showForgotPassword"
      :loading="forgotPasswordLoading"
      @submit="handleForgotPassword"
      @close="showForgotPassword = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

// 导入子组件
import LoginBackground from './components/LoginBackground.vue'
import LoginForm from './components/LoginForm.vue'
import SocialLogin from './components/SocialLogin.vue'
import ForgotPasswordDialog from './components/ForgotPasswordDialog.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 组件引用
const loginFormRef = ref()

// 状态
const loading = ref(false)
const showCaptcha = ref(false)
const captchaUrl = ref('')
const showForgotPassword = ref(false)
const forgotPasswordLoading = ref(false)

// 处理登录
const handleLogin = async (formData: any) => {
  loading.value = true
  
  try {
    await authStore.login({
      username: formData.username,
      password: formData.password,
      captcha: showCaptcha.value ? formData.captcha : undefined,
      rememberMe: formData.remember
    })
    
    ElMessage.success('登录成功')
    
    // 跳转到目标页面或仪表盘
    const redirect = route.query.redirect as string
    router.push(redirect || '/dashboard')
  } catch (error: any) {
    console.error('登录失败:', error)
    
    // 如果需要验证码，显示验证码
    if (error.code === 'CAPTCHA_REQUIRED') {
      showCaptcha.value = true
      refreshCaptcha()
      ElMessage.warning('请输入验证码')
    } else if (error.code === 'INVALID_CREDENTIALS') {
      ElMessage.error('用户名或密码错误')
    } else if (error.code === 'ACCOUNT_LOCKED') {
      ElMessage.error('账户已被锁定，请联系管理员')
    } else if (error.code === 'ACCOUNT_DISABLED') {
      ElMessage.error('账户已被禁用，请联系管理员')
    } else {
      ElMessage.error(error.message || '登录失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}

// 刷新验证码
const refreshCaptcha = async () => {
  try {
    // 实际实现应该调用获取验证码的API
    captchaUrl.value = `/api/captcha?t=${Date.now()}`
  } catch (error) {
    console.error('刷新验证码失败:', error)
    ElMessage.error('刷新验证码失败')
  }
}

// 处理忘记密码
const handleForgotPassword = async (formData: any) => {
  forgotPasswordLoading.value = true
  
  try {
    // 实际实现应该调用忘记密码的API
    await new Promise(resolve => setTimeout(resolve, 2000)) // 模拟API调用
    
    ElMessage.success(`重置邮件已发送到 ${formData.email}`)
    showForgotPassword.value = false
  } catch (error: any) {
    console.error('发送重置邮件失败:', error)
    ElMessage.error(error.message || '发送重置邮件失败，请稍后重试')
  } finally {
    forgotPasswordLoading.value = false
  }
}

// 处理社交登录
const handleSocialLogin = (provider: string) => {
  console.log('社交登录:', provider)
  // 实际实现应该跳转到对应的OAuth授权页面
  ElMessage.info(`${provider}登录功能开发中...`)
}

// 页面挂载时检查是否需要显示验证码
const checkCaptchaRequirement = () => {
  // 可以根据登录失败次数或其他条件决定是否显示验证码
  const failedAttempts = localStorage.getItem('loginFailedAttempts')
  if (failedAttempts && parseInt(failedAttempts) >= 3) {
    showCaptcha.value = true
    refreshCaptcha()
  }
}

// 初始化
checkCaptchaRequirement()
</script>

<style scoped>
.login-container {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(
    135deg,
    #667eea 0%,
    #764ba2 100%
  );
  padding: 20px;
  overflow: hidden;
}

.login-form-container {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 400px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-container {
    padding: 16px;
  }
}
</style>