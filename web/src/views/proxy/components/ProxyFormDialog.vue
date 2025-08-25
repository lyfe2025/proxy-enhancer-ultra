<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑代理' : '添加代理'"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
      label-position="left"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="代理名称" prop="name">
            <el-input
              v-model="formData.name"
              placeholder="请输入代理名称"
              clearable
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="代理类型" prop="type">
            <el-select
              v-model="formData.type"
              placeholder="请选择代理类型"
              style="width: 100%"
            >
              <el-option
                v-for="option in PROXY_TYPE_OPTIONS"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="16">
          <el-form-item label="主机地址" prop="host">
            <el-input
              v-model="formData.host"
              placeholder="请输入主机地址或域名"
              clearable
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="端口" prop="port">
            <el-input-number
              v-model="formData.port"
              :min="PORT_RANGE.MIN"
              :max="PORT_RANGE.MAX"
              placeholder="端口"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" v-if="needsAuth">
        <el-col :span="12">
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="formData.username"
              placeholder="请输入用户名"
              clearable
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="formData.password"
              type="password"
              placeholder="请输入密码"
              show-password
              clearable
            />
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入代理描述（可选）"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>
      
      <el-form-item label="高级设置">
        <el-collapse v-model="activeCollapse">
          <el-collapse-item title="连接设置" name="connection">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="超时时间" prop="timeout">
                  <el-input-number
                    v-model="formData.timeout"
                    :min="1000"
                    :max="60000"
                    :step="1000"
                    placeholder="毫秒"
                    style="width: 100%"
                  />
                  <div class="form-tip">连接超时时间（毫秒）</div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="重试次数" prop="retry_count">
                  <el-input-number
                    v-model="formData.retry_count"
                    :min="0"
                    :max="10"
                    placeholder="次数"
                    style="width: 100%"
                  />
                  <div class="form-tip">连接失败重试次数</div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-collapse-item>
          
          <el-collapse-item title="其他设置" name="other">
            <el-form-item label="启用状态">
              <el-switch
                v-model="formData.enabled"
                active-text="启用"
                inactive-text="禁用"
              />
            </el-form-item>
            
            <el-form-item label="标签">
              <el-input
                v-model="formData.tags"
                placeholder="请输入标签，多个标签用逗号分隔"
                clearable
              />
              <div class="form-tip">用于分类和筛选代理</div>
            </el-form-item>
          </el-collapse-item>
        </el-collapse>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="submitting"
        >
          {{ isEdit ? '更新' : '添加' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { ProxyConfig } from '@/api/proxy'
import { PROXY_TYPE_OPTIONS, PORT_RANGE, DEFAULT_PROXY_FORM, PROXY_FORM_RULES } from '../constants/proxyConstants'

// 定义 props
const props = defineProps<{
  visible: boolean
  editData?: ProxyConfig | null
  submitting?: boolean
}>()

// 定义 emits
const emit = defineEmits<{
  'update:visible': [visible: boolean]
  submit: [data: Partial<ProxyConfig> & {
    timeout?: number
    retry_count?: number
    enabled?: boolean
    tags?: string | string[]
  }]
  close: []
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 对话框显示状态
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 是否为编辑模式
const isEdit = computed(() => !!props.editData)

// 表单数据
const formData = reactive<Partial<ProxyConfig> & {
  timeout?: number
  retry_count?: number
  enabled?: boolean
  tags?: string
}>({ ...DEFAULT_PROXY_FORM })

// 表单验证规则
const formRules: FormRules = PROXY_FORM_RULES

// 折叠面板激活项
const activeCollapse = ref<string[]>([])

// 是否需要认证
const needsAuth = computed(() => {
  return ['http', 'https', 'socks5'].includes(formData.type!)
})

// 监听编辑数据变化
watch(
  () => props.editData,
  (newData: ProxyConfig | null | undefined) => {
    if (newData) {
      Object.assign(formData, {
        ...newData,
        tags: Array.isArray(newData.tags) ? newData.tags.join(', ') : newData.tags || ''
      })
    } else {
      resetForm()
    }
  },
  { immediate: true, deep: true }
)

// 监听对话框显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      nextTick(() => {
        formRef.value?.clearValidate()
      })
    }
  }
)

// 重置表单
const resetForm = () => {
  Object.assign(formData, { ...DEFAULT_PROXY_FORM })
  activeCollapse.value = []
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

// 处理关闭
const handleClose = () => {
  emit('close')
  emit('update:visible', false)
}

// 处理提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    // 处理标签数据
    const submitData = {
      ...formData,
      tags: formData.tags
        ? formData.tags.split(',').map(tag => tag.trim()).filter(tag => tag)
        : []
    }
    
    emit('submit', submitData)
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin-top: 4px;
  line-height: 1.4;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

:deep(.el-collapse-item__header) {
  font-size: 14px;
  font-weight: 500;
}

:deep(.el-collapse-item__content) {
  padding-bottom: 0;
}

:deep(.el-input-number) {
  width: 100%;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .el-dialog {
    width: 95% !important;
    margin: 5vh auto;
  }
  
  :deep(.el-form-item__label) {
    width: 80px !important;
  }
  
  .el-row .el-col {
    margin-bottom: 12px;
  }
}

/* 深色主题适配 */
.dark .form-tip {
  color: var(--el-text-color-secondary);
}
</style>