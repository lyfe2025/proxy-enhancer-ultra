<template>
  <div class="role-management">
    <DataTable
      :data="filteredRoles"
      :columns="tableColumns"
      :loading="loading"
      :show-pagination="false"
      @edit="handleEdit"
      @delete="handleDelete"
      @refresh="handleRefresh"
    >
      <!-- 筛选栏左侧 -->
      <template #filter-left>
        <el-input
          v-model="searchQuery"
          placeholder="搜索角色名称或描述"
          prefix-icon="Search"
          clearable
          style="width: 300px; margin-right: 16px"
          @input="handleSearch"
        />
      </template>

      <!-- 筛选栏右侧 -->
      <template #filter-right>
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新增角色
        </el-button>
        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </template>

      <!-- 角色名称列 -->
      <template #name="{ row }">
        <span class="role-name">{{ row.name }}</span>
      </template>

      <!-- 用户数量列 -->
      <template #user_count="{ row }">
        <el-tag type="info">{{ row.user_count || 0 }}</el-tag>
      </template>

      <!-- 权限列 -->
      <template #permissions="{ row }">
        <div class="permission-tags">
          <el-tag
            v-for="permission in (row.permissions || []).slice(0, 3)"
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

      <!-- 创建时间列 -->
      <template #created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button size="small" @click="handleEdit(row)">编辑</el-button>
        <el-button 
          size="small" 
          type="danger" 
          @click="handleDelete(row)"
          :disabled="(row.user_count || 0) > 0"
        >
          删除
        </el-button>
      </template>
    </DataTable>

    <!-- 角色编辑对话框 -->
    <RoleFormDialog
      v-model:visible="dialogVisible"
      :role="currentRole"
      :permissions="permissions"
      :is-edit="isEdit"
      @save="handleSave"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import DataTable from '@/components/common/DataTable/DataTable.vue'
import RoleFormDialog from './RoleFormDialog.vue'
import type { Role, Permission } from '@/types/system'
import type { TableColumn } from '@/components/common/DataTable/DataTable.vue'

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
const currentRole = ref<Role | null>(null)

// 表格列配置
const tableColumns: TableColumn[] = [
  {
    prop: 'name',
    label: '角色名称',
    width: 150,
    slot: 'name'
  },
  {
    prop: 'description',
    label: '描述',
    width: 200
  },
  {
    prop: 'user_count',
    label: '用户数量',
    width: 100,
    slot: 'user_count'
  },
  {
    prop: 'permissions',
    label: '权限',
    width: 300,
    slot: 'permissions'
  },
  {
    prop: 'created_at',
    label: '创建时间',
    width: 180,
    slot: 'created_at'
  }
]

// 过滤后的角色数据
const filteredRoles = computed(() => {
  return props.roles.filter(role => {
    const matchesSearch = !searchQuery.value || 
      role.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      role.description?.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    return matchesSearch
  })
})

// 监听搜索条件变化
watch(searchQuery, () => {
  // 搜索是在客户端进行的，不需要额外处理
})

// 方法
const showCreateDialog = () => {
  isEdit.value = false
  currentRole.value = null
  dialogVisible.value = true
}

const handleEdit = (role: Role) => {
  isEdit.value = true
  currentRole.value = role
  dialogVisible.value = true
}

const handleDelete = async (role: Role) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${role.name}" 吗？此操作不可恢复。`,
      '删除角色',
      { type: 'warning' }
    )
    
    emit('role-deleted', role.id)
    ElMessage.success('角色删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除角色失败:', error)
      ElMessage.error('删除角色失败，请重试')
    }
  }
}

const handleSave = (role: Role) => {
  if (isEdit.value) {
    emit('role-updated', role)
  } else {
    emit('role-created', role)
  }
  dialogVisible.value = false
}

const handleRefresh = () => {
  emit('refresh-roles')
}

const handleSearch = () => {
  // 搜索逻辑在computed中处理
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.role-management {
  padding: 20px;
}

.role-name {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.permission-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
</style>
