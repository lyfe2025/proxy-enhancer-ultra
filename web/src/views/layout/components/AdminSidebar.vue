<template>
  <el-aside
    :width="sidebarCollapsed ? '64px' : '240px'"
    class="sidebar"
    :class="{ 'sidebar-collapsed': sidebarCollapsed }"
  >
    <!-- Logo区域 -->
    <div class="sidebar-header">
      <div class="logo">
        <div class="logo-icon">
          <el-icon size="24"><Connection /></el-icon>
        </div>
        <transition name="fade">
          <span v-show="!sidebarCollapsed" class="logo-text">
            代理增强器
          </span>
        </transition>
      </div>
    </div>

    <!-- 导航菜单 -->
    <el-scrollbar class="sidebar-menu-container">
      <el-menu
        :default-active="activeMenu"
        :collapse="sidebarCollapsed"
        :unique-opened="true"
        router
        class="sidebar-menu"
      >
        <template v-for="route in menuRoutes" :key="route.path">
          <sidebar-menu-item
            v-if="!route.meta?.hidden"
            :route="route"
            :base-path="route.path"
          />
        </template>
      </el-menu>
    </el-scrollbar>

    <!-- 侧边栏底部 -->
    <div class="sidebar-footer">
      <el-tooltip
        :content="sidebarCollapsed ? '展开侧边栏' : '收起侧边栏'"
        placement="right"
      >
        <el-button
          text
          class="collapse-btn"
          @click="toggleSidebar"
        >
          <el-icon>
            <Expand v-if="sidebarCollapsed" />
            <Fold v-else />
          </el-icon>
        </el-button>
      </el-tooltip>
    </div>
  </el-aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, type RouteRecordRaw } from 'vue-router'
import { Connection, Expand, Fold } from '@element-plus/icons-vue'
import SidebarMenuItem from './SidebarMenuItem.vue'

const props = defineProps<{
  sidebarCollapsed: boolean
  menuRoutes: RouteRecordRaw[]
}>()

const emit = defineEmits<{
  toggleSidebar: []
}>()

const route = useRoute()

// 计算当前激活的菜单项
const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

const toggleSidebar = () => {
  emit('toggleSidebar')
}
</script>

<style scoped>
.sidebar {
  background: var(--el-bg-color);
  border-right: 1px solid var(--el-border-color-light);
  transition: all 0.3s ease;
  position: relative;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid var(--el-border-color-lighter);
  padding: 0 16px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.logo-icon {
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 16px;
  white-space: nowrap;
}

.sidebar-collapsed .logo-text {
  display: none;
}

.sidebar-menu-container {
  flex: 1;
  overflow: hidden;
}

.sidebar-menu {
  border: none;
  background: transparent;
  width: 100%;
}

.sidebar-footer {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-top: 1px solid var(--el-border-color-lighter);
}

.collapse-btn {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.collapse-btn:hover {
  background: var(--el-fill-color-light);
}

/* 动画效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 深色主题适配 */
.dark .sidebar {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
