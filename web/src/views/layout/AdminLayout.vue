<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <el-aside :width="sidebarWidth" class="sidebar">
      <div class="sidebar-header">
        <div class="logo">
          <span class="logo-icon">🚀</span>
          <span v-show="!appStore.sidebarCollapsed" class="logo-text">代理增强器</span>
        </div>
      </div>
      
      <el-scrollbar class="sidebar-menu">
        <el-menu
          :default-active="activeMenu"
          :collapse="appStore.sidebarCollapsed"
          :unique-opened="true"
          router
          class="menu"
        >
          <template v-for="item in menuItems" :key="item.path">
            <el-menu-item
              v-if="!item.children"
              :index="item.path"
              :disabled="!hasPermission(item.permissions)"
            >
              <el-icon><component :is="item.icon" /></el-icon>
              <template #title>{{ item.title }}</template>
            </el-menu-item>
            
            <el-sub-menu
              v-else
              :index="item.path"
              :disabled="!hasPermission(item.permissions)"
            >
              <template #title>
                <el-icon><component :is="item.icon" /></el-icon>
                <span>{{ item.title }}</span>
              </template>
              
              <el-menu-item
                v-for="child in item.children"
                :key="child.path"
                :index="child.path"
                :disabled="!hasPermission(child.permissions)"
              >
                <el-icon><component :is="child.icon" /></el-icon>
                <template #title>{{ child.title }}</template>
              </el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
      </el-scrollbar>
    </el-aside>
    
    <!-- 主内容区 -->
    <el-container class="main-container">
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <el-button
            type="text"
            class="collapse-btn"
            @click="toggleSidebar"
          >
            <el-icon><Expand v-if="appStore.sidebarCollapsed" /><Fold v-else /></el-icon>
          </el-button>
          
          <el-breadcrumb class="breadcrumb" separator="/">
            <el-breadcrumb-item
              v-for="item in appStore.breadcrumbs"
              :key="item.path"
              :to="item.path"
            >
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- 主题切换 -->
          <el-tooltip content="切换主题" placement="bottom">
            <el-button
              type="text"
              class="theme-btn"
              @click="toggleTheme"
            >
              <el-icon><Sunny v-if="appStore.theme === 'dark'" /><Moon v-else /></el-icon>
            </el-button>
          </el-tooltip>
          
          <!-- 全屏切换 -->
          <el-tooltip content="全屏" placement="bottom">
            <el-button
              type="text"
              class="fullscreen-btn"
              @click="toggleFullscreen"
            >
              <el-icon><FullScreen /></el-icon>
            </el-button>
          </el-tooltip>
          
          <!-- 用户菜单 -->
          <el-dropdown class="user-dropdown" @command="handleUserCommand">
            <div class="user-info">
              <el-avatar :size="32" class="user-avatar">
                {{ authStore.user?.username?.charAt(0).toUpperCase() }}
              </el-avatar>
              <span v-show="!appStore.isMobile" class="username">
                {{ authStore.user?.username }}
              </span>
              <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
            </div>
            
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人资料
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>
                  账户设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- 主内容 -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Dashboard,
  Connection,
  Setting,
  DataAnalysis,
  Tools,
  User,
  Expand,
  Fold,
  Sunny,
  Moon,
  FullScreen,
  ArrowDown,
  SwitchButton
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import { useAppStore } from '@/store/app'
import type { MenuItem } from '@/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

// 计算属性
const sidebarWidth = computed(() => {
  return appStore.sidebarCollapsed ? '64px' : '240px'
})

const activeMenu = computed(() => {
  return route.path
})

// 菜单项配置
const menuItems: MenuItem[] = [
  {
    path: '/dashboard',
    title: '仪表盘',
    icon: 'Dashboard',
    permissions: []
  },
  {
    path: '/proxy',
    title: '代理管理',
    icon: 'Connection',
    permissions: ['proxy:read']
  },
  {
    path: '/rules',
    title: '规则配置',
    icon: 'Setting',
    permissions: ['rule:read']
  },
  {
    path: '/data-collection',
    title: '数据收集',
    icon: 'DataAnalysis',
    permissions: ['popup:read'],
    children: [
      {
        path: '/data-collection',
        title: '弹窗管理',
        icon: 'DataAnalysis',
        permissions: ['popup:read']
      },
      {
        path: '/data-collection/submissions',
        title: '提交数据',
        icon: 'DataAnalysis',
        permissions: ['submission:read']
      }
    ]
  },
  {
    path: '/system',
    title: '系统设置',
    icon: 'Tools',
    permissions: ['system:read'],
    children: [
      {
        path: '/system/users',
        title: '用户管理',
        icon: 'User',
        permissions: ['user:read']
      },
      {
        path: '/system/roles',
        title: '角色管理',
        icon: 'Setting',
        permissions: ['role:read']
      },
      {
        path: '/system/permissions',
        title: '权限管理',
        icon: 'Setting',
        permissions: ['permission:read']
      },
      {
        path: '/system/logs',
        title: '系统日志',
        icon: 'DataAnalysis',
        permissions: ['log:read']
      }
    ]
  }
]

// 权限检查
const hasPermission = (permissions?: string[]) => {
  if (!permissions || permissions.length === 0) return true
  return permissions.some(permission => authStore.hasPermission(permission))
}

// 切换侧边栏
const toggleSidebar = () => {
  appStore.toggleSidebar()
}

// 切换主题
const toggleTheme = () => {
  appStore.setTheme(appStore.theme === 'dark' ? 'light' : 'dark')
}

// 全屏切换
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

// 处理用户菜单命令
const handleUserCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      // 跳转到个人资料页面
      ElMessage.info('个人资料功能开发中')
      break
    case 'settings':
      // 跳转到账户设置页面
      ElMessage.info('账户设置功能开发中')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm(
          '确定要退出登录吗？',
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await authStore.logout()
        ElMessage.success('已退出登录')
        router.push('/login')
      } catch {
        // 用户取消
      }
      break
  }
}

// 更新面包屑
const updateBreadcrumbs = () => {
  const breadcrumbs = []
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  
  for (const item of matched) {
    breadcrumbs.push({
      path: item.path,
      title: item.meta.title as string
    })
  }
  
  appStore.setBreadcrumbs(breadcrumbs)
}

// 监听路由变化
watch(
  () => route.path,
  () => {
    updateBreadcrumbs()
  },
  { immediate: true }
)

// 组件挂载
onMounted(() => {
  // 初始化应用状态
  appStore.initApp()
  
  // 检测设备类型
  appStore.detectDevice()
  
  // 监听窗口大小变化
  window.addEventListener('resize', appStore.detectDevice)
})
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  display: flex;
  background-color: #0a0a0a;
}

/* 侧边栏 */
.sidebar {
  background-color: #1a1a1a;
  border-right: 1px solid #333;
  transition: width 0.3s ease;
  overflow: hidden;
}

.sidebar-header {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  border-bottom: 1px solid #333;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #00ff88;
  font-weight: 600;
  font-size: 18px;
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  white-space: nowrap;
  transition: opacity 0.3s ease;
}

.sidebar-menu {
  height: calc(100vh - 60px);
}

.menu {
  border: none;
  background-color: transparent;
}

.menu :deep(.el-menu-item) {
  color: #ccc;
  border-radius: 8px;
  margin: 4px 12px;
  height: 48px;
  line-height: 48px;
}

.menu :deep(.el-menu-item:hover) {
  background-color: rgba(0, 255, 136, 0.1);
  color: #00ff88;
}

.menu :deep(.el-menu-item.is-active) {
  background-color: rgba(0, 255, 136, 0.2);
  color: #00ff88;
  border-right: none;
}

.menu :deep(.el-sub-menu__title) {
  color: #ccc;
  border-radius: 8px;
  margin: 4px 12px;
  height: 48px;
  line-height: 48px;
}

.menu :deep(.el-sub-menu__title:hover) {
  background-color: rgba(0, 255, 136, 0.1);
  color: #00ff88;
}

.menu :deep(.el-sub-menu.is-active .el-sub-menu__title) {
  color: #00ff88;
}

/* 主容器 */
.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 顶部导航 */
.header {
  background-color: #1a1a1a;
  border-bottom: 1px solid #333;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.collapse-btn {
  color: #ccc;
  font-size: 18px;
}

.collapse-btn:hover {
  color: #00ff88;
}

.breadcrumb {
  color: #ccc;
}

.breadcrumb :deep(.el-breadcrumb__item) {
  color: #ccc;
}

.breadcrumb :deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: #00ff88;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.theme-btn,
.fullscreen-btn {
  color: #ccc;
  font-size: 18px;
}

.theme-btn:hover,
.fullscreen-btn:hover {
  color: #00ff88;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  transition: background-color 0.3s ease;
}

.user-info:hover {
  background-color: rgba(0, 255, 136, 0.1);
}

.user-avatar {
  background-color: #00ff88;
  color: #000;
  font-weight: 600;
}

.username {
  color: #ccc;
  font-size: 14px;
}

.dropdown-icon {
  color: #ccc;
  font-size: 12px;
}

/* 主内容 */
.main-content {
  background-color: #0a0a0a;
  padding: 20px;
  overflow-y: auto;
}

/* 过渡动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s ease;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1000;
    height: 100vh;
  }
  
  .main-container {
    margin-left: 0;
  }
  
  .header-left .breadcrumb {
    display: none;
  }
  
  .username {
    display: none;
  }
}

/* 深色主题样式 */
:deep(.el-dropdown-menu) {
  background-color: #1a1a1a;
  border: 1px solid #333;
}

:deep(.el-dropdown-menu__item) {
  color: #ccc;
}

:deep(.el-dropdown-menu__item:hover) {
  background-color: rgba(0, 255, 136, 0.1);
  color: #00ff88;
}

:deep(.el-button--text) {
  background-color: transparent;
  border: none;
}

:deep(.el-breadcrumb__separator) {
  color: #666;
}
</style>