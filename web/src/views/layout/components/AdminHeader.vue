<template>
  <el-header class="header">
    <div class="header-left">
      <!-- 面包屑导航 -->
      <el-breadcrumb separator="/" class="breadcrumb">
        <el-breadcrumb-item
          v-for="item in breadcrumbList"
          :key="item.path"
          :to="item.path"
        >
          <el-icon v-if="item.icon" class="breadcrumb-icon">
            <component :is="item.icon" />
          </el-icon>
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="header-right">
      <!-- 工具栏 -->
      <HeaderToolbar
        :is-fullscreen="isFullscreen"
        :is-dark="isDark"
        @toggle-fullscreen="emit('toggleFullscreen')"
        @toggle-theme="emit('toggleTheme')"
      />

      <!-- 通知中心 -->
      <NotificationCenter
        :notifications="notifications"
        @mark-as-read="emit('markAsRead', $event)"
        @mark-all-as-read="emit('markAllAsRead')"
      />

      <!-- 用户菜单 -->
      <UserMenu
        :user-info="userInfo"
        @logout="emit('logout')"
      />
    </div>
  </el-header>
</template>

<script setup lang="ts">
// 导入子组件
import HeaderToolbar from './header/HeaderToolbar.vue'
import NotificationCenter from './header/NotificationCenter.vue'
import UserMenu from './header/UserMenu.vue'

interface BreadcrumbItem {
  path: string
  title: string
  icon?: string
}

interface NotificationItem {
  id: string
  title: string
  content: string
  type: 'info' | 'warning' | 'error' | 'success'
  read: boolean
  createdAt: string
}

interface UserInfo {
  username: string
  avatar?: string
}

const props = defineProps<{
  breadcrumbList: BreadcrumbItem[]
  notifications: NotificationItem[]
  userInfo: UserInfo
  isFullscreen: boolean
  isDark: boolean
}>()

const emit = defineEmits<{
  toggleFullscreen: []
  toggleTheme: []
  markAsRead: [id: string]
  markAllAsRead: []
  logout: []
}>()
</script>

<style scoped>
.header {
  height: 60px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.header-left {
  flex: 1;
}

.breadcrumb {
  font-size: 14px;
}

.breadcrumb-icon {
  margin-right: 4px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 深色主题适配 */
.dark .header {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    padding: 0 16px;
  }
}
</style>