<template>
  <div class="role-management">
    <!-- 角色操作栏 -->
    <div class="action-bar">
      <div class="action-left">
        <el-input
          v-model="searchQuery"
          placeholder="搜索角色名称或描述"
          prefix-icon="Search"
          clearable
          style="width: 300px; margin-right: 16px"
        />
      </div>
      <div class="action-right">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新增角色
        </el-button>
        <el-button @click="refreshRoles">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 角色列表 -->
    <el-table
      :data="filteredRoles"
      v-loading="loading"
      stripe
      style="width: 100%"
    >
      <el-table-column prop="name" label="角色名称" width="150" />
      <el-table-column prop="description" label="描述" width="200" />
      <el-table-column label="用户数量" width="100">
        <template #default="{ row }">
          <el-tag type="info">{{ row.user_count || 0 }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="权限" width="300">
        <template #default="{ row }">
          <div class="permission-tags">
            <el-tag
              v-for="permission in row.permissions?.slice(0, 3)"
              :key="permission.id"
              size="small"
              style="margin-right: 4px; margin-bottom: 4px"
            >
              {{ permission.name }}
            </el-tag>
            <el-tag
              v-if="row.permissions && row.permissions.length > 3"
              size="small"
              type="info"
            >
              +{{ row.permissions.length - 3 }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="editRole(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteRole(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 角色编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑角色' : '新增角色'"
      width="600px"
      @close="resetForm"
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
                @change="handleCheckAllChange"
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
                    @change="(checked: CheckboxValueType) => handleModuleCheckChange(module, checked as boolean)"
                  >
                    {{ module }}
                  </el-checkbox>
                </div>
                <div class="group-permissions">
                  <el-checkbox
                    v-for="permission in group"
                    :key="permission.id"
                    :value="selectedPermissions.includes(permission.id)"
                    :label="permission.id"
                    @change="(checked: CheckboxValueType) => handlePermissionChange(permission.id, checked as boolean)"
                  >
                    {{ permission.name }}
                  </el-checkbox>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type CheckboxValueType } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'
import type { Role, Permission, RoleForm } from '@/types/system'
// import { roleApi } from '@/api/system' // TODO: 实现API

// Props
interface Props {
  roles: Role[]
  permissions: Permission[]
  loading: boolean
}

const props = defineProps<Props>()

// Emits
interface Emits {
  'role-created': [role: Role]
  'role-updated': [role: Role]
  'role-deleted': [roleId: string]
  'refresh-roles': []
}

const emit = defineEmits<Emits>()

// 响应式数据
const searchQuery = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const selectedPermissions = ref<string[]>([])
const checkAll = ref(false)
const moduleCheckStatus = ref<Record<string, boolean>>({})

// 表单数据
const formData = ref<RoleForm>({
  name: '',
  description: '',
  permission_ids: [],
  status: 'active'
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
const filteredRoles = computed(() => {
  return props.roles.filter(role => {
    const matchesSearch = !searchQuery.value || 
      role.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      role.description?.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    return matchesSearch
  })
})

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

// 监听选中权限变化
watch(selectedPermissions, (newVal) => {
  checkAll.value = newVal.length === props.permissions.length
  updateModuleCheckStatus()
  formData.value.permission_ids = newVal
}, { deep: true })

// 方法
const showCreateDialog = () => {
  isEdit.value = false
  dialogVisible.value = true
}

const editRole = (role: Role) => {
  isEdit.value = true
  formData.value = {
    id: role.id,
    name: role.name,
    description: role.description || '',
    permission_ids: role.permissions?.map(p => p.id) || [],
    status: role.status
  }
  selectedPermissions.value = role.permissions?.map(p => p.id) || []
  dialogVisible.value = true
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

const submitForm = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // TODO: 实现API调用
    if (isEdit.value) {
      // const updatedRole = await roleApi.updateRole(formData.value.id, {
      //   name: formData.value.name,
      //   description: formData.value.description,
      //   permission_ids: selectedPermissions.value
      // })
      // emit('role-updated', updatedRole)
      ElMessage.success('角色更新成功')
    } else {
      // const newRole = await roleApi.createRole({
      //   name: formData.value.name,
      //   description: formData.value.description,
      //   permission_ids: selectedPermissions.value
      // })
      // emit('role-created', newRole)
      ElMessage.success('角色创建成功')
    }
    
    dialogVisible.value = false
  } catch (error) {
    console.error('提交角色表单失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    submitting.value = false
  }
}

const deleteRole = async (role: Role) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${role.name}" 吗？此操作不可恢复。`,
      '删除角色',
      { type: 'warning' }
    )
    
    // TODO: 实现API调用
    // await roleApi.deleteRole(role.id)
    emit('role-deleted', role.id)
    ElMessage.success('角色删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除角色失败:', error)
      ElMessage.error('删除角色失败，请重试')
    }
  }
}

const refreshRoles = () => {
  emit('refresh-roles')
}

const handleCheckAllChange = (checked: CheckboxValueType) => {
  if (checked as boolean) {
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
  updatePermissionSelection()
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
  
  updatePermissionSelection()
}

const updatePermissionSelection = () => {
  updateModuleCheckStatus()
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

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.role-management {
  padding: 20px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px;
  background: var(--el-bg-color-page);
  border-radius: 8px;
}

.action-left {
  display: flex;
  align-items: center;
}

.action-right {
  display: flex;
  gap: 8px;
}

.permission-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

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