<template>
  <el-dialog
    v-model="dialogVisible"
    title="忘记密码"
    width="400px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
    >
      <el-form-item label="邮箱" prop="email">
        <el-input
          v-model="formData.email"
          placeholder="请输入注册邮箱"
          prefix-icon="Message"
          clearable
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ loading ? '发送中...' : '发送重置邮件' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

interface ForgotPasswordData {
  email: string
}

const props = defineProps<{
  loading: boolean
}>()

const emit = defineEmits<{
  submit: [data: ForgotPasswordData]
  close: []
}>()

const dialogVisible = defineModel<boolean>({ required: true })
const formRef = ref<FormInstance>()

const formData = reactive<ForgotPasswordData>({
  email: ''
})

const formRules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
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

const handleClose = () => {
  formRef.value?.resetFields()
  formData.email = ''
  emit('close')
}
</script>

<style scoped>
.dialog-footer {
  text-align: right;
}
</style>
