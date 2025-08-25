<template>
  <div class="stats-grid">
    <div class="stat-card">
      <div class="stat-icon proxy">
        <el-icon><Connection /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-value">{{ stats.totalProxies }}</div>
        <div class="stat-label">代理总数</div>
        <div class="stat-change positive">
          <el-icon><ArrowUp /></el-icon>
          +{{ stats.proxyGrowth }}%
        </div>
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-icon rules">
        <el-icon><SetUp /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-value">{{ stats.activeRules }}</div>
        <div class="stat-label">活跃规则</div>
        <div class="stat-change positive">
          <el-icon><ArrowUp /></el-icon>
          +{{ stats.rulesGrowth }}%
        </div>
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-icon traffic">
        <el-icon><TrendCharts /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-value">{{ formatBytes(stats.totalTraffic) }}</div>
        <div class="stat-label">总流量</div>
        <div class="stat-change positive">
          <el-icon><ArrowUp /></el-icon>
          +{{ stats.trafficGrowth }}%
        </div>
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-icon users">
        <el-icon><User /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-value">{{ stats.onlineUsers }}</div>
        <div class="stat-label">在线用户</div>
        <div class="stat-change neutral">
          <el-icon><Minus /></el-icon>
          {{ stats.userChange }}%
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Connection, SetUp, TrendCharts, User, ArrowUp, Minus } from '@element-plus/icons-vue'

interface DashboardStatsData {
  totalProxies: number
  proxyGrowth: number
  activeRules: number
  rulesGrowth: number
  totalTraffic: number
  trafficGrowth: number
  onlineUsers: number
  userChange: number
}

defineProps<{
  stats: DashboardStatsData
}>()

// 格式化字节数
const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  border: 1px solid var(--el-border-color-lighter);
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  flex-shrink: 0;
}

.stat-icon.proxy {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.rules {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.traffic {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.users {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  line-height: 1;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: var(--el-text-color-regular);
  margin-bottom: 8px;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
}

.stat-change.positive {
  color: var(--el-color-success);
}

.stat-change.negative {
  color: var(--el-color-danger);
}

.stat-change.neutral {
  color: var(--el-color-warning);
}

/* 深色主题适配 */
.dark .stat-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

.dark .stat-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .stat-card {
    padding: 16px;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }
  
  .stat-value {
    font-size: 24px;
  }
}
</style>
