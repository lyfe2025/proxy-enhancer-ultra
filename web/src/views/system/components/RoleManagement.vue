<template>
  <div class="role-management">
    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span>角色列表</span>
          <el-button type="primary" size="small" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            添加角色
          </el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="rolesList" stripe style="width: 100%">
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="permissions" label="权限" min-width="300">
          <template #default="{ row }">
            <el-tag
              v-for="permission in row.permissions.slice(0, 3)"
              :key="permission.id"
              size="small"
              style="margin-right: 4px; margin-bottom: 4px"
            >
              {{ permission.name }}
            </el-tag>
            <el-tag v-if="row.permissions.length > 3" size="small" type="info">
              +{{ row.permissions.length - 3 }} 更多
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="userCount" label="用户数" width="100" />
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            <span>{{ formatDate(row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="editRole(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button size="small" type="danger" @click="deleteRole(row)" :disabled="row.userCount > 0">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑角色对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingRole ? '编辑角色' : '创建角色'"
      width="800px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
          />
        </el-form-item>
        <el-form-item label="权限" prop="permissionIds">
          <el-tree
            ref="permissionTreeRef"
            :data="permissionTree"
            :props="{ children: 'children', label: 'name' }"
            node-key="id"
            show-checkbox
            :default-checked-keys="form.permissionIds"
            @check="handlePermissionCheck"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="saveRole" :loading="saving">
          {{ editingRole ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Delete,
  Edit
} from '@element-plus/icons-vue'
import * as systemApi from '@/api/system'
import type { Role, Permission } from '@/types/system'

// Emits
const emit = defineEmits(['refresh'])

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const showCreateDialog = ref(false)
const editingRole = ref<Role | null>(null)
const formRef = ref()
const permissionTreeRef = ref()

// 列表数据
const rolesList = ref<Role[]>([])
const permissionTree = ref<Permission[]>([])

// 角色表单
const form = reactive({
  name: '',
  description: '',
  permissionIds: [] as string[]
})

// 表单验证规则
const formRules: Record<string, any> = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' }
  ],
  permissionIds: [
    { required: true, message: '请选择权限', trigger: 'change' }
  ]
}

// 获取角色列表
const fetchRoles = async () => {
  try {
    loading.value = true
    const response = await systemApi.getRoleList({})
    rolesList.value = (response.data as any).data || []
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 获取权限树
const fetchPermissionTree = async () => {
  try {
    const response = await systemApi.getPermissionList()
    permissionTree.value = (response.data as any) || []
  } catch (error) {
    ElMessage.error('获取权限列表失败')
  }
}

// 角色操作
const editRole = (role: Role) => {
  editingRole.value = role
  Object.assign(form, {
    name: role.name,
    description: role.description || '',
    permissionIds: role.permissions?.map(p => p.id) || []
  })
  showCreateDialog.value = true
}

const deleteRole = async (role: Role) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${role.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await systemApi.deleteRole(role.id)
    ElMessage.success('删除成功')
    fetchRoles()
    emit('refresh')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 权限选择处理
const handlePermissionCheck = () => {
  const checkedKeys = permissionTreeRef.value.getCheckedKeys()
  const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
  form.permissionIds = [...checkedKeys, ...halfCheckedKeys]
}

// 保存角色
const saveRole = async () => {
  try {
    await formRef.value.validate()
    saving.value = true
    
    if (editingRole.value) {
      // 更新角色
      const updateData = {
        name: form.name,
        description: form.description,
        permissionIds: form.permissionIds
      }
      await systemApi.updateRole(editingRole.value.id, updateData)
      ElMessage.success('更新成功')
    } else {
      // 创建角色
      const createData = {
        name: form.name,
        code: form.name.toLowerCase().replace(/\s+/g, '_'),
        description: form.description,
        status: 'active' as 'active' | 'inactive',
        permissions: []
      }
      await systemApi.createRole(createData)
      ElMessage.success('创建成功')
    }
    
    showCreateDialog.value = false
    fetchRoles()
    emit('refresh')
  } catch (error) {
    ElMessage.error(editingRole.value ? '更新失败' : '创建失败')
  } finally {
    saving.value = false
  }
}

// 重置表单
const resetForm = () => {
  editingRole.value = null
  Object.assign(form, {
    name: '',
    description: '',
    permissionIds: []
  })
  formRef.value?.resetFields()
}

// 工具函数
const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 暴露方法给父组件
defineExpose({
  fetchRoles,
  rolesList
})

// 初始化
onMounted(() => {
  fetchRoles()
  fetchPermissionTree()
})
</script>

<style scoped>
.table-card {
  border: 1px solid var(--el-border-color-light);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 深色主题适配 */
.dark .table-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
