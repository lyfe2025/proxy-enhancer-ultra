<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? '编辑角色' : '新增角色'"
    width="600px"
    @close="handleClose"
    @update:model-value="updateVisible"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
    >
      <el-form-item label="角色名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入角色名称" />
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入角色描述"
        />
      </el-form-item>
      <el-form-item label="权限配置">
        <div class="permission-selector">
          <div class="permission-header">
            <el-checkbox
              :indeterminate="isIndeterminate"
              v-model="checkAll"
              @change="(val: any) => handleCheckAllChange(val as boolean)"
            >
              全选
            </el-checkbox>
            <span class="permission-count">
              已选择 {{ selectedPermissions.length }} / {{ permissions.length }} 个权限
            </span>
          </div>
          <div class="permission-groups">
            <div
              v-for="(group, module) in groupedPermissions"
              :key="module"
              class="permission-group"
            >
              <div class="group-header">
                <el-checkbox
                  :indeterminate="isModuleIndeterminate(module)"
                  v-model="moduleCheckStatus[module]"
                  @change="(checked: any) => handleModuleCheckChange(module, checked as boolean)"
                >
                  {{ module }}
                </el-checkbox>
              </div>
              <div class="group-permissions">
                <el-checkbox
                  v-for="permission in group"
                  :key="permission.id"
                  :model-value="selectedPermissions.includes(permission.id)"
                  :label="permission.name"
                  @change="(checked: any) => handlePermissionChange(permission.id, checked as boolean)"
                />
              </div>
            </div>
          </div>
        </div>
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
import { ref, computed, watch } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import type { Role, Permission } from '@/types/system'

// Props
interface Props {
  visible: boolean
  role?: Role | null
  permissions: Permission[]
  isEdit?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isEdit: false
})

// Emits
interface Emits {
  'update:visible': [visible: boolean]
  'save': [role: Role]
}

const emit = defineEmits<Emits>()

// 响应式数据
const submitting = ref(false)
const formRef = ref<FormInstance>()
const selectedPermissions = ref<string[]>([])
const checkAll = ref(false)
const moduleCheckStatus = ref<Record<string, boolean>>({})

// 表单数据
const formData = ref({
  name: '',
  description: '',
  permission_ids: [] as string[],
  status: 'active' as 'active' | 'inactive'
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '角色名称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入角色描述', trigger: 'blur' },
    { max: 200, message: '描述长度不能超过 200 个字符', trigger: 'blur' }
  ]
}

// 计算属性
const groupedPermissions = computed(() => {
  const groups: Record<string, Permission[]> = {}
  props.permissions.forEach(permission => {
    const module = permission.module || '其他'
    if (!groups[module]) {
      groups[module] = []
    }
    groups[module].push(permission)
  })
  return groups
})

const isIndeterminate = computed(() => {
  const selectedCount = selectedPermissions.value.length
  return selectedCount > 0 && selectedCount < props.permissions.length
})

// 监听角色数据变化
watch(() => props.role, (newRole) => {
  if (newRole && props.isEdit) {
    formData.value = {
      name: newRole.name,
      description: newRole.description || '',
      permission_ids: newRole.permissions?.map(p => p.id) || [],
      status: newRole.status
    }
    selectedPermissions.value = newRole.permissions?.map(p => p.id) || []
  } else {
    resetForm()
  }
}, { immediate: true })

// 监听选中权限变化
watch(selectedPermissions, (newVal) => {
  checkAll.value = newVal.length === props.permissions.length
  updateModuleCheckStatus()
  formData.value.permission_ids = newVal
}, { deep: true })

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
    name: '',
    description: '',
    permission_ids: [],
    status: 'active'
  }
  selectedPermissions.value = []
  checkAll.value = false
  moduleCheckStatus.value = {}
  formRef.value?.resetFields()
}

const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    const roleData: Role = {
      id: props.role?.id || '',
      name: formData.value.name,
      code: formData.value.name.toLowerCase().replace(/\s+/g, '_'),
      description: formData.value.description,
      status: formData.value.status,
      permissions: selectedPermissions.value.map(id => {
        const permission = props.permissions.find(p => p.id === id)
        return permission!
      }).filter(Boolean),
      created_at: props.role?.created_at || new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    
    emit('save', roleData)
    ElMessage.success(props.isEdit ? '角色更新成功' : '角色创建成功')
  } catch (error) {
    console.error('保存角色失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    submitting.value = false
  }
}

const handleCheckAllChange = (checked: boolean | string | number) => {
  if (checked) {
    selectedPermissions.value = props.permissions.map(p => p.id)
  } else {
    selectedPermissions.value = []
  }
}

const handlePermissionChange = (permissionId: string, checked: boolean) => {
  if (checked) {
    if (!selectedPermissions.value.includes(permissionId)) {
      selectedPermissions.value.push(permissionId)
    }
  } else {
    const index = selectedPermissions.value.indexOf(permissionId)
    if (index > -1) {
      selectedPermissions.value.splice(index, 1)
    }
  }
}

const handleModuleCheckChange = (module: string, checked: boolean) => {
  const modulePermissions = groupedPermissions.value[module] || []
  
  modulePermissions.forEach(permission => {
    if (checked) {
      if (!selectedPermissions.value.includes(permission.id)) {
        selectedPermissions.value.push(permission.id)
      }
    } else {
      const index = selectedPermissions.value.indexOf(permission.id)
      if (index > -1) {
        selectedPermissions.value.splice(index, 1)
      }
    }
  })
}

const updateModuleCheckStatus = () => {
  Object.keys(groupedPermissions.value).forEach(module => {
    const modulePermissions = groupedPermissions.value[module] || []
    const modulePermissionIds = modulePermissions.map(p => p.id)
    const selectedCount = modulePermissionIds.filter(id => 
      selectedPermissions.value.includes(id)
    ).length
    
    moduleCheckStatus.value[module] = selectedCount === modulePermissionIds.length
  })
}

const isModuleIndeterminate = (module: string) => {
  const modulePermissions = groupedPermissions.value[module] || []
  const modulePermissionIds = modulePermissions.map(p => p.id)
  const selectedCount = modulePermissionIds.filter(id => 
    selectedPermissions.value.includes(id)
  ).length
  
  return selectedCount > 0 && selectedCount < modulePermissionIds.length
}
</script>

<style scoped>
.permission-selector {
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  padding: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.permission-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.permission-count {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.permission-groups {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.permission-group {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  padding: 12px;
}

.group-header {
  margin-bottom: 8px;
  font-weight: 500;
}

.group-permissions {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 8px;
  margin-left: 20px;
}
</style>
