<template>
  <div id="app" :class="{ 'dark': themeStore.isDark }">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

// 初始化主题
onMounted(() => {
  themeStore.loadSettings()
  themeStore.applyTheme()
})
</script>

<style>
#app {
  width: 100%;
  height: 100vh;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: var(--el-bg-color-page);
  color: var(--el-text-color-primary);
  transition: all 0.3s ease;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
}
</style>