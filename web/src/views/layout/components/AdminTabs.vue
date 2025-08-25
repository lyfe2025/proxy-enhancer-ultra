<template>
  <div class="tabs-container">
    <el-scrollbar class="tabs-scrollbar">
      <div class="tabs-content">
        <div
          v-for="tab in visitedViews"
          :key="tab.path"
          class="tab-item"
          :class="{ active: isActive(tab) }"
          @click="goToTab(tab)"
          @contextmenu.prevent="openContextMenu($event, tab)"
        >
          <el-icon v-if="tab.meta?.icon" class="tab-icon">
            <component :is="tab.meta.icon" />
          </el-icon>
          <span class="tab-title">{{ getTabTitle(tab) }}</span>
          <el-icon
            v-if="!isAffix(tab)"
            class="tab-close"
            @click.stop="closeTab(tab)"
          >
            <Close />
          </el-icon>
        </div>
      </div>
    </el-scrollbar>

    <!-- 右键菜单 -->
    <ul
      v-show="contextMenuVisible"
      :style="{ left: contextMenuLeft + 'px', top: contextMenuTop + 'px' }"
      class="context-menu"
      @click="closeContextMenu"
    >
      <li @click="refreshTab">刷新</li>
      <li @click="closeCurrentTab">关闭</li>
      <li @click="closeOtherTabs">关闭其他</li>
      <li @click="closeAllTabs">关闭全部</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Close } from '@element-plus/icons-vue'

interface TabView {
  name?: string
  path: string
  fullPath: string
  meta?: {
    title?: string
    icon?: string
    affix?: boolean
  }
  query?: Record<string, any>
  params?: Record<string, any>
}

const props = defineProps<{
  visitedViews: TabView[]
  cachedViews: string[]
}>()

const emit = defineEmits<{
  addView: [view: TabView]
  deleteView: [view: TabView]
  deleteOtherViews: [view: TabView]
  deleteAllViews: []
  refreshView: [view: TabView]
}>()

const route = useRoute()
const router = useRouter()

// 右键菜单相关
const contextMenuVisible = ref(false)
const contextMenuLeft = ref(0)
const contextMenuTop = ref(0)
const selectedTab = ref<TabView | null>(null)

onMounted(() => {
  initTags()
  document.addEventListener('click', closeContextMenu)
})

onUnmounted(() => {
  document.removeEventListener('click', closeContextMenu)
})

const isActive = (tab: TabView) => {
  return tab.path === route.path
}

const isAffix = (tab: TabView) => {
  return tab.meta?.affix
}

const getTabTitle = (tab: TabView) => {
  return tab.meta?.title || tab.name || '未命名'
}

const goToTab = (tab: TabView) => {
  if (tab.path !== route.path) {
    router.push({
      path: tab.path,
      query: tab.query
    })
  }
}

const closeTab = (tab: TabView) => {
  if (isAffix(tab)) return
  
  const isCurrentTab = isActive(tab)
  emit('deleteView', tab)
  
  if (isCurrentTab) {
    const index = props.visitedViews.findIndex(v => v.path === tab.path)
    if (index > 0) {
      const prevTab = props.visitedViews[index - 1]
      goToTab(prevTab)
    } else if (props.visitedViews.length > 1) {
      const nextTab = props.visitedViews[1]
      goToTab(nextTab)
    } else {
      router.push('/')
    }
  }
}

const openContextMenu = (event: MouseEvent, tab: TabView) => {
  selectedTab.value = tab
  contextMenuLeft.value = event.clientX
  contextMenuTop.value = event.clientY
  contextMenuVisible.value = true
}

const closeContextMenu = () => {
  contextMenuVisible.value = false
  selectedTab.value = null
}

const refreshTab = () => {
  if (selectedTab.value) {
    emit('refreshView', selectedTab.value)
    if (isActive(selectedTab.value)) {
      // 刷新当前页面
      nextTick(() => {
        router.replace({
          path: '/redirect' + selectedTab.value!.fullPath
        })
      })
    }
  }
}

const closeCurrentTab = () => {
  if (selectedTab.value) {
    closeTab(selectedTab.value)
  }
}

const closeOtherTabs = () => {
  if (selectedTab.value) {
    emit('deleteOtherViews', selectedTab.value)
    if (!isActive(selectedTab.value)) {
      goToTab(selectedTab.value)
    }
  }
}

const closeAllTabs = () => {
  emit('deleteAllViews')
  router.push('/')
}

const initTags = () => {
  const affixTabs = props.visitedViews.filter(tab => tab.meta?.affix)
  for (const tab of affixTabs) {
    if (tab.name) {
      emit('addView', tab)
    }
  }
}
</script>

<style scoped>
.tabs-container {
  height: 40px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  position: relative;
}

.tabs-scrollbar {
  height: 100%;
}

.tabs-content {
  display: flex;
  height: 100%;
  align-items: center;
  padding: 0 16px;
  gap: 4px;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  white-space: nowrap;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  background: var(--el-fill-color-blank);
}

.tab-item:hover {
  background: var(--el-fill-color-light);
}

.tab-item.active {
  background: var(--el-color-primary-light-9);
  border-color: var(--el-color-primary-light-7);
  color: var(--el-color-primary);
}

.tab-icon {
  font-size: 12px;
}

.tab-title {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-close {
  font-size: 12px;
  margin-left: 4px;
  opacity: 0.6;
  transition: opacity 0.3s ease;
  border-radius: 50%;
  padding: 2px;
}

.tab-close:hover {
  opacity: 1;
  background: var(--el-color-danger-light-8);
  color: var(--el-color-danger);
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  z-index: 9999;
  background: var(--el-bg-color-overlay);
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  box-shadow: var(--el-box-shadow);
  padding: 4px 0;
  margin: 0;
  list-style: none;
  min-width: 120px;
}

.context-menu li {
  padding: 8px 16px;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.3s ease;
}

.context-menu li:hover {
  background: var(--el-fill-color-light);
}

/* 深色主题适配 */
.dark .tabs-container {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
