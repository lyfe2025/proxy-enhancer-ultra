<template>
  <div class="user-management">
    <DataTable
      :data="users"
      :columns="tableColumns"
      :loading="loading"
      :pagination="paginationData"
      @edit="handleEdit"
      @delete="handleDelete"
      @refresh="handleRefresh"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <!-- 筛选栏左侧 -->
      <template #filter-left>
        <el-input
          v-model="searchQuery"
          placeholder="搜索用户名或邮箱"
          prefix-icon="Search"
          clearable
          style="width: 300px; margin-right: 16px"
          @input="handleSearch"
        />
        <el-select
          v-model="roleFilter"
          placeholder="筛选角色"
          clearable
          style="width: 150px; margin-right: 16px"
        >
          <el-option label="全部角色" value="" />
          <el-option
            v-for="role in roles"
            :key="role.id"
            :label="role.name"
            :value="role.id"
          />
        </el-select>
        <el-select
          v-model="statusFilter"
          placeholder="筛选状态"
          clearable
          style="width: 120px"
        >
          <el-option label="全部状态" value="" />
          <el-option label="启用" value="true" />
          <el-option label="禁用" value="false" />
        </el-select>
      </template>

      <!-- 筛选栏右侧 -->
      <template #filter-right>
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新增用户
        </el-button>
        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </template>

      <!-- 用户名列 -->
      <template #username="{ row }">
        <div class="user-info">
          <el-avatar :size="32" :src="row.avatar">
            <el-icon><User /></el-icon>
          </el-avatar>
          <span class="username">{{ row.username }}</span>
        </div>
      </template>

      <!-- 角色列 -->
      <template #role="{ row }">
        <el-tag :type="getRoleTagType(row.role?.name)">
          {{ row.role?.name || '未分配' }}
        </el-tag>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <el-switch
          v-model="row.is_active"
          @change="toggleUserStatus(row)"
          :loading="row.statusLoading"
        />
      </template>

      <!-- 创建时间列 -->
      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button size="small" @click="handleEdit(row)">编辑</el-button>
        <el-button size="small" type="warning" @click="resetPassword(row)">重置密码</el-button>
        <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
      </template>
    </DataTable>

    <!-- 用户编辑对话框 -->
    <UserFormDialog
      v-model:visible="dialogVisible"
      :user="currentUser"
      :roles="roles"
      :is-edit="isEdit"
      @save="handleSave"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, User } from '@element-plus/icons-vue'
import DataTable from '@/components/common/DataTable/DataTable.vue'
import UserFormDialog from './UserFormDialog.vue'
import type { User as UserType, Role } from '@/types/system'
import type { TableColumn, PaginationData } from '@/components/common/DataTable/DataTable.vue'

// Props
interface Props {
  users: UserType[]
  roles: Role[]
  loading: boolean
  pagination?: PaginationData
}

const props = defineProps<Props>()

// Emits
interface Emits {
  'user-created': [user: UserType]
  'user-updated': [user: UserType]
  'user-deleted': [userId: string]
  'user-status-toggled': [userId: string, status: boolean]
  'refresh-users': []
  'page-change': [page: number]
  'size-change': [size: number]
  'search': [query: string, filters: any]
}

const emit = defineEmits<Emits>()

// 响应式数据
const searchQuery = ref('')
const roleFilter = ref('')
const statusFilter = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentUser = ref<UserType | null>(null)

// 表格列配置
const tableColumns: TableColumn[] = [
  {
    prop: 'username',
    label: '用户名',
    width: 200,
    slot: 'username'
  },
  {
    prop: 'email',
    label: '邮箱',
    width: 200
  },
  {
    prop: 'role',
    label: '角色',
    width: 120,
    slot: 'role'
  },
  {
    prop: 'is_active',
    label: '状态',
    width: 100,
    slot: 'status'
  },
  {
    prop: 'created_at',
    label: '创建时间',
    width: 180,
    slot: 'created_at'
  }
]

// 分页数据
const paginationData = computed<PaginationData>(() => ({
  page: props.pagination?.page || 1,
  size: props.pagination?.size || 20,
  total: props.pagination?.total || 0
}))

// 过滤后的用户数据
const filteredUsers = computed(() => {
  return props.users.filter(user => {
    const matchesSearch = !searchQuery.value || 
      user.username.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      user.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesRole = !roleFilter.value || user.role_id === roleFilter.value
    const matchesStatus = statusFilter.value === '' || user.is_active.toString() === statusFilter.value
    
    return matchesSearch && matchesRole && matchesStatus
  })
})

// 监听筛选条件变化
watch([searchQuery, roleFilter, statusFilter], () => {
  handleSearch()
})

// 方法
const showCreateDialog = () => {
  isEdit.value = false
  currentUser.value = null
  dialogVisible.value = true
}

const handleEdit = (user: UserType) => {
  isEdit.value = true
  currentUser.value = user
  dialogVisible.value = true
}

const handleDelete = async (user: UserType) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.username}" 吗？此操作不可恢复。`,
      '删除用户',
      { type: 'warning' }
    )
    
    emit('user-deleted', user.id)
    ElMessage.success('用户删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error('删除用户失败，请重试')
    }
  }
}

const handleSave = (user: UserType) => {
  if (isEdit.value) {
    emit('user-updated', user)
  } else {
    emit('user-created', user)
  }
  dialogVisible.value = false
}

const toggleUserStatus = async (user: UserType & { statusLoading?: boolean }) => {
  try {
    user.statusLoading = true
    emit('user-status-toggled', user.id, user.is_active)
    ElMessage.success(`用户已${user.is_active ? '启用' : '禁用'}`)
  } catch (error) {
    user.is_active = !user.is_active // 回滚状态
    console.error('切换用户状态失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    user.statusLoading = false
  }
}

const resetPassword = async (user: UserType) => {
  try {
    await ElMessageBox.confirm(
      `确定要重置用户 "${user.username}" 的密码吗？`,
      '重置密码',
      { type: 'warning' }
    )
    
    // TODO: 实现API调用
    ElMessage.success('密码重置成功，新密码已发送到用户邮箱')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置密码失败:', error)
      ElMessage.error('重置密码失败，请重试')
    }
  }
}

const handleRefresh = () => {
  emit('refresh-users')
}

const handleSearch = () => {
  const filters = {
    role: roleFilter.value,
    status: statusFilter.value
  }
  emit('search', searchQuery.value, filters)
}

const handlePageChange = (page: number) => {
  emit('page-change', page)
}

const handleSizeChange = (size: number) => {
  emit('size-change', size)
}

const getRoleTagType = (roleName?: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  if (!roleName) return 'primary'
  const roleTypes: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    '超级管理员': 'danger',
    '管理员': 'warning',
    '普通用户': 'info'
  }
  return roleTypes[roleName] || 'primary'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.user-management {
  padding: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  font-weight: 500;
}
</style>
