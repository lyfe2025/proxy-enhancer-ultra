<template>
  <div class="submission-management">
    <!-- 搜索筛选 -->
    <SubmissionFilter
      v-model:search-form="searchForm"
      :popups-list="popupsList"
      @search="handleSearch"
      @reset="resetSearch"
      @export="handleExport"
    />

    <!-- 提交数据列表 -->
    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span>提交数据</span>
          <div class="header-actions">
            <el-button size="small" @click="batchProcess" :disabled="!selectedItems.length">
              <el-icon><CircleCheck /></el-icon>
              批量处理
            </el-button>
            <el-button size="small" type="danger" @click="batchDelete" :disabled="!selectedItems.length">
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="submissionsList"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="popupName" label="弹窗名称" width="150" />
        <el-table-column prop="data" label="提交数据" min-width="300">
          <template #default="{ row }">
            <div class="submission-data">
              <div
                v-for="(value, key) in row.data"
                :key="key"
                class="data-item"
              >
                <span class="data-key">{{ key }}:</span>
                <span class="data-value">{{ value }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getSubmissionStatusTagType(row.status)" size="small">
              {{ getSubmissionStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="userAgent" label="用户代理" width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.userAgent" placement="top">
              <span class="truncated">{{ row.userAgent }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="submittedAt" label="提交时间" width="160">
          <template #default="{ row }">
            <span>{{ formatDate(row.submittedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">
              <el-icon><View /></el-icon>
              详情
            </el-button>
            <el-button size="small" @click="updateStatus(row)">
              <el-icon><Edit /></el-icon>
              处理
            </el-button>
            <el-button size="small" type="danger" @click="deleteItem(row)">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  CircleCheck,
  Delete,
  View,
  Edit
} from '@element-plus/icons-vue'
import SubmissionFilter from './SubmissionFilter.vue'

interface SubmissionItem {
  id: number
  popupName: string
  data: Record<string, any>
  status: string
  ip: string
  userAgent: string
  submittedAt: string
}

interface PopupItem {
  id: number
  name: string
}

interface SearchForm {
  popupId: number
  status: string
  dateRange: [string, string] | null
}

interface PaginationData {
  page: number
  size: number
  total: number
}

const props = defineProps<{
  submissionsList: SubmissionItem[]
  popupsList: PopupItem[]
  loading: boolean
  pagination: PaginationData
}>()

const emit = defineEmits<{
  search: [form: SearchForm]
  resetSearch: []
  export: []
  batchProcess: [ids: number[]]
  batchDelete: [ids: number[]]
  viewDetail: [item: SubmissionItem]
  updateStatus: [item: SubmissionItem]
  delete: [item: SubmissionItem]
  sizeChange: [size: number]
  currentChange: [page: number]
}>()

const searchForm = ref<SearchForm>({
  popupId: 0,
  status: '',
  dateRange: null
})

const selectedItems = ref<SubmissionItem[]>([])

const handleSearch = () => {
  emit('search', searchForm.value)
}

const resetSearch = () => {
  searchForm.value = {
    popupId: 0,
    status: '',
    dateRange: null
  }
  emit('resetSearch')
}

const handleExport = () => {
  emit('export')
}

const handleSelectionChange = (selection: SubmissionItem[]) => {
  selectedItems.value = selection
}

const batchProcess = async () => {
  try {
    await ElMessageBox.confirm('确认批量处理选中的提交数据吗？', '提示', {
      type: 'warning'
    })
    emit('batchProcess', selectedItems.value.map(item => item.id))
  } catch {}
}

const batchDelete = async () => {
  try {
    await ElMessageBox.confirm('确认批量删除选中的提交数据吗？此操作不可恢复！', '危险操作', {
      type: 'error'
    })
    emit('batchDelete', selectedItems.value.map(item => item.id))
  } catch {}
}

const viewDetail = (item: SubmissionItem) => {
  emit('viewDetail', item)
}

const updateStatus = (item: SubmissionItem) => {
  emit('updateStatus', item)
}

const deleteItem = async (item: SubmissionItem) => {
  try {
    await ElMessageBox.confirm('确认删除这条提交数据吗？此操作不可恢复！', '危险操作', {
      type: 'error'
    })
    emit('delete', item)
  } catch {}
}

const handleSizeChange = (size: number) => {
  emit('sizeChange', size)
}

const handleCurrentChange = (page: number) => {
  emit('currentChange', page)
}

// 工具函数
const getSubmissionStatusTagType = (status: string): "success" | "info" | "warning" | "danger" | "primary" => {
  const statusMap: Record<string, "success" | "info" | "warning" | "danger" | "primary"> = {
    pending: 'warning',
    processed: 'success',
    ignored: 'info'
  }
  return statusMap[status] || 'info'
}

const getSubmissionStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    processed: '已处理',
    ignored: '已忽略'
  }
  return statusMap[status] || status
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 组件销毁时清理弹窗
onUnmounted(() => {
  try {
    ElMessageBox.close()
  } catch (error) {
    // 忽略关闭错误
  }
})
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

.submission-data {
  max-height: 100px;
  overflow-y: auto;
}

.data-item {
  display: flex;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 12px;
}

.data-key {
  font-weight: 600;
  color: var(--el-text-color-regular);
  min-width: 60px;
}

.data-value {
  color: var(--el-text-color-primary);
  word-break: break-all;
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

/* 深色主题适配 */
.dark .table-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
