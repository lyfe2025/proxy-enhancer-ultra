<template>
  <div class="proxy-list">
    <el-card class="list-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><List /></el-icon>
            代理列表
          </span>
          <div class="header-actions">
            <el-button
              v-if="selectedProxies.length > 0"
              type="warning"
              size="small"
              @click="handleBatchTest"
            >
              <el-icon><VideoPlay /></el-icon>
              批量测试 ({{ selectedProxies.length }})
            </el-button>
            <el-button
              v-if="selectedProxies.length > 0"
              type="danger"
              size="small"
              @click="handleBatchDelete"
            >
              <el-icon><Delete /></el-icon>
              批量删除 ({{ selectedProxies.length }})
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        :data="proxyList"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="代理名称" min-width="150">
          <template #default="{ row }">
            <div class="proxy-name">
              <span class="name-text">{{ row.name }}</span>
              <el-tag
                :type="getTypeTagType(row.type)"
                size="small"
                class="type-tag"
              >
                {{ row.type.toUpperCase() }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="host" label="主机地址" min-width="150" />
        
        <el-table-column prop="port" label="端口" width="80" />
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="response_time" label="响应时间" width="100">
          <template #default="{ row }">
            <span v-if="row.response_time">
              {{ row.response_time }}ms
            </span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="traffic" label="流量" width="100">
          <template #default="{ row }">
            <span v-if="row.traffic">
              {{ formatBytes(row.traffic) }}
            </span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="last_used" label="最后使用" width="150">
          <template #default="{ row }">
            <span v-if="row.last_used">
              {{ formatDate(row.last_used) }}
            </span>
            <span v-else class="text-muted">从未使用</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="description" label="描述" min-width="150">
          <template #default="{ row }">
            <span v-if="row.description" class="description-text">
              {{ row.description }}
            </span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                type="primary"
                size="small"
                @click="handleTest(row)"
                :loading="row.testing"
              >
                <el-icon><VideoPlay /></el-icon>
                测试
              </el-button>
              <el-button
                type="warning"
                size="small"
                @click="handleEdit(row)"
              >
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(row)"
              >
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPagination.page"
          v-model:page-size="currentPagination.size"
          :page-sizes="PAGINATION_CONFIG.PAGE_SIZES"
          :total="currentPagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { List, VideoPlay, Edit, Delete } from '@element-plus/icons-vue'
import type { ProxyConfig } from '@/api/proxy'
import { getTypeTagType, getStatusTagType, getStatusText, formatBytes, formatDate } from '../utils/proxyUtils'
import { PAGINATION_CONFIG } from '../constants/proxyConstants'

// 定义 Pagination 接口
interface Pagination {
  page: number
  size: number
  total: number
}

// 定义 props
const props = defineProps<{
  proxyList: ProxyConfig[]
  loading: boolean
  pagination: Pagination
}>()

// 定义 emits
const emit = defineEmits<{
  'update:pagination': [pagination: Pagination]
  test: [proxy: ProxyConfig]
  'batch-test': [proxies: ProxyConfig[]]
  edit: [proxy: ProxyConfig]
  delete: [proxy: ProxyConfig]
  'batch-delete': [proxies: ProxyConfig[]]
  'selection-change': [selection: ProxyConfig[]]
}>()

// 本地分页数据
const currentPagination = reactive({ ...props.pagination })

// 选中的代理
const selectedProxies = ref<ProxyConfig[]>([])

// 监听 props 分页变化
watch(
  () => props.pagination,
  (newValue) => {
    Object.assign(currentPagination, newValue)
  },
  { deep: true }
)

// 选择变化处理
const handleSelectionChange = (selection: ProxyConfig[]) => {
  selectedProxies.value = selection
  emit('selection-change', selection)
}

// 测试代理
const handleTest = (proxy: ProxyConfig) => {
  emit('test', proxy)
}

// 批量测试
const handleBatchTest = () => {
  if (selectedProxies.value.length > 0) {
    emit('batch-test', selectedProxies.value)
  }
}

// 编辑代理
const handleEdit = (proxy: ProxyConfig) => {
  emit('edit', proxy)
}

// 删除代理
const handleDelete = (proxy: ProxyConfig) => {
  emit('delete', proxy)
}

// 批量删除
const handleBatchDelete = () => {
  if (selectedProxies.value.length > 0) {
    emit('batch-delete', selectedProxies.value)
  }
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  currentPagination.size = size
  currentPagination.page = 1
  emit('update:pagination', { ...currentPagination })
}

// 当前页变化
const handleCurrentChange = (page: number) => {
  currentPagination.page = page
  emit('update:pagination', { ...currentPagination })
}
</script>

<style scoped>
.proxy-list {
  margin-bottom: 24px;
}

.list-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.header-actions {
  display: flex;
  gap: 8px;
}

.proxy-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name-text {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.type-tag {
  font-size: 10px;
  padding: 2px 6px;
}

.description-text {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}

.text-muted {
  color: var(--el-text-color-placeholder);
}

.action-buttons {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.action-buttons .el-button {
  padding: 4px 8px;
  font-size: 12px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--el-border-color-lighter);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }
  
  .action-buttons .el-button {
    width: 100%;
    justify-content: center;
  }
  
  .pagination-wrapper {
    overflow-x: auto;
  }
}

/* 深色主题适配 */
.dark .list-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
  box-shadow: 0 2px 8px rgba(255, 255, 255, 0.1);
}

.dark .pagination-wrapper {
  border-top-color: var(--el-border-color);
}
</style>