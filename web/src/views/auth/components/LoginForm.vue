<template>
  <div class="login-form">
    <!-- Logo和标题 -->
    <div class="login-header">
      <div class="logo">
        <div class="logo-icon">
          <el-icon size="40"><Connection /></el-icon>
        </div>
        <h1 class="logo-text">代理增强器</h1>
      </div>
      <p class="login-subtitle">智能代理管理平台</p>
    </div>

    <!-- 登录表单 -->
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      class="login-form-content"
      size="large"
      @keyup.enter="handleSubmit"
    >
      <el-form-item prop="username">
        <el-input
          v-model="formData.username"
          placeholder="请输入用户名"
          prefix-icon="User"
          clearable
        />
      </el-form-item>

      <el-form-item prop="password">
        <el-input
          v-model="formData.password"
          type="password"
          placeholder="请输入密码"
          prefix-icon="Lock"
          show-password
          clearable
        />
      </el-form-item>

      <!-- 验证码 -->
      <el-form-item v-if="showCaptcha" prop="captcha">
        <div class="captcha-container">
          <el-input
            v-model="formData.captcha"
            placeholder="请输入验证码"
            prefix-icon="Picture"
            clearable
            class="captcha-input"
          />
          <div class="captcha-image" @click="refreshCaptcha">
            <img :src="captchaUrl" alt="验证码" />
          </div>
        </div>
      </el-form-item>

      <!-- 记住我 -->
      <el-form-item>
        <div class="login-options">
          <el-checkbox v-model="formData.remember">记住我</el-checkbox>
          <el-link type="primary" @click="handleForgotPassword">
            忘记密码？
          </el-link>
        </div>
      </el-form-item>

      <!-- 登录按钮 -->
      <el-form-item>
        <el-button
          type="primary"
          class="login-button"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ loading ? '登录中...' : '登录' }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Connection } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

interface LoginFormData {
  username: string
  password: string
  captcha: string
  remember: boolean
}

const props = defineProps<{
  loading: boolean
  showCaptcha: boolean
  captchaUrl: string
}>()

const emit = defineEmits<{
  submit: [data: LoginFormData]
  refreshCaptcha: []
  forgotPassword: []
}>()

const formRef = ref<FormInstance>()

const formData = reactive<LoginFormData>({
  username: '',
  password: '',
  captcha: '',
  remember: false
})

const formRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 4, message: '验证码长度为 4 个字符', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    emit('submit', { ...formData })
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const refreshCaptcha = () => {
  emit('refreshCaptcha')
}

const handleForgotPassword = () => {
  emit('forgotPassword')
}

// 暴露重置表单方法
const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(formData, {
    username: '',
    password: '',
    captcha: '',
    remember: false
  })
}

defineExpose({
  resetForm
})
</script>

<style scoped>
.login-form {
  width: 100%;
  max-width: 400px;
  background: var(--el-bg-color);
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  animation: fadeInUp 0.6s ease-out;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin-bottom: 8px;
}

.logo-icon {
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: var(--el-color-primary-light-9);
  border-radius: 50%;
}

.logo-text {
  font-size: 28px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin: 0;
  background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-subtitle {
  color: var(--el-text-color-regular);
  font-size: 14px;
  margin: 0;
}

.login-form-content {
  width: 100%;
}

.captcha-container {
  display: flex;
  gap: 8px;
  align-items: center;
}

.captcha-input {
  flex: 1;
}

.captcha-image {
  width: 100px;
  height: 40px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;
}

.captcha-image:hover {
  border-color: var(--el-color-primary);
  transform: scale(1.02);
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 深色主题适配 */
.dark .login-form {
  background: var(--el-bg-color-overlay);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-form {
    padding: 24px;
  }
  
  .logo-text {
    font-size: 24px;
  }
}
</style>
