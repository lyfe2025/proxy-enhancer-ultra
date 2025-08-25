<template>
  <div class="not-found">
    <!-- 背景装饰形状 -->
    <div class="background-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
      <div class="shape shape-4"></div>
      <div class="shape shape-5"></div>
    </div>

    <div class="error-container">
      <!-- 404动画和图标 -->
      <ErrorAnimation />

      <!-- 错误信息 -->
      <ErrorContent
        :title="errorTitle"
        :description="errorDescription"
        :reasons="errorReasons"
      />

      <!-- 操作按钮 -->
      <ErrorActions
        @go-home="handleGoHome"
        @go-back="handleGoBack"
        @refresh="handleRefresh"
      />

      <!-- 搜索建议 -->
      <SearchSuggestion @search="handleSearch" />

      <!-- 快速链接 -->
      <QuickLinks :links="quickLinks" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

// 导入子组件
import ErrorAnimation from './components/ErrorAnimation.vue'
import ErrorContent from './components/ErrorContent.vue'
import ErrorActions from './components/ErrorActions.vue'
import SearchSuggestion from './components/SearchSuggestion.vue'
import QuickLinks from './components/QuickLinks.vue'

// 导入图标
import {
  House,
  User,
  Setting,
  Document,
  DataAnalysis,
  Connection
} from '@element-plus/icons-vue'

const route = useRoute()

// 错误信息配置
const errorTitle = ref('页面未找到')
const errorDescription = ref(`
  抱歉，您访问的页面不存在或已被移除。
  请检查URL是否正确，或返回首页继续浏览。
`)
const errorReasons = ref([
  '页面链接已过期或被删除',
  'URL地址输入错误',
  '您没有访问该页面的权限',
  '服务器配置问题'
])

// 快速链接配置
const quickLinks = ref([
  { path: '/', title: '首页', icon: House },
  { path: '/dashboard', title: '仪表盘', icon: DataAnalysis },
  { path: '/proxy', title: '代理管理', icon: Connection },
  { path: '/rules', title: '规则配置', icon: Setting },
  { path: '/data', title: '数据收集', icon: Document },
  { path: '/profile', title: '个人中心', icon: User }
])

// 事件处理
const handleGoHome = () => {
  ElMessage.success('正在跳转到首页...')
}

const handleGoBack = () => {
  ElMessage.info('正在返回上一页...')
}

const handleRefresh = () => {
  ElMessage.info('正在刷新页面...')
}

const handleSearch = (keyword: string) => {
  ElMessage.success(`正在搜索: ${keyword}`)
  console.log('搜索关键词:', keyword)
}

// 页面加载时记录404事件
const logNotFoundEvent = () => {
  console.warn('404 Page Not Found:', {
    path: route.path,
    fullPath: route.fullPath,
    query: route.query,
    timestamp: new Date().toISOString(),
    userAgent: navigator.userAgent,
    referrer: document.referrer
  })
  
  // 这里可以发送到分析服务
  // analytics.track('page_not_found', { path: route.path })
}

// 初始化
logNotFoundEvent()
</script>

<style scoped>
.not-found {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  overflow: hidden;
  text-align: center;
  padding: 20px;
}

.error-container {
  position: relative;
  z-index: 2;
  max-width: 800px;
  width: 100%;
}

/* 背景装饰形状 */
.background-shapes {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
}

.shape {
  position: absolute;
  background: linear-gradient(45deg, rgba(0, 255, 136, 0.1), rgba(0, 136, 255, 0.1));
  border-radius: 50%;
  animation: float 6s ease-in-out infinite;
}

.shape-1 {
  width: 80px;
  height: 80px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 120px;
  height: 120px;
  top: 20%;
  right: 15%;
  animation-delay: 1s;
}

.shape-3 {
  width: 60px;
  height: 60px;
  bottom: 30%;
  left: 20%;
  animation-delay: 2s;
}

.shape-4 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  right: 10%;
  animation-delay: 3s;
}

.shape-5 {
  width: 40px;
  height: 40px;
  top: 50%;
  left: 50%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.3;
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
    opacity: 0.6;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .error-container {
    padding: 20px 16px;
  }
}

/* 深色主题适配 */
.dark .not-found {
  background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
}
</style>