<template>
  <el-card class="status-card">
    <template #header>
      <div class="card-header">
        <span class="card-title">
          <el-icon><Monitor /></el-icon>
          系统状态
        </span>
        <el-button text size="small" @click="refreshStatus">
          刷新
        </el-button>
      </div>
    </template>
    <div class="status-list">
      <div
        v-for="service in services"
        :key="service.name"
        class="status-item"
      >
        <div class="service-info">
          <div class="service-name">{{ service.name }}</div>
          <div class="service-desc">{{ service.description }}</div>
        </div>
        <div class="service-status">
          <el-tag
            :type="getStatusType(service.status)"
            size="small"
            effect="dark"
          >
            {{ getStatusText(service.status) }}
          </el-tag>
        </div>
      </div>
      <div v-if="!services.length" class="empty-state">
        <el-empty description="暂无服务信息" />
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Monitor } from '@element-plus/icons-vue'

interface ServiceItem {
  name: string
  description: string
  status: 'running' | 'stopped' | 'error' | 'warning'
}

defineProps<{
  services: ServiceItem[]
}>()

const emit = defineEmits<{
  refresh: []
}>()

const refreshStatus = () => {
  emit('refresh')
}

const getStatusType = (status: string): "success" | "info" | "warning" | "danger" => {
  const statusMap: Record<string, "success" | "info" | "warning" | "danger"> = {
    running: 'success',
    stopped: 'info',
    warning: 'warning',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    warning: '警告',
    error: '错误'
  }
  return statusMap[status] || status
}
</script>

<style scoped>
.status-card {
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

.status-list {
  max-height: 400px;
  overflow-y: auto;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.status-item:last-child {
  border-bottom: none;
}

.service-name {
  font-weight: 500;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
}

.service-desc {
  font-size: 13px;
  color: var(--el-text-color-regular);
}

.empty-state {
  padding: 40px 20px;
  text-align: center;
}

/* 深色主题适配 */
.dark .status-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
