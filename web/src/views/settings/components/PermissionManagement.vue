<template>
  <div class="permission-management">
    <!-- 权限筛选栏 -->
    <div class="filter-bar">
      <div class="filter-left">
        <el-input
          v-model="searchQuery"
          placeholder="搜索权限名称或代码"
          prefix-icon="Search"
          clearable
          style="width: 300px; margin-right: 16px"
        />
        <el-select
          v-model="moduleFilter"
          placeholder="筛选模块"
          clearable
          style="width: 150px"
        >
          <el-option label="全部模块" value="" />
          <el-option
            v-for="module in modules"
            :key="module"
            :label="module"
            :value="module"
          />
        </el-select>
      </div>
      <div class="filter-right">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新增权限
        </el-button>
        <el-button @click="refreshPermissions">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 权限列表 -->
    <el-table
      :data="filteredPermissions"
      v-loading="loading"
      stripe
      style="width: 100%"
    >
      <el-table-column prop="name" label="权限名称" width="180" />
      <el-table-column prop="code" label="权限代码" width="200">
        <template #default="{ row }">
          <el-tag type="info" size="small">{{ row.code }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="所属模块" width="120">
        <template #default="{ row }">
          <el-tag :type="getModuleTagType(row.module)">
            {{ row.module || '未分类' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" width="250" show-overflow-tooltip />
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="editPermission(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deletePermission(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 权限编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑权限' : '新增权限'"
      width="500px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入权限名称" />
        </el-form-item>
        <el-form-item label="权限代码" prop="code">
          <el-input
            v-model="formData.code"
            placeholder="请输入权限代码，如：user:create"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="所属模块" prop="module">
          <el-select
            v-model="formData.module"
            placeholder="请选择或输入模块名称"
            filterable
            allow-create
            style="width: 100%"
          >
            <el-option
              v-for="module in modules"
              :key="module"
              :label="module"
              :value="module"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入权限描述"
          />
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
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'
import type { Permission, PermissionForm } from '@/types/system'
// import { permissionApi } from '@/api/system' // TODO: 实现API

// Props
interface Props {
  permissions: Permission[]
  loading: boolean
}

const props = defineProps<Props>()

// Emits
interface Emits {
  'permission-created': [permission: Permission]
  'permission-updated': [permission: Permission]
  'permission-deleted': [permissionId: string]
  'refresh-permissions': []
}

const emit = defineEmits<Emits>()

// 响应式数据
const searchQuery = ref('')
const moduleFilter = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

// 表单数据
const formData = ref<PermissionForm>({
  name: '',
  code: '',
  type: 'menu',
  module: '',
  description: '',
  sort: 0,
  status: 'active'
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' },
    { min: 2, max: 50, message: '权限名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入权限代码', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z][a-zA-Z0-9_]*:[a-zA-Z][a-zA-Z0-9_]*$/,
      message: '权限代码格式应为：模块:操作，如 user:create',
      trigger: 'blur'
    }
  ],
  module: [
    { required: true, message: '请选择或输入模块名称', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入权限描述', trigger: 'blur' },
    { max: 200, message: '描述长度不能超过 200 个字符', trigger: 'blur' }
  ]
}

// 计算属性
const filteredPermissions = computed(() => {
  return props.permissions.filter(permission => {
    const matchesSearch = !searchQuery.value || 
      permission.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      permission.code.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesModule = !moduleFilter.value || permission.module === moduleFilter.value
    
    return matchesSearch && matchesModule
  })
})

const modules = computed(() => {
  const moduleSet = new Set<string>()
  props.permissions.forEach(permission => {
    if (permission.module) {
      moduleSet.add(permission.module)
    }
  })
  return Array.from(moduleSet).sort()
})

// 方法
const showCreateDialog = () => {
  isEdit.value = false
  dialogVisible.value = true
}

const editPermission = (permission: Permission) => {
  isEdit.value = true
  formData.value = {
    id: permission.id,
    name: permission.name,
    code: permission.code,
    type: permission.type,
    module: permission.module,
    description: permission.description || '',
    sort: permission.sort,
    status: permission.status
  }
  dialogVisible.value = true
}

const resetForm = () => {
  formData.value = {
    name: '',
    code: '',
    type: 'menu',
    module: '',
    description: '',
    sort: 0,
    status: 'active'
  }
  formRef.value?.resetFields()
}

const submitForm = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // TODO: 实现API调用
    if (isEdit.value) {
      // const updatedPermission = await permissionApi.updatePermission(formData.value.id, {
      //   name: formData.value.name,
      //   module: formData.value.module,
      //   description: formData.value.description
      // })
      // emit('permission-updated', updatedPermission)
      ElMessage.success('权限更新成功')
    } else {
      // const newPermission = await permissionApi.createPermission({
      //   name: formData.value.name,
      //   code: formData.value.code,
      //   module: formData.value.module,
      //   description: formData.value.description
      // })
      // emit('permission-created', newPermission)
      ElMessage.success('权限创建成功')
    }
    
    dialogVisible.value = false
  } catch (error) {
    console.error('提交权限表单失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    submitting.value = false
  }
}

const deletePermission = async (permission: Permission) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除权限 "${permission.name}" 吗？此操作不可恢复。`,
      '删除权限',
      { type: 'warning' }
    )
    
    // TODO: 实现API调用
    // await permissionApi.deletePermission(permission.id)
    emit('permission-deleted', permission.id)
    ElMessage.success('权限删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除权限失败:', error)
      ElMessage.error('删除权限失败，请重试')
    }
  }
}

const refreshPermissions = () => {
  emit('refresh-permissions')
}

const getModuleTagType = (module?: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  if (!module) return 'info'
  const moduleTypes: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    '用户管理': 'primary',
    '角色管理': 'success',
    '权限管理': 'warning',
    '系统管理': 'danger',
    '代理管理': 'info'
  }
  return moduleTypes[module] || 'primary'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.permission-management {
  padding: 20px;
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px;
  background: var(--el-bg-color-page);
  border-radius: 8px;
}

.filter-left {
  display: flex;
  align-items: center;
}

.filter-right {
  display: flex;
  gap: 8px;
}

:deep(.el-table .el-table__cell) {
  padding: 12px 0;
}

:deep(.el-tag) {
  border-radius: 4px;
}
</style>