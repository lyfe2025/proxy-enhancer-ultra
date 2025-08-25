<template>
  <div class="system-settings">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>系统设置</h2>
        <p>管理用户、角色、权限和系统配置</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateUserDialog">
          <el-icon><Plus /></el-icon>
          新增用户
        </el-button>
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
        <!-- 用户筛选 -->
        <div class="filter-bar">
          <div class="filter-left">
            <el-input
              v-model="userSearchQuery"
              placeholder="搜索用户名或邮箱"
              style="width: 300px"
              clearable
              @input="handleUserSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            
            <el-select
              v-model="userRoleFilter"
              placeholder="角色筛选"
              style="width: 150px"
              clearable
              @change="handleUserFilter"
            >
              <el-option
                v-for="role in roles"
                :key="role.id"
                :label="role.name"
                :value="role.id"
              />
            </el-select>
            
            <el-select
              v-model="userStatusFilter"
              placeholder="状态筛选"
              style="width: 120px"
              clearable
              @change="handleUserFilter"
            >
              <el-option label="启用" value="active" />
              <el-option label="禁用" value="inactive" />
            </el-select>
          </div>
          
          <div class="filter-right">
            <el-button @click="refreshUsers">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>

        <!-- 用户列表 -->
        <div class="users-table">
          <el-table
            v-loading="usersLoading"
            :data="filteredUsers"
            row-key="id"
            stripe
            style="width: 100%"
          >
            <el-table-column prop="username" label="用户名" min-width="120">
              <template #default="{ row }">
                <div class="user-info">
                  <el-avatar :size="32" :src="row.avatar" :alt="row.username">
                    <el-icon><User /></el-icon>
                  </el-avatar>
                  <span class="username">{{ row.username }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="email" label="邮箱" min-width="180">
              <template #default="{ row }">
                <span class="email">{{ row.email }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="role_name" label="角色" width="120">
              <template #default="{ row }">
                <el-tag :type="getRoleTagType(row.role_name)">{{ row.role_name }}</el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-switch
                  v-model="row.is_active"
                  @change="toggleUserStatus(row)"
                  :loading="row.toggling"
                  active-color="#00ff88"
                  inactive-color="#666"
                />
              </template>
            </el-table-column>
            
            <el-table-column prop="last_login" label="最后登录" width="180">
              <template #default="{ row }">
                <span v-if="row.last_login" class="last-login">
                  {{ formatDate(row.last_login) }}
                </span>
                <span v-else class="never-login">从未登录</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="text" size="small" @click="editUser(row)">
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-button>
                  
                  <el-button type="text" size="small" @click="resetPassword(row)">
                    <el-icon><Key /></el-icon>
                    重置密码
                  </el-button>
                  
                  <el-button
                    type="text"
                    size="small"
                    @click="deleteUser(row)"
                    style="color: #f56c6c"
                    :disabled="row.id === currentUser.id"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- 角色管理 -->
      <el-tab-pane label="角色管理" name="roles">
        <!-- 角色操作栏 -->
        <div class="filter-bar">
          <div class="filter-left">
            <el-input
              v-model="roleSearchQuery"
              placeholder="搜索角色名称或描述"
              style="width: 300px"
              clearable
              @input="handleRoleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
          
          <div class="filter-right">
            <el-button type="primary" @click="showCreateRoleDialog">
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
        <div class="roles-table">
          <el-table
            v-loading="rolesLoading"
            :data="filteredRoles"
            row-key="id"
            stripe
            style="width: 100%"
          >
            <el-table-column prop="name" label="角色名称" width="150">
              <template #default="{ row }">
                <div class="role-name">
                  <el-icon class="role-icon"><UserFilled /></el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="description" label="描述" min-width="200">
              <template #default="{ row }">
                <span class="role-description">{{ row.description || '暂无描述' }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="user_count" label="用户数量" width="100">
              <template #default="{ row }">
                <el-tag size="small">{{ row.user_count || 0 }}</el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="permissions" label="权限" min-width="300">
              <template #default="{ row }">
                <div class="permissions-list">
                  <el-tag
                    v-for="permission in row.permissions?.slice(0, 3)"
                    :key="permission.id"
                    size="small"
                    class="permission-tag"
                  >
                    {{ permission.name }}
                  </el-tag>
                  <el-tag
                    v-if="row.permissions?.length > 3"
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
            
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="text" size="small" @click="editRole(row)">
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-button>
                  
                  <el-button
                    type="text"
                    size="small"
                    @click="deleteRole(row)"
                    style="color: #f56c6c"
                    :disabled="row.is_system"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- 权限管理 -->
      <el-tab-pane label="权限管理" name="permissions">
        <!-- 权限操作栏 -->
        <div class="filter-bar">
          <div class="filter-left">
            <el-input
              v-model="permissionSearchQuery"
              placeholder="搜索权限名称或描述"
              style="width: 300px"
              clearable
              @input="handlePermissionSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            
            <el-select
              v-model="permissionModuleFilter"
              placeholder="模块筛选"
              style="width: 150px"
              clearable
              @change="handlePermissionFilter"
            >
              <el-option label="用户管理" value="user" />
              <el-option label="代理管理" value="proxy" />
              <el-option label="规则管理" value="rule" />
              <el-option label="数据收集" value="data" />
              <el-option label="系统设置" value="system" />
            </el-select>
          </div>
          
          <div class="filter-right">
            <el-button type="primary" @click="showCreatePermissionDialog">
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
        <div class="permissions-table">
          <el-table
            v-loading="permissionsLoading"
            :data="filteredPermissions"
            row-key="id"
            stripe
            style="width: 100%"
          >
            <el-table-column prop="name" label="权限名称" width="200">
              <template #default="{ row }">
                <div class="permission-name">
                  <el-icon class="permission-icon"><Lock /></el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="code" label="权限代码" width="180">
              <template #default="{ row }">
                <code class="permission-code">{{ row.code }}</code>
              </template>
            </el-table-column>
            
            <el-table-column prop="module" label="所属模块" width="120">
              <template #default="{ row }">
                <el-tag :type="getModuleTagType(row.module)">{{ getModuleText(row.module) }}</el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="description" label="描述" min-width="200">
              <template #default="{ row }">
                <span class="permission-description">{{ row.description || '暂无描述' }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="text" size="small" @click="editPermission(row)">
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-button>
                  
                  <el-button
                    type="text"
                    size="small"
                    @click="deletePermission(row)"
                    style="color: #f56c6c"
                    :disabled="row.is_system"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- 系统配置 -->
      <el-tab-pane label="系统配置" name="config">
        <div class="config-sections">
          <!-- 基础配置 -->
          <div class="config-section">
            <div class="section-header">
              <h3>基础配置</h3>
              <p>系统基本设置和参数配置</p>
            </div>
            
            <el-form
              ref="configFormRef"
              :model="configForm"
              label-width="150px"
              class="config-form"
            >
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="系统名称">
                    <el-input v-model="configForm.system_name" placeholder="请输入系统名称" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="系统版本">
                    <el-input v-model="configForm.system_version" placeholder="请输入系统版本" readonly />
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="会话超时时间">
                    <el-input-number
                      v-model="configForm.session_timeout"
                      :min="5"
                      :max="1440"
                      style="width: 100%"
                    />
                    <div class="form-tip">分钟，建议设置为30-120分钟</div>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="最大登录尝试次数">
                    <el-input-number
                      v-model="configForm.max_login_attempts"
                      :min="3"
                      :max="10"
                      style="width: 100%"
                    />
                    <div class="form-tip">次，超过后将锁定账户</div>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="密码最小长度">
                    <el-input-number
                      v-model="configForm.password_min_length"
                      :min="6"
                      :max="20"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="启用双因子认证">
                    <el-switch
                      v-model="configForm.enable_2fa"
                      active-color="#00ff88"
                      inactive-color="#666"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-form-item>
                <el-button type="primary" @click="saveConfig" :loading="configSaving">
                  保存配置
                </el-button>
                <el-button @click="resetConfig">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 日志配置 -->
          <div class="config-section">
            <div class="section-header">
              <h3>日志配置</h3>
              <p>系统日志记录和存储设置</p>
            </div>
            
            <el-form
              :model="logConfigForm"
              label-width="150px"
              class="config-form"
            >
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="日志级别">
                    <el-select v-model="logConfigForm.log_level" style="width: 100%">
                      <el-option label="DEBUG" value="debug" />
                      <el-option label="INFO" value="info" />
                      <el-option label="WARN" value="warn" />
                      <el-option label="ERROR" value="error" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="日志保留天数">
                    <el-input-number
                      v-model="logConfigForm.log_retention_days"
                      :min="1"
                      :max="365"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="启用访问日志">
                    <el-switch
                      v-model="logConfigForm.enable_access_log"
                      active-color="#00ff88"
                      inactive-color="#666"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="启用错误日志">
                    <el-switch
                      v-model="logConfigForm.enable_error_log"
                      active-color="#00ff88"
                      inactive-color="#666"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-form-item>
                <el-button type="primary" @click="saveLogConfig" :loading="logConfigSaving">
                  保存日志配置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 创建/编辑用户对话框 -->
    <el-dialog
      v-model="userDialogVisible"
      :title="isEditUser ? '编辑用户' : '新增用户'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userFormRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="userForm.username" placeholder="请输入用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="userForm.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20" v-if="!isEditUser">
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input
                v-model="userForm.password"
                type="password"
                placeholder="请输入密码"
                show-password
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="确认密码" prop="confirm_password">
              <el-input
                v-model="userForm.confirm_password"
                type="password"
                placeholder="请确认密码"
                show-password
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="角色" prop="role_id">
              <el-select v-model="userForm.role_id" placeholder="请选择角色" style="width: 100%">
                <el-option
                  v-for="role in roles"
                  :key="role.id"
                  :label="role.name"
                  :value="role.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch
                v-model="userForm.is_active"
                active-text="启用"
                inactive-text="禁用"
                active-color="#00ff88"
                inactive-color="#666"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="userDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitUserForm" :loading="userSubmitting">
            {{ isEditUser ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 创建/编辑角色对话框 -->
    <el-dialog
      v-model="roleDialogVisible"
      :title="isEditRole ? '编辑角色' : '新增角色'"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="roleFormRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="角色名称" prop="name">
              <el-input v-model="roleForm.name" placeholder="请输入角色名称" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="roleForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
          />
        </el-form-item>
        
        <el-form-item label="权限" prop="permission_ids">
          <div class="permissions-selector">
            <div class="permissions-header">
              <el-checkbox
                v-model="selectAllPermissions"
                @change="handleSelectAllPermissions"
                :indeterminate="isIndeterminate"
              >
                全选
              </el-checkbox>
              <span class="selected-count">已选择 {{ roleForm.permission_ids.length }} 个权限</span>
            </div>
            
            <div class="permissions-groups">
              <div v-for="(group, module) in groupedPermissions" :key="module" class="permission-group">
                <div class="group-header">
                  <el-checkbox
                    :model-value="isModuleAllSelected(module)"
                    @change="(checked) => handleModuleSelectAll(module, checked)"
                    :indeterminate="isModuleIndeterminate(module)"
                  >
                    {{ getModuleText(module) }}
                  </el-checkbox>
                </div>
                
                <div class="group-permissions">
                  <el-checkbox
                    v-for="permission in group"
                    :key="permission.id"
                    v-model="roleForm.permission_ids"
                    :label="permission.id"
                    :value="permission.id"
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
        <div class="dialog-footer">
          <el-button @click="roleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitRoleForm" :loading="roleSubmitting">
            {{ isEditRole ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 创建/编辑权限对话框 -->
    <el-dialog
      v-model="permissionDialogVisible"
      :title="isEditPermission ? '编辑权限' : '新增权限'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="permissionFormRef"
        :model="permissionForm"
        :rules="permissionFormRules"
        label-width="100px"
      >
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="permissionForm.name" placeholder="请输入权限名称" />
        </el-form-item>
        
        <el-form-item label="权限代码" prop="code">
          <el-input v-model="permissionForm.code" placeholder="例如: user:create" />
          <div class="form-tip">格式: 模块:操作，例如 user:create, proxy:read</div>
        </el-form-item>
        
        <el-form-item label="所属模块" prop="module">
          <el-select v-model="permissionForm.module" placeholder="请选择模块" style="width: 100%">
            <el-option label="用户管理" value="user" />
            <el-option label="代理管理" value="proxy" />
            <el-option label="规则管理" value="rule" />
            <el-option label="数据收集" value="data" />
            <el-option label="系统设置" value="system" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="permissionForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入权限描述"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="permissionDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPermissionForm" :loading="permissionSubmitting">
            {{ isEditPermission ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Delete,
  Edit,
  User,
  UserFilled,
  Lock,
  Key,
  InfoFilled
} from '@element-plus/icons-vue'
import { usersApi, rolesApi, permissionsApi, systemApi } from '@/api'
import { useAuthStore } from '@/store'
import type { User, Role, Permission } from '@/types'

// Store
const authStore = useAuthStore()
const currentUser = computed(() => authStore.user)

// 响应式数据
const activeTab = ref('users')
const usersLoading = ref(false)
const rolesLoading = ref(false)
const permissionsLoading = ref(false)
const userSubmitting = ref(false)
const roleSubmitting = ref(false)
const permissionSubmitting = ref(false)
const configSaving = ref(false)
const logConfigSaving = ref(false)

// 对话框状态
const userDialogVisible = ref(false)
const roleDialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const systemInfoDialogVisible = ref(false)
const isEditUser = ref(false)
const isEditRole = ref(false)
const isEditPermission = ref(false)

// 搜索和筛选
const userSearchQuery = ref('')
const userRoleFilter = ref('')
const userStatusFilter = ref('')
const roleSearchQuery = ref('')
const permissionSearchQuery = ref('')
const permissionModuleFilter = ref('')

// 数据
const users = ref<User[]>([])
const roles = ref<Role[]>([])
const permissions = ref<Permission[]>([])
const systemInfo = ref({})

// 权限选择
const selectAllPermissions = ref(false)
const isIndeterminate = ref(false)

// 表单引用和数据
const userFormRef = ref()
const roleFormRef = ref()
const permissionFormRef = ref()
const configFormRef = ref()

const userForm = reactive({
  id: '',
  username: '',
  email: '',
  password: '',
  confirm_password: '',
  role_id: '',
  is_active: true
})

const roleForm = reactive({
  id: '',
  name: '',
  description: '',
  permission_ids: [] as string[]
})

const permissionForm = reactive({
  id: '',
  name: '',
  code: '',
  module: '',
  description: ''
})

const configForm = reactive({
  system_name: 'Proxy Enhancer Ultra',
  system_version: '1.0.0',
  session_timeout: 60,
  max_login_attempts: 5,
  password_min_length: 8,
  enable_2fa: false
})

const logConfigForm = reactive({
  log_level: 'info',
  log_retention_days: 30,
  enable_access_log: true,
  enable_error_log: true
})

// 表单验证规则
const userFormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== userForm.password) {
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

const roleFormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ]
}

const permissionFormRules = {
  name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入权限代码', trigger: 'blur' },
    { pattern: /^[a-z]+:[a-z]+$/, message: '格式应为 模块:操作', trigger: 'blur' }
  ],
  module: [
    { required: true, message: '请选择所属模块', trigger: 'change' }
  ]
}

// 计算属性
const filteredUsers = computed(() => {
  let result = users.value
  
  // 搜索过滤
  if (userSearchQuery.value) {
    const query = userSearchQuery.value.toLowerCase()
    result = result.filter(user => 
      user.username.toLowerCase().includes(query) ||
      user.email.toLowerCase().includes(query)
    )
  }
  
  // 角色过滤
  if (userRoleFilter.value) {
    result = result.filter(user => user.role_id === userRoleFilter.value)
  }
  
  // 状态过滤
  if (userStatusFilter.value) {
    const isActive = userStatusFilter.value === 'active'
    result = result.filter(user => user.is_active === isActive)
  }
  
  return result
})

const filteredRoles = computed(() => {
  let result = roles.value
  
  if (roleSearchQuery.value) {
    const query = roleSearchQuery.value.toLowerCase()
    result = result.filter(role => 
      role.name.toLowerCase().includes(query) ||
      (role.description && role.description.toLowerCase().includes(query))
    )
  }
  
  return result
})

const filteredPermissions = computed(() => {
  let result = permissions.value
  
  // 搜索过滤
  if (permissionSearchQuery.value) {
    const query = permissionSearchQuery.value.toLowerCase()
    result = result.filter(permission => 
      permission.name.toLowerCase().includes(query) ||
      permission.code.toLowerCase().includes(query) ||
      (permission.description && permission.description.toLowerCase().includes(query))
    )
  }
  
  // 模块过滤
  if (permissionModuleFilter.value) {
    result = result.filter(permission => permission.module === permissionModuleFilter.value)
  }
  
  return result
})

const groupedPermissions = computed(() => {
  const groups: Record<string, Permission[]> = {}
  
  permissions.value.forEach(permission => {
    if (!groups[permission.module]) {
      groups[permission.module] = []
    }
    groups[permission.module].push(permission)
  })
  
  return groups
})

// 获取数据
const fetchUsers = async () => {
  try {
    usersLoading.value = true
    const response = await usersApi.getUsers()
    
    users.value = response.data.map(user => ({
      ...user,
      toggling: false
    }))
    
  } catch (error: any) {
    ElMessage.error(error.message || '获取用户列表失败')
  } finally {
    usersLoading.value = false
  }
}

const fetchRoles = async () => {
  try {
    rolesLoading.value = true
    const response = await rolesApi.getRoles()
    roles.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '获取角色列表失败')
  } finally {
    rolesLoading.value = false
  }
}

const fetchPermissions = async () => {
  try {
    permissionsLoading.value = true
    const response = await permissionsApi.getPermissions()
    permissions.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '获取权限列表失败')
  } finally {
    permissionsLoading.value = false
  }
}

const fetchSystemInfo = async () => {
  try {
    const response = await systemApi.getSystemInfo()
    systemInfo.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '获取系统信息失败')
  }
}

// 搜索处理
const handleUserSearch = () => {
  // 实时搜索，无需额外处理
}

const handleRoleSearch = () => {
  // 实时搜索，无需额外处理
}

const handlePermissionSearch = () => {
  // 实时搜索，无需额外处理
}

// 筛选处理
const handleUserFilter = () => {
  // 实时筛选，无需额外处理
}

const handlePermissionFilter = () => {
  // 实时筛选，无需额外处理
}

// 刷新数据
const refreshUsers = () => {
  fetchUsers()
}

const refreshRoles = () => {
  fetchRoles()
}

const refreshPermissions = () => {
  fetchPermissions()
}

// 用户管理
const showCreateUserDialog = () => {
  isEditUser.value = false
  resetUserForm()
  userDialogVisible.value = true
}

const editUser = (user: User) => {
  isEditUser.value = true
  Object.assign(userForm, {
    id: user.id,
    username: user.username,
    email: user.email,
    password: '',
    confirm_password: '',
    role_id: user.role_id,
    is_active: user.is_active
  })
  userDialogVisible.value = true
}

const resetUserForm = () => {
  Object.assign(userForm, {
    id: '',
    username: '',
    email: '',
    password: '',
    confirm_password: '',
    role_id: '',
    is_active: true
  })
  
  if (userFormRef.value) {
    userFormRef.value.clearValidate()
  }
}

const submitUserForm = async () => {
  if (!userFormRef.value) return
  
  try {
    await userFormRef.value.validate()
    userSubmitting.value = true
    
    const data = { ...userForm }
    delete data.id
    delete data.confirm_password
    
    if (isEditUser.value) {
      delete data.password // 编辑时不更新密码
      await usersApi.updateUser(userForm.id, data)
      ElMessage.success('用户更新成功')
    } else {
      await usersApi.createUser(data)
      ElMessage.success('用户创建成功')
    }
    
    userDialogVisible.value = false
    fetchUsers()
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    userSubmitting.value = false
  }
}

const toggleUserStatus = async (user: User) => {
  try {
    user.toggling = true
    
    if (user.is_active) {
      await usersApi.enableUser(user.id)
      ElMessage.success('用户已启用')
    } else {
      await usersApi.disableUser(user.id)
      ElMessage.success('用户已禁用')
    }
    
  } catch (error: any) {
    user.is_active = !user.is_active // 回滚状态
    ElMessage.error(error.message || '操作失败')
  } finally {
    user.toggling = false
  }
}

const resetPassword = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要重置用户 "${user.username}" 的密码吗？新密码将发送到用户邮箱。`,
      '确认重置密码',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await usersApi.resetPassword(user.id)
    ElMessage.success('密码重置成功，新密码已发送到用户邮箱')
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '重置密码失败')
    }
  }
}

const deleteUser = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.username}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await usersApi.deleteUser(user.id)
    ElMessage.success('用户删除成功')
    fetchUsers()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 角色管理
const showCreateRoleDialog = () => {
  isEditRole.value = false
  resetRoleForm()
  roleDialogVisible.value = true
}

const editRole = (role: Role) => {
  isEditRole.value = true
  Object.assign(roleForm, {
    id: role.id,
    name: role.name,
    description: role.description,
    permission_ids: role.permissions?.map(p => p.id) || []
  })
  updatePermissionSelection()
  roleDialogVisible.value = true
}

const resetRoleForm = () => {
  Object.assign(roleForm, {
    id: '',
    name: '',
    description: '',
    permission_ids: []
  })
  
  selectAllPermissions.value = false
  isIndeterminate.value = false
  
  if (roleFormRef.value) {
    roleFormRef.value.clearValidate()
  }
}

const submitRoleForm = async () => {
  if (!roleFormRef.value) return
  
  try {
    await roleFormRef.value.validate()
    roleSubmitting.value = true
    
    const data = { ...roleForm }
    delete data.id
    
    if (isEditRole.value) {
      await rolesApi.updateRole(roleForm.id, data)
      ElMessage.success('角色更新成功')
    } else {
      await rolesApi.createRole(data)
      ElMessage.success('角色创建成功')
    }
    
    roleDialogVisible.value = false
    fetchRoles()
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    roleSubmitting.value = false
  }
}

const deleteRole = async (role: Role) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${role.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await rolesApi.deleteRole(role.id)
    ElMessage.success('角色删除成功')
    fetchRoles()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 权限选择处理
const updatePermissionSelection = () => {
  const selectedCount = roleForm.permission_ids.length
  const totalCount = permissions.value.length
  
  selectAllPermissions.value = selectedCount === totalCount
  isIndeterminate.value = selectedCount > 0 && selectedCount < totalCount
}

const handleSelectAllPermissions = (checked: boolean) => {
  if (checked) {
    roleForm.permission_ids = permissions.value.map(p => p.id)
  } else {
    roleForm.permission_ids = []
  }
  isIndeterminate.value = false
}

const isModuleAllSelected = (module: string) => {
  const modulePermissions = groupedPermissions.value[module] || []
  return modulePermissions.every(p => roleForm.permission_ids.includes(p.id))
}

const isModuleIndeterminate = (module: string) => {
  const modulePermissions = groupedPermissions.value[module] || []
  const selectedInModule = modulePermissions.filter(p => roleForm.permission_ids.includes(p.id))
  return selectedInModule.length > 0 && selectedInModule.length < modulePermissions.length
}

const handleModuleSelectAll = (module: string, checked: boolean) => {
  const modulePermissions = groupedPermissions.value[module] || []
  
  if (checked) {
    // 添加该模块的所有权限
    modulePermissions.forEach(p => {
      if (!roleForm.permission_ids.includes(p.id)) {
        roleForm.permission_ids.push(p.id)
      }
    })
  } else {
    // 移除该模块的所有权限
    roleForm.permission_ids = roleForm.permission_ids.filter(id => 
      !modulePermissions.some(p => p.id === id)
    )
  }
  
  updatePermissionSelection()
}

// 权限管理
const showCreatePermissionDialog = () => {
  isEditPermission.value = false
  resetPermissionForm()
  permissionDialogVisible.value = true
}

const editPermission = (permission: Permission) => {
  isEditPermission.value = true
  Object.assign(permissionForm, {
    id: permission.id,
    name: permission.name,
    code: permission.code,
    module: permission.module,
    description: permission.description
  })
  permissionDialogVisible.value = true
}

const resetPermissionForm = () => {
  Object.assign(permissionForm, {
    id: '',
    name: '',
    code: '',
    module: '',
    description: ''
  })
  
  if (permissionFormRef.value) {
    permissionFormRef.value.clearValidate()
  }
}

const submitPermissionForm = async () => {
  if (!permissionFormRef.value) return
  
  try {
    await permissionFormRef.value.validate()
    permissionSubmitting.value = true
    
    const data = { ...permissionForm }
    delete data.id
    
    if (isEditPermission.value) {
      await permissionsApi.updatePermission(permissionForm.id, data)
      ElMessage.success('权限更新成功')
    } else {
      await permissionsApi.createPermission(data)
      ElMessage.success('权限创建成功')
    }
    
    permissionDialogVisible.value = false
    fetchPermissions()
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    permissionSubmitting.value = false
  }
}

const deletePermission = async (permission: Permission) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除权限 "${permission.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await permissionsApi.deletePermission(permission.id)
    ElMessage.success('权限删除成功')
    fetchPermissions()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 系统配置
const saveConfig = async () => {
  try {
    configSaving.value = true
    await systemApi.updateConfig(configForm)
    ElMessage.success('配置保存成功')
  } catch (error: any) {
    ElMessage.error(error.message || '保存配置失败')
  } finally {
    configSaving.value = false
  }
}

const resetConfig = () => {
  Object.assign(configForm, {
    system_name: 'Proxy Enhancer Ultra',
    system_version: '1.0.0',
    session_timeout: 60,
    max_login_attempts: 5,
    password_min_length: 8,
    enable_2fa: false
  })
}

const saveLogConfig = async () => {
  try {
    logConfigSaving.value = true
    await systemApi.updateLogConfig(logConfigForm)
    ElMessage.success('日志配置保存成功')
  } catch (error: any) {
    ElMessage.error(error.message || '保存日志配置失败')
  } finally {
    logConfigSaving.value = false
  }
}

// 显示系统信息
const showSystemInfoDialog = () => {
  fetchSystemInfo()
  systemInfoDialogVisible.value = true
}

// 工具函数
const getRoleTagType = (roleName: string) => {
  const typeMap: Record<string, string> = {
    '超级管理员': 'danger',
    '管理员': 'warning',
    '操作员': 'success',
    '普通用户': 'info'
  }
  return typeMap[roleName] || 'info'
}

const getModuleTagType = (module: string) => {
  const typeMap: Record<string, string> = {
    user: 'primary',
    proxy: 'success',
    rule: 'warning',
    data: 'info',
    system: 'danger'
  }
  return typeMap[module] || 'info'
}

const getModuleText = (module: string) => {
  const textMap: Record<string, string> = {
    user: '用户管理',
    proxy: '代理管理',
    rule: '规则管理',
    data: '数据收集',
    system: '系统设置'
  }
  return textMap[module] || module
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载
onMounted(() => {
  fetchUsers()
  fetchRoles()
  fetchPermissions()
})
</script>

<style scoped>
.system-settings {
  padding: 0;
}

/* 页面头部 */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding: 24px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%);
  border: 1px solid #333;
  border-radius: 12px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  color: #fff;
  font-size: 24px;
  font-weight: 600;
}

.header-left p {
  margin: 0;
  color: #888;
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* 标签页 */
.settings-tabs {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  overflow: hidden;
}

/* 筛选栏 */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 16px 20px;
  background-color: #2a2a2a;
  border-bottom: 1px solid #333;
}

.filter-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 表格 */
.users-table,
.roles-table,
.permissions-table {
  padding: 0 20px 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  font-weight: 500;
  color: #fff;
}

.email {
  color: #888;
}

.last-login {
  color: #888;
  font-size: 12px;
}

.never-login {
  color: #666;
  font-size: 12px;
  font-style: italic;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.role-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.role-icon {
  color: #00ff88;
}

.role-description {
  color: #888;
}

.permissions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.permission-tag {
  margin: 0;
}

.permission-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.permission-icon {
  color: #00ff88;
}

.permission-code {
  background-color: #333;
  color: #00ff88;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.permission-description {
  color: #888;
}

/* 配置部分 */
.config-sections {
  padding: 20px;
}

.config-section {
  margin-bottom: 40px;
  background-color: #2a2a2a;
  border: 1px solid #333;
  border-radius: 8px;
  overflow: hidden;
}

.section-header {
  padding: 20px;
  background-color: #333;
  border-bottom: 1px solid #444;
}

.section-header h3 {
  margin: 0 0 8px 0;
  color: #fff;
  font-size: 18px;
  font-weight: 600;
}

.section-header p {
  margin: 0;
  color: #888;
  font-size: 14px;
}

.config-form {
  padding: 20px;
}

.form-tip {
  margin-top: 4px;
  color: #666;
  font-size: 12px;
}

/* 对话框 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.permissions-selector {
  border: 1px solid #333;
  border-radius: 6px;
  background-color: #2a2a2a;
}

.permissions-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background-color: #333;
  border-bottom: 1px solid #444;
}

.selected-count {
  color: #888;
  font-size: 12px;
}

.permissions-groups {
  max-height: 300px;
  overflow-y: auto;
  padding: 16px;
}

.permission-group {
  margin-bottom: 16px;
}

.permission-group:last-child {
  margin-bottom: 0;
}

.group-header {
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid #333;
}

.group-permissions {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 8px;
  padding-left: 20px;
}

/* 系统信息 */
.system-info {
  color: #fff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .header-right {
    width: 100%;
    justify-content: flex-start;
  }
  
  .filter-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .filter-left {
    width: 100%;
    flex-wrap: wrap;
  }
  
  .filter-right {
    width: 100%;
    justify-content: flex-start;
  }
  
  .action-buttons {
    flex-direction: column;
    gap: 4px;
  }
  
  .group-permissions {
    grid-template-columns: 1fr;
  }
}

/* 深色主题样式 */
:deep(.el-table) {
  background-color: #1a1a1a;
  color: #fff;
}

:deep(.el-table__header) {
  background-color: #2a2a2a;
}

:deep(.el-table th) {
  background-color: #2a2a2a;
  color: #fff;
  border-bottom: 1px solid #333;
}

:deep(.el-table td) {
  border-bottom: 1px solid #333;
}

:deep(.el-table__row) {
  background-color: #1a1a1a;
}

:deep(.el-table__row:hover) {
  background-color: #2a2a2a;
}

:deep(.el-table__row--striped) {
  background-color: #222;
}

:deep(.el-table__row--striped:hover) {
  background-color: #2a2a2a;
}

:deep(.el-tabs__header) {
  background-color: #2a2a2a;
  margin: 0;
}

:deep(.el-tabs__nav-wrap) {
  background-color: #2a2a2a;
}

:deep(.el-tabs__item) {
  color: #888;
  border-bottom: 1px solid #333;
}

:deep(.el-tabs__item.is-active) {
  color: #00ff88;
  border-bottom-color: #00ff88;
}

:deep(.el-tabs__item:hover) {
  color: #00ff88;
}

:deep(.el-tabs__content) {
  background-color: #1a1a1a;
}

:deep(.el-tab-pane) {
  background-color: #1a1a1a;
}

:deep(.el-input__wrapper) {
  background-color: #333;
  border: 1px solid #444;
}

:deep(.el-input__inner) {
  color: #fff;
  background-color: transparent;
}

:deep(.el-input__wrapper:hover) {
  border-color: #00ff88;
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #00ff88;
  box-shadow: 0 0 0 1px #00ff88;
}

:deep(.el-select .el-input__wrapper) {
  background-color: #333;
}

:deep(.el-button--primary) {
  background-color: #00ff88;
  border-color: #00ff88;
  color: #000;
}

:deep(.el-button--primary:hover) {
  background-color: #00cc6a;
  border-color: #00cc6a;
}

:deep(.el-button) {
  background-color: #333;
  border-color: #444;
  color: #fff;
}

:deep(.el-button:hover) {
  background-color: #444;
  border-color: #555;
}

:deep(.el-button--text) {
  color: #00ff88;
}

:deep(.el-button--text:hover) {
  color: #00cc6a;
  background-color: rgba(0, 255, 136, 0.1);
}

:deep(.el-tag) {
  background-color: #333;
  border-color: #444;
  color: #fff;
}

:deep(.el-tag--primary) {
  background-color: rgba(0, 255, 136, 0.2);
  border-color: #00ff88;
  color: #00ff88;
}

:deep(.el-tag--success) {
  background-color: rgba(103, 194, 58, 0.2);
  border-color: #67c23a;
  color: #67c23a;
}

:deep(.el-tag--warning) {
  background-color: rgba(230, 162, 60, 0.2);
  border-color: #e6a23c;
  color: #e6a23c;
}

:deep(.el-tag--danger) {
  background-color: rgba(245, 108, 108, 0.2);
  border-color: #f56c6c;
  color: #f56c6c;
}

:deep(.el-tag--info) {
  background-color: rgba(144, 147, 153, 0.2);
  border-color: #909399;
  color: #909399;
}

:deep(.el-switch__core) {
  background-color: #666;
}

:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #00ff88;
}

:deep(.el-avatar) {
  background-color: #333;
  color: #888;
}

:deep(.el-dialog) {
  background-color: #1a1a1a;
  border: 1px solid #333;
}

:deep(.el-dialog__header) {
  background-color: #2a2a2a;
  border-bottom: 1px solid #333;
  padding: 20px;
}

:deep(.el-dialog__title) {
  color: #fff;
}

:deep(.el-dialog__body) {
  background-color: #1a1a1a;
  color: #fff;
}

:deep(.el-form-item__label) {
  color: #fff;
}

:deep(.el-textarea__inner) {
  background-color: #333;
  border-color: #444;
  color: #fff;
}

:deep(.el-textarea__inner:hover) {
  border-color: #00ff88;
}

:deep(.el-textarea__inner:focus) {
  border-color: #00ff88;
  box-shadow: 0 0 0 1px #00ff88;
}

:deep(.el-checkbox) {
  color: #fff;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #00ff88;
  border-color: #00ff88;
}

:deep(.el-checkbox__inner:hover) {
  border-color: #00ff88;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__wrapper) {
  background-color: #333;
}

:deep(.el-descriptions) {
  background-color: #1a1a1a;
}

:deep(.el-descriptions__header) {
  background-color: #2a2a2a;
}

:deep(.el-descriptions__body) {
  background-color: #1a1a1a;
}

:deep(.el-descriptions-item__label) {
  background-color: #2a2a2a;
  color: #fff;
}

:deep(.el-descriptions-item__content) {
  background-color: #1a1a1a;
  color: #fff;
}
</style>