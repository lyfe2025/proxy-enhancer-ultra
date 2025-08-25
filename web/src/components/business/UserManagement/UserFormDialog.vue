<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? '编辑用户' : '新增用户'"
    width="500px"
    @close="handleClose"
    @update:model-value="updateVisible"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
    >
      <el-form-item label="用户名" prop="username">
        <el-input v-model="formData.username" :disabled="isEdit" />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="formData.email" />
      </el-form-item>
      <el-form-item v-if="!isEdit" label="密码" prop="password">
        <el-input v-model="formData.password" type="password" show-password />
      </el-form-item>
      <el-form-item v-if="!isEdit" label="确认密码" prop="confirm_password">
        <el-input v-model="formData.confirm_password" type="password" show-password />
      </el-form-item>
      <el-form-item label="角色" prop="role_id">
        <el-select v-model="formData.role_id" placeholder="请选择角色" style="width: 100%">
          <el-option
            v-for="role in roles"
            :key="role.id"
            :label="role.name"
            :value="role.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-switch v-model="formData.is_active" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="submitting">
        {{ isEdit ? '更新' : '创建' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import type { User, Role } from '@/types/system'

// Props
interface Props {
  visible: boolean
  user?: User | null
  roles: Role[]
  isEdit?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isEdit: false
})

// Emits
interface Emits {
  'update:visible': [visible: boolean]
  'save': [user: User]
}

const emit = defineEmits<Emits>()

// 响应式数据
const submitting = ref(false)
const formRef = ref<FormInstance>()

// 表单数据类型
interface UserForm {
  id?: string
  username: string
  email: string
  password?: string
  confirm_password?: string
  role_id: string
  status: 'active' | 'inactive' | 'locked'
  is_active: boolean
  phone?: string
  department?: string
  position?: string
  avatar?: string
}

// 表单数据
const formData = ref<UserForm>({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
  role_id: '',
  status: 'active',
  is_active: true
})

// 表单验证规则
const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== formData.value.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  role_id: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 监听用户数据变化
watch(() => props.user, (newUser) => {
  if (newUser && props.isEdit) {
    formData.value = {
      id: newUser.id,
      username: newUser.username,
      email: newUser.email,
      password: '',
      confirm_password: '',
      role_id: newUser.role_id,
      status: newUser.status,
      is_active: newUser.is_active,
      phone: newUser.phone,
      department: newUser.department,
      position: newUser.position,
      avatar: newUser.avatar
    }
  } else {
    resetForm()
  }
}, { immediate: true })

// 方法
const updateVisible = (visible: boolean) => {
  emit('update:visible', visible)
}

const handleClose = () => {
  resetForm()
  emit('update:visible', false)
}

const resetForm = () => {
  formData.value = {
    username: '',
    email: '',
    password: '',
    confirm_password: '',
    role_id: '',
    status: 'active',
    is_active: true
  }
  formRef.value?.resetFields()
}

const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    const userData: User = {
      id: formData.value.id || '',
      username: formData.value.username,
      email: formData.value.email,
      nickname: formData.value.username,
      avatar: formData.value.avatar || '',
      status: formData.value.status,
      is_active: formData.value.is_active,
      role_id: formData.value.role_id,
      phone: formData.value.phone || '',
      department: formData.value.department || '',
      position: formData.value.position || '',
      created_at: props.user?.created_at || new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    
    emit('save', userData)
    ElMessage.success(props.isEdit ? '用户更新成功' : '用户创建成功')
  } catch (error) {
    console.error('保存用户失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
/* 组件特定样式 */
</style>
