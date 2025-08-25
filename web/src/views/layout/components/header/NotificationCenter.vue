<template>
  <el-popover placement="bottom" width="320" trigger="click">
    <template #reference>
      <el-badge :value="unreadCount" :hidden="unreadCount === 0">
        <el-button text circle>
          <el-icon size="18"><Bell /></el-icon>
        </el-button>
      </el-badge>
    </template>

    <div class="notifications-panel">
      <div class="notifications-header">
        <span class="notifications-title">通知中心</span>
        <el-button text size="small" @click="markAllAsRead">
          全部已读
        </el-button>
      </div>

      <el-scrollbar max-height="300px">
        <div class="notifications-list">
          <div
            v-for="notification in notifications"
            :key="notification.id"
            class="notification-item"
            :class="{ 'notification-unread': !notification.read }"
            @click="markAsRead(notification.id)"
          >
            <div class="notification-content">
              <div class="notification-title">{{ notification.title }}</div>
              <div class="notification-desc">{{ notification.content }}</div>
              <div class="notification-time">{{ formatTime(notification.createdAt) }}</div>
            </div>
          </div>
          <div v-if="!notifications.length" class="notifications-empty">
            <el-empty description="暂无通知" />
          </div>
        </div>
      </el-scrollbar>
    </div>
  </el-popover>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Bell } from '@element-plus/icons-vue'
import { formatDistanceToNow } from 'date-fns'
import { zhCN } from 'date-fns/locale'

interface NotificationItem {
  id: string
  title: string
  content: string
  type: 'info' | 'warning' | 'error' | 'success'
  read: boolean
  createdAt: string
}

const props = defineProps<{
  notifications: NotificationItem[]
}>()

const emit = defineEmits<{
  markAsRead: [id: string]
  markAllAsRead: []
}>()

// 计算未读通知数量
const unreadCount = computed(() => {
  return props.notifications.filter(n => !n.read).length
})

const markAsRead = (id: string) => {
  emit('markAsRead', id)
}

const markAllAsRead = () => {
  emit('markAllAsRead')
}

const formatTime = (timeString: string) => {
  return formatDistanceToNow(new Date(timeString), {
    addSuffix: true,
    locale: zhCN
  })
}
</script>

<style scoped>
/* 通知面板样式 */
.notifications-panel {
  padding: 0;
}

.notifications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.notifications-title {
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.notifications-list {
  padding: 8px 0;
}

.notification-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.3s ease;
  border-left: 3px solid transparent;
}

.notification-item:hover {
  background: var(--el-fill-color-light);
}

.notification-unread {
  border-left-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.notification-title {
  font-weight: 500;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
  font-size: 13px;
}

.notification-desc {
  color: var(--el-text-color-regular);
  font-size: 12px;
  margin-bottom: 4px;
  line-height: 1.4;
}

.notification-time {
  color: var(--el-text-color-placeholder);
  font-size: 11px;
}

.notifications-empty {
  padding: 20px;
  text-align: center;
}
</style>
