<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <AdminSidebar
      :sidebar-collapsed="themeStore.sidebarCollapsed"
      :menu-routes="menuRoutes"
      @toggle-sidebar="toggleSidebar"
    />

    <!-- 主内容区域 -->
    <el-container class="main-container">
      <!-- 顶部导航栏 -->
      <AdminHeader
        :breadcrumb-list="breadcrumbList"
        :notifications="notifications"
        :user-info="userInfo"
        :is-fullscreen="isFullscreen"
        :is-dark="themeStore.isDark"
        @toggle-fullscreen="toggleFullscreen"
        @toggle-theme="toggleTheme"
        @mark-as-read="markNotificationAsRead"
        @mark-all-as-read="markAllNotificationsAsRead"
        @logout="handleLogout"
      />

      <!-- 标签页 -->
      <AdminTabs
        :visited-views="visitedViews"
        :cached-views="cachedViews"
        @add-view="addVisitedView"
        @delete-view="deleteVisitedView"
        @delete-other-views="deleteOtherViews"
        @delete-all-views="deleteAllViews"
        @refresh-view="refreshView"
      />

      <!-- 主要内容区域 -->
      <el-main class="main-content">
        <router-view v-slot="{ Component, route }">
          <transition name="fade-transform" mode="out-in">
            <keep-alive :include="cachedViews">
              <component :is="Component" :key="route.path" />
            </keep-alive>
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { menuRoutes } from '@/router'

// 导入子组件
import AdminSidebar from './components/AdminSidebar.vue'
import AdminHeader from './components/AdminHeader.vue'
import AdminTabs from './components/AdminTabs.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()

// 全屏状态
const isFullscreen = ref(false)

// 用户信息
const userInfo = computed(() => ({
  username: authStore.userInfo?.username || '用户',
  avatar: authStore.userInfo?.avatar
}))

// 通知数据
const notifications = ref([
  {
    id: '1',
    title: '系统更新',
    content: '系统已更新到最新版本 v2.1.0',
    type: 'info' as const,
    read: false,
    createdAt: new Date().toISOString()
  },
  {
    id: '2',
    title: '代理异常',
    content: '代理服务器 proxy-01 连接异常，请及时处理',
    type: 'warning' as const,
    read: false,
    createdAt: new Date(Date.now() - 1000 * 60 * 30).toISOString()
  }
])

// 标签页相关
const visitedViews = ref<any[]>([])
const cachedViews = ref<string[]>([])

// 面包屑导航
const breadcrumbList = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  const first = matched[0]

  if (!first || (first.name !== 'Dashboard' && first.path !== '/')) {
    matched.unshift({ path: '/', meta: { title: '首页', icon: 'House' } } as any)
  }

  return matched.map(item => ({
    path: item.path,
    title: item.meta?.title || '',
    icon: item.meta?.icon
  }))
})

onMounted(() => {
  initTags()
  addVisitedView(route as any)
})

// 监听路由变化
watch(route, (newRoute) => {
  addVisitedView(newRoute as any)
})

// 初始化固定标签
const initTags = () => {
  const affixTags = [
    { path: '/', meta: { title: '首页', affix: true, icon: 'House' } }
  ]
  for (const tag of affixTags) {
    if (tag.meta?.affix) {
      addVisitedView(tag as any)
    }
  }
}

// 标签页管理方法
const addVisitedView = (view: any) => {
  if (visitedViews.value.some(v => v.path === view.path)) return
  
  visitedViews.value.push({
    name: view.name,
    path: view.path,
    fullPath: view.fullPath,
    meta: view.meta,
    query: view.query,
    params: view.params
  })

  if (view.name && !view.meta?.noCache) {
    cachedViews.value.push(view.name)
  }
}

const deleteVisitedView = (view: any) => {
  const index = visitedViews.value.findIndex(v => v.path === view.path)
  if (index > -1) {
    visitedViews.value.splice(index, 1)
  }
  
  if (view.name) {
    const cacheIndex = cachedViews.value.indexOf(view.name)
    if (cacheIndex > -1) {
      cachedViews.value.splice(cacheIndex, 1)
    }
  }
}

const deleteOtherViews = (view: any) => {
  visitedViews.value = visitedViews.value.filter(v => v.meta?.affix || v.path === view.path)
  cachedViews.value = cachedViews.value.filter(name => {
    const tab = visitedViews.value.find(v => v.name === name)
    return tab?.meta?.affix || tab?.path === view.path
  })
}

const deleteAllViews = () => {
  visitedViews.value = visitedViews.value.filter(v => v.meta?.affix)
  cachedViews.value = cachedViews.value.filter(name => {
    const tab = visitedViews.value.find(v => v.name === name)
    return tab?.meta?.affix
  })
}

const refreshView = (view: any) => {
  if (view.name) {
    const index = cachedViews.value.indexOf(view.name)
    if (index > -1) {
      cachedViews.value.splice(index, 1)
      setTimeout(() => {
        cachedViews.value.push(view.name)
      }, 100)
    }
  }
}

// 侧边栏切换
const toggleSidebar = () => {
  themeStore.toggleSidebar()
}

// 全屏切换
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

// 主题切换
const toggleTheme = () => {
  themeStore.toggleDark()
}

// 通知管理
const markNotificationAsRead = (id: string) => {
  const notification = notifications.value.find(n => n.id === id)
  if (notification) {
    notification.read = true
  }
}

const markAllNotificationsAsRead = () => {
  notifications.value.forEach(notification => {
    notification.read = true
  })
}

// 退出登录
const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
    ElMessage.success('已退出登录')
  } catch (error) {
    console.error('退出登录失败:', error)
    ElMessage.error('退出登录失败')
  }
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  display: flex;
  background: var(--el-bg-color-page);
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 24px;
  overflow: auto;
  background: var(--el-bg-color-page);
}

/* 页面切换动画 */
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

/* 深色主题适配 */
.dark .admin-layout {
  background: var(--el-bg-color-page);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-content {
    padding: 16px;
  }
}
</style>