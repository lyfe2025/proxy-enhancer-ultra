<template>
  <el-card class="table-card">
    <template #header>
      <div class="card-header">
        <span>代理列表</span>
        <div class="header-actions">
          <el-button size="small" @click="handleBatchTest" :disabled="!selectedProxies.length">
            <el-icon><VideoPlay /></el-icon>
            批量测试
          </el-button>
          <el-button size="small" type="danger" @click="handleBatchDelete" :disabled="!selectedProxies.length">
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
        </div>
      </div>
    </template>

    <el-table
      v-loading="loading"
      :data="proxyList"
      @selection-change="handleSelectionChange"
      stripe
      style="width: 100%"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="代理名称" min-width="150">
        <template #default="{ row }">
          <div class="proxy-name">
            <el-icon class="proxy-icon"><Connection /></el-icon>
            <span>{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)">{{ row.type.toUpperCase() }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="host" label="地址" min-width="200" />
      <el-table-column prop="port" label="端口" width="80" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="responseTime" label="响应时间" width="120">
        <template #default="{ row }">
          <span v-if="row.responseTime">{{ row.responseTime }}ms</span>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="lastCheck" label="最后检测" width="160">
        <template #default="{ row }">
          <span v-if="row.lastCheck">{{ formatDate(row.lastCheck) }}</span>
          <span v-else class="text-muted">未检测</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleTest(row)">
            <el-icon><VideoPlay /></el-icon>
            测试
          </el-button>
          <el-button size="small" @click="handleEdit(row)">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">
            <el-icon><Delete /></el-icon>
            删除
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
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  Connection,
  VideoPlay,
  Edit,
  Delete
} from '@element-plus/icons-vue'
import type { ProxyConfig } from '@/api/proxy'
import { getTypeTagType, getStatusTagType, getStatusText, formatDate } from '../utils/proxyUtils'

interface Pagination {
  page: number
  size: number
  total: number
}

const props = defineProps<{
  proxyList: ProxyConfig[]
  loading: boolean
  pagination: Pagination
}>()

const emit = defineEmits<{
  'update:pagination': [value: Pagination]
  test: [proxy: ProxyConfig]
  batchTest: [proxies: ProxyConfig[]]
  edit: [proxy: ProxyConfig]
  delete: [proxy: ProxyConfig]
  batchDelete: [proxies: ProxyConfig[]]
  selectionChange: [selection: ProxyConfig[]]
}>()

const selectedProxies = ref<ProxyConfig[]>([])

const handleSelectionChange = (selection: ProxyConfig[]) => {
  selectedProxies.value = selection
  emit('selectionChange', selection)
}

const handleTest = (proxy: ProxyConfig) => {
  emit('test', proxy)
}

const handleBatchTest = () => {
  emit('batchTest', selectedProxies.value)
}

const handleEdit = (proxy: ProxyConfig) => {
  emit('edit', proxy)
}

const handleDelete = (proxy: ProxyConfig) => {
  emit('delete', proxy)
}

const handleBatchDelete = () => {
  emit('batchDelete', selectedProxies.value)
}

const handleSizeChange = (size: number) => {
  const newPagination = { ...props.pagination, size, page: 1 }
  emit('update:pagination', newPagination)
}

const handleCurrentChange = (page: number) => {
  const newPagination = { ...props.pagination, page }
  emit('update:pagination', newPagination)
}


</script>

<style scoped>
.table-card {
  border: 1px solid var(--el-border-color-light);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.proxy-icon {
  color: var(--el-color-primary);
}

.text-muted {
  color: var(--el-text-color-placeholder);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 深色主题适配 */
.dark .table-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
