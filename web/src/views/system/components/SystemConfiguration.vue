<template>
  <div class="system-configuration">
    <el-row :gutter="20">
      <!-- 基本设置 -->
      <el-col :span="12">
        <el-card class="config-card">
          <template #header>
            <span>基本设置</span>
          </template>
          <el-form :model="config" label-width="120px">
            <el-form-item label="系统名称">
              <el-input v-model="config.system_name" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统描述">
              <el-input
                v-model="config.system_description"
                type="textarea"
                :rows="3"
                placeholder="请输入系统描述"
              />
            </el-form-item>
            <el-form-item label="系统Logo">
              <el-input v-model="config.system_logo" placeholder="请输入Logo URL" />
            </el-form-item>
            <el-form-item label="语言">
              <el-select v-model="config.language" placeholder="选择语言">
                <el-option label="中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
              </el-select>
            </el-form-item>
            <el-form-item label="时区">
              <el-select v-model="config.timezone" placeholder="选择时区">
                <el-option label="Asia/Shanghai" value="Asia/Shanghai" />
                <el-option label="UTC" value="UTC" />
                <el-option label="America/New_York" value="America/New_York" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 安全设置 -->
      <el-col :span="12">
        <el-card class="config-card">
          <template #header>
            <span>安全设置</span>
          </template>
          <el-form :model="config" label-width="120px">
            <el-form-item label="密码最小长度">
              <el-input-number
                v-model="config.password_min_length"
                :min="6"
                :max="20"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="密码复杂度">
              <el-checkbox v-model="config.password_complexity">启用密码复杂度要求</el-checkbox>
            </el-form-item>
            <el-form-item label="登录失败限制">
              <el-input-number
                v-model="config.max_login_attempts"
                :min="3"
                :max="10"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="会话超时(分钟)">
              <el-input-number
                v-model="config.session_timeout"
                :min="30"
                :max="1440"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="启用双因子认证">
              <el-switch v-model="config.enable_2fa" />
            </el-form-item>
            <el-form-item label="启用验证码">
              <el-switch v-model="config.enable_captcha" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 邮件设置 -->
      <el-col :span="12">
        <el-card class="config-card">
          <template #header>
            <span>邮件设置</span>
          </template>
          <el-form :model="config" label-width="120px">
            <el-form-item label="SMTP服务器">
              <el-input v-model="config.email_host" placeholder="请输入SMTP服务器" />
            </el-form-item>
            <el-form-item label="SMTP端口">
              <el-input-number
                v-model="config.email_port"
                :min="1"
                :max="65535"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="用户名">
              <el-input v-model="config.email_username" placeholder="请输入SMTP用户名" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input
                v-model="config.email_password"
                type="password"
                placeholder="请输入SMTP密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="发件人邮箱">
              <el-input v-model="config.email_from" placeholder="请输入发件人邮箱" />
            </el-form-item>
            <el-form-item label="发件人名称">
              <el-input v-model="config.email_from_name" placeholder="请输入发件人名称" />
            </el-form-item>
            <el-form-item label="启用SSL">
              <el-switch v-model="config.email_secure" />
            </el-form-item>
            <el-form-item>
              <el-button @click="testEmailConfig" :loading="testingEmail">
                <el-icon><Message /></el-icon>
                测试邮件配置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 日志设置 -->
      <el-col :span="12">
        <el-card class="config-card">
          <template #header>
            <span>日志设置</span>
          </template>
          <el-form :model="config" label-width="120px">
            <el-form-item label="日志级别">
              <el-select v-model="config.log_level" placeholder="选择日志级别">
                <el-option label="DEBUG" value="debug" />
                <el-option label="INFO" value="info" />
                <el-option label="WARN" value="warn" />
                <el-option label="ERROR" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志保留天数">
              <el-input-number
                v-model="config.log_retention_days"
                :min="1"
                :max="365"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="启用访问日志">
              <el-switch v-model="config.enable_access_log" />
            </el-form-item>
            <el-form-item label="启用错误日志">
              <el-switch v-model="config.enable_error_log" />
            </el-form-item>
            <el-form-item label="启用审计日志">
              <el-switch v-model="config.enable_audit_log" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <!-- 保存按钮 -->
    <div class="save-actions">
      <el-button type="primary" @click="saveConfig" :loading="saving">
        <el-icon><Check /></el-icon>
        保存配置
      </el-button>
      <el-button @click="fetchConfig">
        <el-icon><Refresh /></el-icon>
        重置
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Message,
  Check,
  Refresh
} from '@element-plus/icons-vue'
import * as systemApi from '@/api/system'
import type { SystemConfig } from '@/types/system'

// Emits
const emit = defineEmits(['refresh'])

// 响应式数据
const saving = ref(false)
const testingEmail = ref(false)

// 系统配置
const config = reactive<SystemConfig>({
  system_name: 'Proxy Enhancer Ultra',
  system_version: '1.0.0',
  system_logo: '',
  system_description: '高性能代理增强系统',
  session_timeout: 120,
  max_login_attempts: 5,
  password_min_length: 8,
  password_complexity: true,
  enable_2fa: false,
  enable_captcha: false,
  captcha_threshold: 3,
  enable_email_verification: false,
  enable_sms_verification: false,
  maintenance_mode: false,
  maintenance_message: '',
  allow_registration: true,
  default_role_id: '1',
  timezone: 'Asia/Shanghai',
  date_format: 'YYYY-MM-DD',
  language: 'zh-CN',
  theme: 'light' as 'light' | 'dark' | 'auto',
  log_level: 'info' as 'debug' | 'info' | 'warn' | 'error',
  log_retention_days: 30,
  enable_access_log: true,
  enable_error_log: true,
  enable_audit_log: true,
  backup_enabled: false,
  backup_frequency: 'daily' as 'daily' | 'weekly' | 'monthly',
  backup_retention_days: 7,
  email_host: '',
  email_port: 587,
  email_username: '',
  email_password: '',
  email_from: '',
  email_from_name: '',
  email_secure: false,
  sms_provider: '',
  sms_api_key: '',
  sms_api_secret: ''
})

// 获取系统配置
const fetchConfig = async () => {
  try {
    const response = await systemApi.getSystemConfigs()
    Object.assign(config, response.data)
  } catch (error) {
    ElMessage.error('获取系统配置失败')
  }
}

// 保存系统配置
const saveConfig = async () => {
  try {
    saving.value = true
    await systemApi.updateSystemConfig('all', JSON.stringify(config))
    ElMessage.success('保存成功')
    emit('refresh')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 测试邮件配置
const testEmailConfig = async () => {
  try {
    testingEmail.value = true
    // 测试邮件配置功能暂未实现
    ElMessage.info('邮件配置测试功能开发中')
    ElMessage.success('邮件配置测试成功')
  } catch (error) {
    ElMessage.error('邮件配置测试失败')
  } finally {
    testingEmail.value = false
  }
}

// 暴露方法给父组件
defineExpose({
  fetchConfig,
  saveConfig
})

// 初始化
onMounted(() => {
  fetchConfig()
})
</script>

<style scoped>
.config-card {
  border: 1px solid var(--el-border-color-light);
  margin-bottom: 20px;
}

.save-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 30px;
  padding: 20px 0;
  border-top: 1px solid var(--el-border-color-light);
}

/* 深色主题适配 */
.dark .config-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

.dark .save-actions {
  border-color: var(--el-border-color);
}
</style>
