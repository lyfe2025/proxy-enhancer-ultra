<template>
  <div class="header-toolbar">
    <!-- 全屏切换 -->
    <el-tooltip :content="isFullscreen ? '退出全屏' : '进入全屏'" placement="bottom">
      <el-button text circle @click="toggleFullscreen">
        <el-icon size="18"><FullScreen /></el-icon>
      </el-button>
    </el-tooltip>

    <!-- 主题切换 -->
    <el-tooltip :content="isDark ? '切换到浅色模式' : '切换到深色模式'" placement="bottom">
      <el-button text circle @click="toggleTheme">
        <el-icon size="18">
          <Sunny v-if="isDark" />
          <Moon v-else />
        </el-icon>
      </el-button>
    </el-tooltip>
  </div>
</template>

<script setup lang="ts">
import {
  FullScreen,
  Sunny,
  Moon
} from '@element-plus/icons-vue'

const props = defineProps<{
  isFullscreen: boolean
  isDark: boolean
}>()

const emit = defineEmits<{
  toggleFullscreen: []
  toggleTheme: []
}>()

const toggleFullscreen = () => {
  emit('toggleFullscreen')
}

const toggleTheme = () => {
  emit('toggleTheme')
}
</script>

<style scoped>
.header-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
