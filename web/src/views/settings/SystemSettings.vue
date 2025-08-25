<template>
  <div class="system-settings">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>系统设置</h2>
        <p>管理用户、角色、权限和系统配置</p>
      </div>
      <div class="header-right">
        <el-button @click="showSystemInfoDialog">
          <el-icon><InfoFilled /></el-icon>
          系统信息
        </el-button>
      </div>
    </div>

    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="settings-tabs">
      <!-- 用户管理 -->
      <el-tab-pane label="用户管理" name="users">
        <UserManagement
          :users="users"
          :roles="roles"
          :loading="usersLoading"
          @user-created="handleUserCreated"
          @user-updated="handleUserUpdated"
          @user-deleted="handleUserDeleted"
          @user-status-toggled="handleUserStatusToggled"
          @refresh-users="handleUsersRefresh"
        />
      </el-tab-pane>

      <!-- 角色管理 -->
      <el-tab-pane label="角色管理" name="roles">
        <RoleManagement
          :roles="roles"
          :permissions="permissions"
          :loading="rolesLoading"
          @role-created="handleRoleCreated"
          @role-updated="handleRoleUpdated"
          @role-deleted="handleRoleDeleted"
          @refresh-roles="handleRolesRefresh"
        />
      </el-tab-pane>

      <!-- 权限管理 -->
      <el-tab-pane label="权限管理" name="permissions">
        <PermissionManagement
          :permissions="permissions"
          :loading="permissionsLoading"
          @permission-created="handlePermissionCreated"
          @permission-updated="handlePermissionUpdated"
          @permission-deleted="handlePermissionDeleted"
          @refresh="handlePermissionsRefresh"
        />
      </el-tab-pane>

      <!-- 系统配置 -->
      <el-tab-pane label="系统配置" name="config">
        <SystemConfiguration
          :config="systemConfig"
          :loading="configLoading"
          @config-updated="handleConfigUpdated"
          @config-reset="handleConfigReset"
          @save-config="handleSaveConfig"
          @test-email="handleTestEmail"
        />
      </el-tab-pane>
    </el-tabs>

    <!-- 系统信息对话框 -->
    <el-dialog
      v-model="systemInfoDialogVisible"
      title="系统信息"
      width="600px"
    >
      <div class="system-info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="系统名称">{{ systemInfo.name }}</el-descriptions-item>
          <el-descriptions-item label="系统版本">{{ systemInfo.version }}</el-descriptions-item>
          <el-descriptions-item label="Go版本">{{ systemInfo.go_version }}</el-descriptions-item>
          <el-descriptions-item label="运行时间">{{ systemInfo.uptime }}</el-descriptions-item>
          <el-descriptions-item label="CPU使用率">{{ systemInfo.cpu_usage }}%</el-descriptions-item>
          <el-descriptions-item label="内存使用率">{{ systemInfo.memory_usage }}%</el-descriptions-item>
          <el-descriptions-item label="磁盘使用率">{{ systemInfo.disk_usage }}%</el-descriptions-item>
          <el-descriptions-item label="数据库连接">{{ systemInfo.db_status }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'
import { UserManagement, RoleManagement, SystemConfiguration } from '@/components'
import PermissionManagement from './components/PermissionManagement.vue'
import { getUserList, getRoleList, getPermissionList, getSystemStats } from '@/api/system'
import { useAuthStore } from '@/store'
import type { User, Role, Permission, SystemConfig, SystemInfo } from '@/types/system'

// Store
const authStore = useAuthStore()
const currentUser = computed(() => authStore.user)

// 响应式数据
const activeTab = ref('users')
const systemInfoDialogVisible = ref(false)

// 数据状态
const users = ref<User[]>([])
const roles = ref<Role[]>([])
const permissions = ref<Permission[]>([])
const systemConfig = ref<SystemConfig>({} as SystemConfig)
const systemInfo = ref<SystemInfo>({} as SystemInfo)

// 加载状态
const usersLoading = ref(false)
const rolesLoading = ref(false)
const permissionsLoading = ref(false)
const configLoading = ref(false)

// 数据获取方法
const fetchUsers = async () => {
  try {
    usersLoading.value = true
    const response = await getUserList({ page: 1, pageSize: 1000 })
    const responseData = response.data as any
    users.value = Array.isArray(responseData) ? responseData : (responseData?.data || [])
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  } finally {
    usersLoading.value = false
  }
}

const fetchRoles = async () => {
  try {
    rolesLoading.value = true
    const response = await getRoleList({ page: 1, pageSize: 1000 })
    const responseData = response.data as any
    roles.value = Array.isArray(responseData) ? responseData : (responseData?.data || [])
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  } finally {
    rolesLoading.value = false
  }
}

const fetchPermissions = async () => {
  try {
    permissionsLoading.value = true
    const response = await getPermissionList()
    permissions.value = response.data?.data || response.data || []
  } catch (error) {
    ElMessage.error('获取权限列表失败')
  } finally {
    permissionsLoading.value = false
  }
}

const fetchSystemInfo = async () => {
  try {
    const response = await getSystemStats()
    systemInfo.value = response.data?.data || response.data || {
      name: '',
      version: '',
      go_version: '',
      build_time: '',
      git_commit: '',
      uptime: '',
      cpu_usage: 0,
      memory_usage: 0,
      memory_total: 0,
      memory_used: 0,
      disk_usage: 0,
      disk_total: 0,
      disk_used: 0,
      db_status: '',
      db_version: '',
      online_users: 0,
      total_users: 0,
      total_roles: 0,
      total_permissions: 0,
      system_load: [],
      network_in: 0,
      network_out: 0
    }
  } catch (error) {
    ElMessage.error('获取系统信息失败')
  }
}

// 系统信息对话框
const showSystemInfoDialog = async () => {
  await fetchSystemInfo()
  systemInfoDialogVisible.value = true
}

// 用户管理事件处理
const handleUserCreated = (user: User) => {
  users.value.push(user)
  ElMessage.success('用户创建成功')
}

const handleUserUpdated = (updatedUser: User) => {
  const index = users.value.findIndex(u => u.id === updatedUser.id)
  if (index !== -1) {
    users.value[index] = updatedUser
  }
  ElMessage.success('用户更新成功')
}

const handleUserDeleted = (userId: string) => {
  users.value = users.value.filter(u => u.id !== userId)
  ElMessage.success('用户删除成功')
}

const handleUserStatusToggled = (userId: string, status: boolean) => {
  const user = users.value.find(u => u.id === userId)
  if (user) {
    user.is_active = status
  }
  ElMessage.success(`用户已${status ? '启用' : '禁用'}`)
}

const handlePasswordReset = (userId: string) => {
  ElMessage.success('密码重置成功')
}

const handleUsersRefresh = () => {
  fetchUsers()
}

// 角色管理事件处理
const handleRoleCreated = (role: Role) => {
  roles.value.push(role)
  ElMessage.success('角色创建成功')
}

const handleRoleUpdated = (updatedRole: Role) => {
  const index = roles.value.findIndex(r => r.id === updatedRole.id)
  if (index !== -1) {
    roles.value[index] = updatedRole
  }
  ElMessage.success('角色更新成功')
}

const handleRoleDeleted = (roleId: string) => {
  roles.value = roles.value.filter(r => r.id !== roleId)
  ElMessage.success('角色删除成功')
}

const handleRolesRefresh = () => {
  fetchRoles()
}

// 权限管理事件处理
const handlePermissionCreated = (permission: Permission) => {
  permissions.value.push(permission)
  ElMessage.success('权限创建成功')
}

const handlePermissionUpdated = (updatedPermission: Permission) => {
  const index = permissions.value.findIndex(p => p.id === updatedPermission.id)
  if (index !== -1) {
    permissions.value[index] = updatedPermission
  }
  ElMessage.success('权限更新成功')
}

const handlePermissionDeleted = (permissionId: string) => {
  permissions.value = permissions.value.filter(p => p.id !== permissionId)
  ElMessage.success('权限删除成功')
}

const handlePermissionsRefresh = () => {
  fetchPermissions()
}

// 系统配置事件处理
const handleConfigUpdated = (config: SystemConfig) => {
  systemConfig.value = config
  ElMessage.success('系统配置更新成功')
}

const handleConfigReset = () => {
  // 重新获取配置数据
  // fetchSystemConfig() // TODO: 实现配置获取
  ElMessage.info('系统配置已重置')
}

const handleSaveConfig = async (config: SystemConfig) => {
  try {
    // TODO: 实现保存配置API
    systemConfig.value = config
    ElMessage.success('系统配置保存成功')
  } catch (error) {
    console.error('保存系统配置失败:', error)
    ElMessage.error('保存系统配置失败')
  }
}

const handleTestEmail = async (emailConfig: any) => {
  try {
    // TODO: 实现邮件测试API
    ElMessage.info('邮件配置测试功能开发中')
  } catch (error) {
    console.error('邮件测试失败:', error)
    ElMessage.error('邮件测试失败')
  }
}

// 初始化
onMounted(() => {
  fetchUsers()
  fetchRoles()
  fetchPermissions()
})
</script>

<style scoped>
.system-settings {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.header-left h2 {
  margin: 0 0 8px 0;
  color: #303133;
  font-size: 24px;
  font-weight: 600;
}

.header-left p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.settings-tabs {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.system-info {
  padding: 20px;
}
</style>