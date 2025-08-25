<template>
  <div class="operation-logs">
    <!-- 搜索筛选 -->
    <el-card class="filter-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="用户">
          <el-input
            v-model="searchForm.username"
            placeholder="输入用户名"
            clearable
          />
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select v-model="searchForm.action" placeholder="选择操作类型" clearable>
            <el-option label="登录" value="login" />
            <el-option label="登出" value="logout" />
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshLeft /></el-icon>
            重置
          </el-button>
          <el-button @click="exportLogs">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 日志列表 -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="logsList" stripe style="width: 100%">
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="action" label="操作" width="100">
          <template #default="{ row }">
            <el-tag :type="getActionTagType(row.action)" size="small">
              {{ getActionText(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="资源" width="150" />
        <el-table-column prop="action" label="描述" min-width="200" />
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="user_agent" label="用户代理" width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.user_agent" placement="top">
              <span class="truncated">{{ row.user_agent }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="160">
          <template #default="{ row }">
            <span>{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">
              <el-icon><View /></el-icon>
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="操作日志详情"
      width="600px"
    >
      <div v-if="selectedLog" class="log-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="用户">{{ selectedLog.username }}</el-descriptions-item>
          <el-descriptions-item label="操作">
            <el-tag :type="getActionTagType(selectedLog.action)">
              {{ getActionText(selectedLog.action) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="模块">{{ selectedLog.resource }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ selectedLog.ip_address }}</el-descriptions-item>
          <el-descriptions-item label="操作时间" :span="2">
            {{ formatDate(selectedLog.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="用户代理" :span="2">
            {{ selectedLog.user_agent }}
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ selectedLog.action }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div v-if="selectedLog.error_message" class="log-details">
          <h4>错误信息</h4>
          <pre>{{ selectedLog.error_message }}</pre>
        </div>
      </div>
      <template #footer>
        <el-button @click="showDetailDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Search,
  RefreshLeft,
  Download,
  View
} from '@element-plus/icons-vue'
import * as systemApi from '@/api/system'
import type { OperationLog } from '@/types/system'

// Emits
const emit = defineEmits(['refresh'])

// 响应式数据
const loading = ref(false)
const showDetailDialog = ref(false)
const selectedLog = ref<OperationLog | null>(null)

// 搜索表单
const searchForm = reactive({
  username: '',
  action: '',
  dateRange: null as [string, string] | null
})

// 列表数据
const logsList = ref<OperationLog[]>([])

// 分页
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 获取操作日志
const fetchLogs = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      size: pagination.size,
      username: searchForm.username,
      action: searchForm.action,
      startTime: searchForm.dateRange?.[0],
      endTime: searchForm.dateRange?.[1]
    }
    const response = await systemApi.getOperationLogs(params)
    logsList.value = (response.data as any).data || []
    pagination.total = (response.data as any).total || 0
  } catch (error) {
    ElMessage.error('获取操作日志失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  fetchLogs()
}

const resetSearch = () => {
  Object.assign(searchForm, {
    username: '',
    action: '',
    dateRange: null
  })
  handleSearch()
}

// 分页处理
const handleSizeChange = (size: number) => {
  pagination.size = size
  fetchLogs()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchLogs()
}

// 查看日志详情
const viewDetail = (log: OperationLog) => {
  selectedLog.value = log
  showDetailDialog.value = true
}

// 导出日志
const exportLogs = async () => {
  try {
    const params = {
      username: searchForm.username,
      action: searchForm.action,
      startTime: searchForm.dateRange?.[0],
      endTime: searchForm.dateRange?.[1]
    }
    const response = await systemApi.getOperationLogs(params)
    // 处理文件下载
    const blob = new Blob([JSON.stringify(response.data)], { type: 'application/json' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `operation_logs_${new Date().toISOString().split('T')[0]}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 工具函数
const getActionTagType = (action: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const actionMap: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    login: 'success',
    logout: 'info',
    create: 'primary',
    update: 'warning',
    delete: 'danger'
  }
  return actionMap[action] || 'info'
}

const getActionText = (action: string) => {
  const actionMap: Record<string, string> = {
    login: '登录',
    logout: '登出',
    create: '创建',
    update: '更新',
    delete: '删除'
  }
  return actionMap[action] || action
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 暴露方法给父组件
defineExpose({
  fetchLogs
})

// 初始化
onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.filter-card {
  margin-bottom: 20px;
}

.table-card {
  border: 1px solid var(--el-border-color-light);
}

.truncated {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.log-detail {
  padding: 20px 0;
}

.log-details {
  margin-top: 20px;
}

.log-details h4 {
  margin: 0 0 12px 0;
  color: var(--el-text-color-primary);
}

.log-details pre {
  background: var(--el-bg-color-page);
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  padding: 12px;
  font-size: 12px;
  max-height: 300px;
  overflow-y: auto;
}

/* 深色主题适配 */
.dark .table-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

.dark .filter-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

.dark .log-details pre {
  background: var(--el-bg-color);
  border-color: var(--el-border-color);
  color: var(--el-text-color-primary);
}
</style>
