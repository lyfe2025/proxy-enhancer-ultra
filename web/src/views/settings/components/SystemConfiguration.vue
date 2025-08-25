<template>
  <div class="system-configuration">
    <!-- 配置导航 -->
    <div class="config-nav">
      <el-tabs v-model="activeTab" type="card">
        <el-tab-pane label="基本设置" name="basic" />
        <el-tab-pane label="安全设置" name="security" />
        <el-tab-pane label="邮件设置" name="email" />
        <el-tab-pane label="短信设置" name="sms" />
        <el-tab-pane label="存储设置" name="storage" />
        <el-tab-pane label="日志设置" name="logging" />
        <el-tab-pane label="监控设置" name="monitoring" />
      </el-tabs>
    </div>

    <!-- 配置内容 -->
    <div class="config-content">
      <el-form
        ref="formRef"
        :model="configData"
        :rules="formRules"
        label-width="120px"
        v-loading="loading"
      >
        <!-- 基本设置 -->
        <div v-show="activeTab === 'basic'" class="config-section">
          <h3>基本设置</h3>
          <el-form-item label="系统名称" prop="system_name">
            <el-input v-model="configData.system_name" placeholder="请输入系统名称" />
          </el-form-item>
          <el-form-item label="系统描述" prop="system_description">
            <el-input
              v-model="configData.system_description"
              type="textarea"
              :rows="3"
              placeholder="请输入系统描述"
            />
          </el-form-item>
          <el-form-item label="系统版本" prop="system_version">
            <el-input v-model="configData.system_version" placeholder="请输入系统版本" />
          </el-form-item>
          <el-form-item label="维护模式">
            <el-switch
              v-model="configData.maintenance_mode"
              active-text="开启"
              inactive-text="关闭"
            />
          </el-form-item>
          <el-form-item label="维护提示" v-if="configData.maintenance_mode">
            <el-input
              v-model="configData.maintenance_message"
              type="textarea"
              :rows="2"
              placeholder="请输入维护提示信息"
            />
          </el-form-item>
        </div>

        <!-- 安全设置 -->
        <div v-show="activeTab === 'security'" class="config-section">
          <h3>安全设置</h3>
          <el-form-item label="密码最小长度" prop="password_min_length">
            <el-input-number
              v-model="configData.password_min_length"
              :min="6"
              :max="20"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="密码复杂度">
            <el-switch
              v-model="configData.password_complexity"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
          <el-form-item label="登录失败限制" prop="max_login_attempts">
            <el-input-number
              v-model="configData.max_login_attempts"
              :min="3"
              :max="10"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="会话超时(分钟)" prop="session_timeout">
            <el-input-number
              v-model="configData.session_timeout"
              :min="15"
              :max="1440"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="双因子认证">
            <el-switch
              v-model="configData.enable_2fa"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
          <el-form-item label="验证码">
            <el-switch
              v-model="configData.enable_captcha"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
        </div>

        <!-- 邮件设置 -->
        <div v-show="activeTab === 'email'" class="config-section">
          <h3>邮件设置</h3>
          <el-form-item label="SMTP服务器" prop="email_host">
            <el-input v-model="configData.email_host" placeholder="请输入SMTP服务器地址" />
          </el-form-item>
          <el-form-item label="SMTP端口" prop="email_port">
            <el-input-number
              v-model="configData.email_port"
              :min="1"
              :max="65535"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="用户名" prop="email_username">
            <el-input v-model="configData.email_username" placeholder="请输入SMTP用户名" />
          </el-form-item>
          <el-form-item label="密码" prop="email_password">
            <el-input
              v-model="configData.email_password"
              type="password"
              placeholder="请输入SMTP密码"
              show-password
            />
          </el-form-item>
          <el-form-item label="发件人邮箱" prop="email_from">
            <el-input v-model="configData.email_from" placeholder="请输入发件人邮箱" />
          </el-form-item>
          <el-form-item label="发件人名称" prop="email_from_name">
            <el-input v-model="configData.email_from_name" placeholder="请输入发件人名称" />
          </el-form-item>
          <el-form-item label="SSL加密">
            <el-switch
              v-model="configData.email_secure"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
        </div>

        <!-- 短信设置 -->
        <div v-show="activeTab === 'sms'" class="config-section">
          <h3>短信设置</h3>
          <el-form-item label="启用短信">
            <el-switch
              v-model="configData.sms_enabled"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
          <template v-if="configData.sms_enabled">
            <el-form-item label="短信服务商" prop="sms_provider">
              <el-select v-model="configData.sms_provider" placeholder="请选择短信服务商">
                <el-option label="阿里云" value="aliyun" />
                <el-option label="腾讯云" value="tencent" />
                <el-option label="华为云" value="huawei" />
              </el-select>
            </el-form-item>
            <el-form-item label="Access Key" prop="sms_access_key">
              <el-input v-model="configData.sms_access_key" placeholder="请输入Access Key" />
            </el-form-item>
            <el-form-item label="Secret Key" prop="sms_secret_key">
              <el-input
                v-model="configData.sms_secret_key"
                type="password"
                placeholder="请输入Secret Key"
                show-password
              />
            </el-form-item>
            <el-form-item label="签名" prop="sms_signature">
              <el-input v-model="configData.sms_signature" placeholder="请输入短信签名" />
            </el-form-item>
          </template>
        </div>

        <!-- 存储设置 -->
        <div v-show="activeTab === 'storage'" class="config-section">
          <h3>存储设置</h3>
          <el-form-item label="存储类型" prop="storage_type">
            <el-select v-model="configData.storage_type" placeholder="请选择存储类型">
              <el-option label="本地存储" value="local" />
              <el-option label="阿里云OSS" value="aliyun" />
              <el-option label="腾讯云COS" value="tencent" />
              <el-option label="七牛云" value="qiniu" />
            </el-select>
          </el-form-item>
          <template v-if="configData.storage_type !== 'local'">
            <el-form-item label="Bucket名称" prop="storage_bucket">
              <el-input v-model="configData.storage_bucket" placeholder="请输入Bucket名称" />
            </el-form-item>
            <el-form-item label="访问域名" prop="storage_domain">
              <el-input v-model="configData.storage_domain" placeholder="请输入访问域名" />
            </el-form-item>
            <el-form-item label="Access Key" prop="storage_access_key">
              <el-input v-model="configData.storage_access_key" placeholder="请输入Access Key" />
            </el-form-item>
            <el-form-item label="Secret Key" prop="storage_secret_key">
              <el-input
                v-model="configData.storage_secret_key"
                type="password"
                placeholder="请输入Secret Key"
                show-password
              />
            </el-form-item>
          </template>
        </div>

        <!-- 日志设置 -->
        <div v-show="activeTab === 'logging'" class="config-section">
          <h3>日志设置</h3>
          <el-form-item label="日志级别" prop="logging.level">
            <el-select v-model="configData.log_level" placeholder="请选择日志级别">
              <el-option label="DEBUG" value="debug" />
              <el-option label="INFO" value="info" />
              <el-option label="WARN" value="warn" />
              <el-option label="ERROR" value="error" />
            </el-select>
          </el-form-item>
          <el-form-item label="日志保留天数" prop="logging.retention_days">
            <el-input-number
              v-model="configData.log_retention_days"
              :min="1"
              :max="365"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="启用文件日志">
            <el-switch
              v-model="configData.enable_access_log"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
          <el-form-item label="启用数据库日志">
            <el-switch
              v-model="configData.enable_audit_log"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
        </div>

        <!-- 监控设置 -->
        <div v-show="activeTab === 'monitoring'" class="config-section">
          <h3>监控设置</h3>
          <el-form-item label="启用监控">
            <el-switch
              v-model="configData.monitoring!.enabled"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
          <template v-if="configData.monitoring?.enabled">
            <el-form-item label="数据收集间隔(秒)" prop="monitoring.collection_interval">
              <el-input-number
                v-model="configData.monitoring!.collection_interval"
                :min="10"
                :max="3600"
                controls-position="right"
              />
            </el-form-item>
            <el-form-item label="数据保留天数" prop="monitoring.retention_days">
              <el-input-number
                v-model="configData.monitoring!.retention_days"
                :min="1"
                :max="365"
                controls-position="right"
              />
            </el-form-item>
            <el-form-item label="告警阈值">
              <div class="threshold-group">
                <div class="threshold-item">
                  <label>CPU使用率(%):</label>
                  <el-input-number
                    v-model="configData.monitoring!.cpu_threshold"
                    :min="50"
                    :max="95"
                    controls-position="right"
                  />
                </div>
                <div class="threshold-item">
                  <label>内存使用率(%):</label>
                  <el-input-number
                    v-model="configData.monitoring!.memory_threshold"
                    :min="50"
                    :max="95"
                    controls-position="right"
                  />
                </div>
                <div class="threshold-item">
                  <label>磁盘使用率(%):</label>
                  <el-input-number
                    v-model="configData.monitoring!.disk_threshold"
                    :min="70"
                    :max="95"
                    controls-position="right"
                  />
                </div>
              </div>
            </el-form-item>
          </template>
        </div>
      </el-form>
    </div>

    <!-- 操作按钮 -->
    <div class="config-actions">
      <el-button @click="resetConfig">重置</el-button>
      <el-button type="primary" @click="saveConfig" :loading="saving">
        保存配置
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import type { SystemConfig } from '@/types/system'
// import { systemApi } from '@/api/system'

// Props
interface Props {
  config: SystemConfig
  loading: boolean
}

const props = defineProps<Props>()

// Emits
interface Emits {
  'config-updated': [config: SystemConfig]
  'config-reset': []
}

const emit = defineEmits<Emits>()

// 响应式数据
const activeTab = ref('basic')
const saving = ref(false)
const formRef = ref<FormInstance>()

// 配置数据
const configData = ref<SystemConfig>({
  system_name: '',
  system_description: '',
  system_version: '',
  maintenance_mode: false,
  maintenance_message: '',
  password_min_length: 8,
  password_complexity: false,
  max_login_attempts: 5,
  session_timeout: 120,
  enable_2fa: false,
  enable_captcha: false,
  email_host: '',
  email_port: 587,
  email_username: '',
  email_password: '',
  email_from: '',
  email_from_name: '',
  email_secure: false,
  sms_enabled: false,
  sms_provider: '',
  sms_access_key: '',
  sms_secret_key: '',
  sms_signature: '',
  storage_type: 'local',
  storage_bucket: '',
  storage_domain: '',
  storage_access_key: '',
  storage_secret_key: '',
  log_level: 'info',
  log_retention_days: 30,
  enable_access_log: true,
  enable_error_log: true,
  enable_audit_log: true,
  monitoring: {
    enabled: false,
    collection_interval: 60,
    retention_days: 30,
    cpu_threshold: 80,
    memory_threshold: 80,
    disk_threshold: 85
  }
})

// 表单验证规则
const formRules = {
  system_name: [
    { required: true, message: '请输入系统名称', trigger: 'blur' }
  ],
  system_version: [
    { required: true, message: '请输入系统版本', trigger: 'blur' }
  ],
  password_min_length: [
    { required: true, message: '请设置密码最小长度', trigger: 'blur' }
  ],
  email_host: [
    { required: true, message: '请输入SMTP服务器地址', trigger: 'blur' }
  ],
  email_from: [
    { required: true, message: '请输入发件人邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

// 监听props变化
watch(
  () => props.config,
  (newConfig) => {
    if (newConfig) {
      configData.value = { ...newConfig }
    }
  },
  { immediate: true, deep: true }
)

// 方法
const saveConfig = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    saving.value = true
    
    // TODO: 实现API调用
    // const updatedConfig = await systemApi.updateConfig(configData.value)
    emit('config-updated', configData.value)
    ElMessage.success('配置保存成功')
  } catch (error) {
    console.error('保存配置失败:', error)
    ElMessage.error('保存配置失败，请重试')
  } finally {
    saving.value = false
  }
}

const resetConfig = () => {
  configData.value = { ...props.config }
  formRef.value?.resetFields()
  emit('config-reset')
  ElMessage.info('配置已重置')
}
</script>

<style scoped>
.system-configuration {
  padding: 20px;
}

.config-nav {
  margin-bottom: 20px;
}

.config-content {
  background: var(--el-bg-color-page);
  padding: 24px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.config-section h3 {
  margin: 0 0 20px 0;
  color: var(--el-text-color-primary);
  font-size: 16px;
  font-weight: 600;
}

.threshold-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.threshold-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.threshold-item label {
  min-width: 120px;
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.config-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px;
  background: var(--el-bg-color-page);
  border-radius: 8px;
}

:deep(.el-tabs__item) {
  padding: 0 20px;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>