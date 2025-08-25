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
  background: var(--bg-secondary);
  border-radius: var(--radius-xl);
  padding: 40px;
  box-shadow: var(--shadow-xl), var(--shadow-glow);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-secondary);
  animation: fadeInUp 0.6s ease-out;
  position: relative;
  overflow: hidden;
}

.login-form::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--primary-color), transparent);
  opacity: 0.3;
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
  color: var(--primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: var(--primary-alpha-10);
  border: 1px solid var(--primary-alpha-20);
  border-radius: var(--radius-full);
  transition: all var(--transition-normal);
}

.logo-icon:hover {
  background: var(--primary-alpha-20);
  box-shadow: var(--shadow-glow);
  transform: scale(1.05);
}

.logo-text {
  font-size: var(--font-3xl);
  font-weight: var(--font-bold);
  color: var(--text-primary);
  margin: 0;
  background: linear-gradient(135deg, var(--primary-color), #4ade80);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-family: 'Inter', 'SF Pro Display', system-ui, sans-serif;
  letter-spacing: -0.5px;
}

.login-subtitle {
  color: var(--text-secondary);
  font-size: var(--font-sm);
  margin: 0;
  font-weight: var(--font-medium);
}

.login-form-content {
  width: 100%;
}

.captcha-container {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.captcha-input {
  flex: 1;
}

.captcha-image {
  width: 100px;
  height: 40px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  overflow: hidden;
  transition: all var(--transition-normal);
}

.captcha-image:hover {
  border-color: var(--primary-color);
  transform: scale(1.02);
  box-shadow: 0 0 10px var(--primary-alpha-30);
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
  height: 48px;
  font-size: var(--font-md);
  font-weight: var(--font-semibold);
  background: linear-gradient(135deg, var(--primary-color), #4ade80);
  border: none;
  border-radius: var(--radius-lg);
  color: var(--text-inverse);
  cursor: pointer;
  transition: all var(--transition-normal);
  position: relative;
  overflow: hidden;
  font-family: 'Inter', system-ui, sans-serif;
  letter-spacing: 0.3px;
}

.login-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-glow-strong), var(--shadow-lg);
  background: linear-gradient(135deg, #00ff88, #66ffaa);
}

.login-button:hover::before {
  left: 100%;
}

.login-button:active {
  transform: translateY(-1px);
}

.login-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.login-button:disabled:hover {
  transform: none;
  box-shadow: none;
}

/* Element Plus 组件样式覆盖 */
:deep(.el-input__wrapper) {
  background-color: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  box-shadow: none;
  transition: all var(--transition-normal);
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--border-secondary);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 1px var(--primary-alpha-30);
}

:deep(.el-input__inner) {
  color: var(--text-primary);
  font-weight: var(--font-medium);
}

:deep(.el-input__inner::placeholder) {
  color: var(--text-tertiary);
}

:deep(.el-checkbox__label) {
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

:deep(.el-link.el-link--primary) {
  color: var(--primary-color);
  font-weight: var(--font-medium);
}

:deep(.el-link.el-link--primary:hover) {
  color: var(--primary-light);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-form {
    padding: var(--spacing-lg);
    margin: var(--spacing-md);
  }
  
  .logo-text {
    font-size: var(--font-2xl);
  }
  
  .login-button {
    height: 44px;
    font-size: var(--font-sm);
  }
}

@media (max-width: 480px) {
  .login-form {
    padding: var(--spacing-md);
  }
  
  .login-header {
    margin-bottom: var(--spacing-lg);
  }
}
</style>
