<template>
  <template v-if="!route.children || route.children.length === 0">
    <!-- 单级菜单项 -->
    <el-menu-item
      :index="resolvePath(route.path)"
      :class="{ 'is-active': isActive }"
    >
      <el-icon v-if="route.meta?.icon">
        <component :is="route.meta.icon" />
      </el-icon>
      <template #title>
        <span>{{ route.meta?.title }}</span>
      </template>
    </el-menu-item>
  </template>

  <template v-else>
    <!-- 多级菜单项 -->
    <el-sub-menu
      :index="resolvePath(route.path)"
      :class="{ 'is-active': hasActiveChild }"
    >
      <template #title>
        <el-icon v-if="route.meta?.icon">
          <component :is="route.meta.icon" />
        </el-icon>
        <span>{{ route.meta?.title }}</span>
      </template>
      
      <template v-for="child in route.children" :key="child.path">
        <sidebar-menu-item
          v-if="!child.meta?.hidden"
          :route="child"
          :base-path="resolvePath(child.path)"
        />
      </template>
    </el-sub-menu>
  </template>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

interface Props {
  route: RouteRecordRaw
  basePath: string
}

const props = defineProps<Props>()
const currentRoute = useRoute()

// 解析完整路径
const resolvePath = (path: string) => {
  if (path.startsWith('/')) {
    return path
  }
  return `${props.basePath}/${path}`.replace(/\/+/g, '/')
}

// 检查当前菜单项是否激活
const isActive = computed(() => {
  const resolvedPath = resolvePath(props.route.path)
  return currentRoute.path === resolvedPath
})

// 检查是否有激活的子菜单
const hasActiveChild = computed(() => {
  if (!props.route.children) return false
  
  return props.route.children.some(child => {
    const childPath = resolvePath(child.path)
    return currentRoute.path.startsWith(childPath)
  })
})
</script>

<style scoped>
.el-menu-item,
.el-sub-menu {
  border-radius: 6px;
  margin: 2px 8px;
}

.el-menu-item.is-active {
  background: var(--el-color-primary-light-9) !important;
  color: var(--el-color-primary) !important;
}

.el-sub-menu.is-active > .el-sub-menu__title {
  color: var(--el-color-primary) !important;
}

:deep(.el-menu-item) {
  height: 48px;
  line-height: 48px;
}

:deep(.el-sub-menu__title) {
  height: 48px;
  line-height: 48px;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background: var(--el-fill-color-light) !important;
}

:deep(.el-menu-item.is-active) {
  background: var(--el-color-primary-light-9) !important;
  color: var(--el-color-primary) !important;
}

:deep(.el-menu-item.is-active::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--el-color-primary);
  border-radius: 0 2px 2px 0;
}
</style>