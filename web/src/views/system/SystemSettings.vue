<template>
  <div class="system-settings">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">
          <el-icon><Setting /></el-icon>
          系统设置
        </h1>
        <p class="page-description">管理系统配置、用户权限和安全设置</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="saveAllSettings" :loading="saving">
          <el-icon><Check /></el-icon>
          保存所有设置
        </el-button>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 设置标签页 -->
    <el-tabs v-model="activeTab" class="settings-tabs">
      <!-- 用户管理 -->
      <el-tab-pane label="用户管理" name="users">
        <div class="tab-content">
          <UserManagement 
            ref="userManagementRef" 
            :users="usersList"
            :roles="rolesList"
            :loading="usersLoading"
            @user-created="handleUserCreated"
            @user-updated="handleUserUpdated"
            @user-deleted="handleUserDeleted"
            @user-status-toggled="handleUserStatusToggled"
            @refresh-users="handleRefresh"
          />
        </div>
      </el-tab-pane>

      <!-- 角色管理 -->
      <el-tab-pane label="角色管理" name="roles">
        <div class="tab-content">
          <RoleManagement 
            ref="roleManagementRef"
            :roles="rolesList"
            :permissions="permissionsList"
            :loading="rolesLoading"
            @role-created="handleRoleCreated"
            @role-updated="handleRoleUpdated"
            @role-deleted="handleRoleDeleted"
            @refresh-roles="handleRefresh"
          />
        </div>
      </el-tab-pane>

      <!-- 系统配置 -->
      <el-tab-pane label="系统配置" name="config">
        <div class="tab-content">
          <SystemConfiguration 
            ref="systemConfigurationRef"
            :config="systemConfig"
            :loading="configLoading"
            @config-updated="handleConfigUpdated"
            @config-reset="handleConfigReset"
            @save-config="handleSaveConfig"
            @test-email="handleTestEmail"
          />
        </div>
      </el-tab-pane>

      <!-- 操作日志 -->
      <el-tab-pane label="操作日志" name="logs">
        <div class="tab-content">
          <OperationLogs 
            ref="operationLogsRef"
            @refresh="handleRefresh"
          />
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Setting,
  Check,
  Refresh
} from '@element-plus/icons-vue'
import * as systemApi from '@/api/system'
import type { Role } from '@/types/system'

// 导入子组件
import { UserManagement, RoleManagement, SystemConfiguration } from '@/components'
import OperationLogs from './components/OperationLogs.vue'

// 响应式数据
const activeTab = ref('users')
const saving = ref(false)

// 子组件引用
const userManagementRef = ref()
const roleManagementRef = ref()
const systemConfigurationRef = ref()
const operationLogsRef = ref()

// 数据列表
const usersList = ref<any[]>([])
const rolesList = ref<Role[]>([])
const permissionsList = ref<any[]>([])
const systemConfig = ref<any>({})

// 加载状态
const usersLoading = ref(false)
const rolesLoading = ref(false)
const configLoading = ref(false)

// 获取用户列表
const fetchUsersList = async () => {
  try {
    usersLoading.value = true
    const response = await systemApi.getUserList({})
    usersList.value = (response.data as any).data || []
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    usersLoading.value = false
  }
}

// 获取角色列表
const fetchRolesList = async () => {
  try {
    rolesLoading.value = true
    const response = await systemApi.getRoleList({})
    rolesList.value = (response.data as any).data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  } finally {
    rolesLoading.value = false
  }
}

// 获取权限列表
const fetchPermissionsList = async () => {
  try {
    const response = await systemApi.getPermissionList()
    permissionsList.value = (response.data as any).data || []
  } catch (error) {
    console.error('获取权限列表失败:', error)
    ElMessage.error('获取权限列表失败')
  }
}

// 获取系统配置
const fetchSystemConfig = async () => {
  try {
    configLoading.value = true
    const response = await systemApi.getSystemConfigs()
    systemConfig.value = response.data || {}
  } catch (error) {
    console.error('获取系统配置失败:', error)
    ElMessage.error('获取系统配置失败')
  } finally {
    configLoading.value = false
  }
}

// 刷新所有数据
const refreshData = () => {
  // 刷新当前标签页的数据
  switch (activeTab.value) {
    case 'users':
      userManagementRef.value?.fetchUsers()
      break
    case 'roles':
      roleManagementRef.value?.fetchRoles()
      break
    case 'config':
      systemConfigurationRef.value?.fetchConfig()
      break
    case 'logs':
      operationLogsRef.value?.fetchLogs()
      break
  }
  
  // 刷新角色列表
  fetchRolesList()
}

// 保存所有设置
const saveAllSettings = async () => {
  try {
    saving.value = true
    
    // 根据当前活动标签保存对应设置
    switch (activeTab.value) {
      case 'config':
        await systemConfigurationRef.value?.saveConfig()
        break
      default:
        ElMessage.info('当前标签页没有需要保存的配置')
        break
    }
    
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 用户管理事件处理
const handleUserCreated = (user: any) => {
  usersList.value.push(user)
  ElMessage.success('用户创建成功')
}

const handleUserUpdated = (updatedUser: any) => {
  const index = usersList.value.findIndex(u => u.id === updatedUser.id)
  if (index !== -1) {
    usersList.value[index] = updatedUser
  }
  ElMessage.success('用户更新成功')
}

const handleUserDeleted = (userId: string) => {
  usersList.value = usersList.value.filter(u => u.id !== userId)
  ElMessage.success('用户删除成功')
}

const handleUserStatusToggled = (userId: string, status: boolean) => {
  const user = usersList.value.find(u => u.id === userId)
  if (user) {
    user.is_active = status
  }
  ElMessage.success(`用户已${status ? '启用' : '禁用'}`)
}

// 角色管理事件处理
const handleRoleCreated = (role: any) => {
  rolesList.value.push(role)
  ElMessage.success('角色创建成功')
}

const handleRoleUpdated = (updatedRole: any) => {
  const index = rolesList.value.findIndex(r => r.id === updatedRole.id)
  if (index !== -1) {
    rolesList.value[index] = updatedRole
  }
  ElMessage.success('角色更新成功')
}

const handleRoleDeleted = (roleId: string) => {
  rolesList.value = rolesList.value.filter(r => r.id !== roleId)
  ElMessage.success('角色删除成功')
}

// 系统配置事件处理
const handleConfigUpdated = (config: any) => {
  systemConfig.value = config
  ElMessage.success('系统配置更新成功')
}

const handleConfigReset = () => {
  fetchSystemConfig()
  ElMessage.info('系统配置已重置')
}

const handleSaveConfig = async (config: any) => {
  try {
    await systemApi.updateSystemConfig('all', JSON.stringify(config))
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

// 处理子组件的刷新事件
const handleRefresh = () => {
  // 当子组件数据变化时，刷新相关数据
  fetchUsersList()
  fetchRolesList()
  fetchPermissionsList()
  fetchSystemConfig()
}

// 初始化
onMounted(() => {
  fetchUsersList()
  fetchRolesList()
  fetchPermissionsList()
  fetchSystemConfig()
})
</script>

<style scoped>
.system-settings {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left {
  flex: 1;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.page-description {
  margin: 0;
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 12px;
}

.settings-tabs {
  margin-bottom: 20px;
}

.tab-content {
  padding: 0;
}

/* 深色主题适配 */
.dark .system-settings {
  background: var(--el-bg-color);
}
</style>