<template>
  <el-card class="quick-actions-card">
    <template #header>
      <div class="card-header">
        <span class="card-title">
          <el-icon><Operation /></el-icon>
          快速操作
        </span>
      </div>
    </template>
    <div class="actions-grid">
      <div class="action-item" @click="createProxy">
        <div class="action-icon">
          <el-icon><Plus /></el-icon>
        </div>
        <div class="action-text">新建代理</div>
      </div>
      <div class="action-item" @click="createRule">
        <div class="action-icon">
          <el-icon><Setting /></el-icon>
        </div>
        <div class="action-text">创建规则</div>
      </div>
      <div class="action-item" @click="viewLogs">
        <div class="action-icon">
          <el-icon><Document /></el-icon>
        </div>
        <div class="action-text">查看日志</div>
      </div>
      <div class="action-item" @click="systemSettings">
        <div class="action-icon">
          <el-icon><Tools /></el-icon>
        </div>
        <div class="action-text">系统设置</div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Operation, Plus, Setting, Document, Tools } from '@element-plus/icons-vue'

const emit = defineEmits<{
  createProxy: []
  createRule: []
  viewLogs: []
  systemSettings: []
}>()

const createProxy = () => {
  emit('createProxy')
}

const createRule = () => {
  emit('createRule')
}

const viewLogs = () => {
  emit('viewLogs')
}

const systemSettings = () => {
  emit('systemSettings')
}
</script>

<style scoped>
.quick-actions-card {
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

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  padding: 20px 0;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px;
  border: 2px dashed var(--el-border-color);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-item:hover {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
  transform: translateY(-2px);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--el-color-primary-light-8);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: var(--el-color-primary);
  transition: all 0.3s ease;
}

.action-item:hover .action-icon {
  background: var(--el-color-primary);
  color: white;
}

.action-text {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

/* 深色主题适配 */
.dark .quick-actions-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

.dark .action-item:hover {
  background: var(--el-color-primary-dark-2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .actions-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
