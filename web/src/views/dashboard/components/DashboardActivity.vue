<template>
  <el-card class="activity-card">
    <template #header>
      <div class="card-header">
        <span class="card-title">
          <el-icon><Clock /></el-icon>
          最近活动
        </span>
        <el-button text size="small" @click="refreshActivity">
          刷新
        </el-button>
      </div>
    </template>
    <div class="activity-list">
      <div
        v-for="activity in activities"
        :key="activity.id"
        class="activity-item"
      >
        <div :class="['activity-icon', activity.type]">
          <el-icon v-if="activity.type === 'success'"><SuccessFilled /></el-icon>
          <el-icon v-else-if="activity.type === 'warning'"><WarningFilled /></el-icon>
          <el-icon v-else-if="activity.type === 'error'"><CircleCloseFilled /></el-icon>
          <el-icon v-else><InfoFilled /></el-icon>
        </div>
        <div class="activity-content">
          <div class="activity-title">{{ activity.title }}</div>
          <div class="activity-desc">{{ activity.description }}</div>
          <div class="activity-time">{{ formatTime(activity.time) }}</div>
        </div>
      </div>
      <div v-if="!activities.length" class="empty-state">
        <el-empty description="暂无活动记录" />
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Clock, SuccessFilled, WarningFilled, CircleCloseFilled, InfoFilled } from '@element-plus/icons-vue'

interface ActivityItem {
  id: string
  type: 'success' | 'warning' | 'error' | 'info'
  title: string
  description: string
  time: string
}

defineProps<{
  activities: ActivityItem[]
}>()

const emit = defineEmits<{
  refresh: []
}>()

const refreshActivity = () => {
  emit('refresh')
}

const formatTime = (timeString: string) => {
  const now = new Date()
  const time = new Date(timeString)
  const diff = now.getTime() - time.getTime()
  
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return time.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.activity-card {
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.activity-list {
  max-height: 400px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: white;
  flex-shrink: 0;
}

.activity-icon.success {
  background: var(--el-color-success);
}

.activity-icon.warning {
  background: var(--el-color-warning);
}

.activity-icon.error {
  background: var(--el-color-danger);
}

.activity-icon.info {
  background: var(--el-color-info);
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-weight: 500;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
}

.activity-desc {
  font-size: 13px;
  color: var(--el-text-color-regular);
  margin-bottom: 4px;
}

.activity-time {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}

.empty-state {
  padding: 40px 20px;
  text-align: center;
}

/* 深色主题适配 */
.dark .activity-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
