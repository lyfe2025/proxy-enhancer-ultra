<template>
  <el-dialog
    v-model="visible"
    :title="editingProxy ? '编辑代理' : '添加代理'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="formRules"
      label-width="100px"
    >
      <el-form-item label="代理名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入代理名称" />
      </el-form-item>
      <el-form-item label="代理类型" prop="type">
        <el-select v-model="form.type" placeholder="选择代理类型">
          <el-option label="HTTP" value="http" />
          <el-option label="HTTPS" value="https" />
          <el-option label="SOCKS5" value="socks5" />
          <el-option label="SOCKS4" value="socks4" />
        </el-select>
      </el-form-item>
      <el-form-item label="服务器地址" prop="host">
        <el-input v-model="form.host" placeholder="请输入服务器地址" />
      </el-form-item>
      <el-form-item label="端口" prop="port">
        <el-input-number
          v-model="form.port"
          :min="1"
          :max="65535"
          placeholder="请输入端口号"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="用户名">
        <el-input v-model="form.username" placeholder="请输入用户名（可选）" />
      </el-form-item>
      <el-form-item label="密码">
        <el-input
          v-model="form.password"
          type="password"
          placeholder="请输入密码（可选）"
          show-password
        />
      </el-form-item>
      <el-form-item label="描述">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入代理描述（可选）"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        {{ editingProxy ? '更新' : '添加' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormItemRule } from 'element-plus'
import type { ProxyConfig } from '@/api/proxy'

interface ProxyForm {
  name: string
  type: string
  host: string
  port: number
  username: string
  password: string
  description: string
}

const props = defineProps<{
  modelValue: boolean
  editingProxy?: ProxyConfig | null
  saving?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  save: [form: ProxyForm]
  close: []
}>()

const visible = ref(props.modelValue)
const formRef = ref()

const form = reactive<ProxyForm>({
  name: '',
  type: 'http',
  host: '',
  port: 8080,
  username: '',
  password: '',
  description: ''
})

// 表单验证规则
const formRules: Record<string, FormItemRule[]> = {
  name: [
    { required: true, message: '请输入代理名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择代理类型', trigger: 'change' }
  ],
  host: [
    { required: true, message: '请输入服务器地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: '端口号必须在1-65535之间', trigger: 'blur' }
  ]
}

// 监听modelValue变化
watch(
  () => props.modelValue,
  (newValue) => {
    visible.value = newValue
    if (newValue && props.editingProxy) {
      // 编辑模式，填充表单数据
      Object.assign(form, {
        name: props.editingProxy.name,
        type: props.editingProxy.type,
        host: props.editingProxy.host,
        port: props.editingProxy.port,
        username: props.editingProxy.username || '',
        password: props.editingProxy.password || '',
        description: props.editingProxy.description || ''
      })
    } else if (newValue) {
      // 新增模式，重置表单
      resetForm()
    }
  }
)

// 监听visible变化
watch(visible, (newValue) => {
  emit('update:modelValue', newValue)
})

const handleSave = async () => {
  try {
    await formRef.value.validate()
    emit('save', { ...form })
  } catch (error) {
    // 验证失败
  }
}

const handleClose = () => {
  visible.value = false
  emit('close')
}

const resetForm = () => {
  Object.assign(form, {
    name: '',
    type: 'http',
    host: '',
    port: 8080,
    username: '',
    password: '',
    description: ''
  })
  formRef.value?.resetFields()
}

// 暴露重置方法给父组件
defineExpose({
  resetForm
})
</script>

<style scoped>
/* 组件样式在这里 */
</style>
