<template>
  <div class="popup-management">
    <!-- 搜索筛选 -->
    <PopupFilter
      v-model:search-form="searchForm"
      @search="handleSearch"
      @reset="resetSearch"
    />

    <!-- 弹窗列表 -->
    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span>弹窗列表</span>
          <div class="header-actions">
            <el-button size="small" @click="batchEnable" :disabled="!selectedItems.length">
              <el-icon><CircleCheck /></el-icon>
              批量启用
            </el-button>
            <el-button size="small" @click="batchDisable" :disabled="!selectedItems.length">
              <el-icon><CircleClose /></el-icon>
              批量禁用
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
        :data="popupsList"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="弹窗名称" min-width="150">
          <template #default="{ row }">
            <div class="popup-name">
              <el-icon class="popup-icon"><Document /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getPopupTypeTagType(row.type)">{{ getPopupTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="展示次数" width="120">
          <template #default="{ row }">
            <span>{{ row.views || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="submissions" label="提交次数" width="120">
          <template #default="{ row }">
            <span>{{ row.submissions || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="conversionRate" label="转化率" width="100">
          <template #default="{ row }">
            <span>{{ ((row.submissions || 0) / Math.max(row.views || 1, 1) * 100).toFixed(1) }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            <span>{{ formatDate(row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="preview(row)">
              <el-icon><View /></el-icon>
              预览
            </el-button>
            <el-button size="small" @click="edit(row)">
              <el-icon><Edit /></el-icon>
              编辑
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
  Document,
  CircleCheck,
  CircleClose,
  Delete,
  View,
  Edit
} from '@element-plus/icons-vue'
import PopupFilter from './PopupFilter.vue'

interface PopupItem {
  id: number
  name: string
  type: string
  status: string
  views: number
  submissions: number
  createdAt: string
}

interface SearchForm {
  keyword: string
  type: string
  status: string
}

interface PaginationData {
  page: number
  size: number
  total: number
}

const props = defineProps<{
  popupsList: PopupItem[]
  loading: boolean
  pagination: PaginationData
}>()

const emit = defineEmits<{
  search: [form: SearchForm]
  resetSearch: []
  batchEnable: [ids: number[]]
  batchDisable: [ids: number[]]
  batchDelete: [ids: number[]]
  preview: [item: PopupItem]
  edit: [item: PopupItem]
  delete: [item: PopupItem]
  sizeChange: [size: number]
  currentChange: [page: number]
}>()

const searchForm = ref<SearchForm>({
  keyword: '',
  type: '',
  status: ''
})

const selectedItems = ref<PopupItem[]>([])

const handleSearch = () => {
  emit('search', searchForm.value)
}

const resetSearch = () => {
  searchForm.value = {
    keyword: '',
    type: '',
    status: ''
  }
  emit('resetSearch')
}

const handleSelectionChange = (selection: PopupItem[]) => {
  selectedItems.value = selection
}

const batchEnable = async () => {
  try {
    await ElMessageBox.confirm('确认批量启用选中的弹窗吗？', '提示', {
      type: 'warning'
    })
    emit('batchEnable', selectedItems.value.map(item => item.id))
  } catch {}
}

const batchDisable = async () => {
  try {
    await ElMessageBox.confirm('确认批量禁用选中的弹窗吗？', '提示', {
      type: 'warning'
    })
    emit('batchDisable', selectedItems.value.map(item => item.id))
  } catch {}
}

const batchDelete = async () => {
  try {
    await ElMessageBox.confirm('确认批量删除选中的弹窗吗？此操作不可恢复！', '危险操作', {
      type: 'error'
    })
    emit('batchDelete', selectedItems.value.map(item => item.id))
  } catch {}
}

const preview = (item: PopupItem) => {
  emit('preview', item)
}

const edit = (item: PopupItem) => {
  emit('edit', item)
}

const deleteItem = async (item: PopupItem) => {
  try {
    await ElMessageBox.confirm(`确认删除弹窗"${item.name}"吗？此操作不可恢复！`, '危险操作', {
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
const getPopupTypeTagType = (type: string): "success" | "info" | "warning" | "danger" | "primary" => {
  const typeMap: Record<string, "success" | "info" | "warning" | "danger" | "primary"> = {
    form: 'primary',
    survey: 'success',
    feedback: 'warning',
    newsletter: 'info'
  }
  return typeMap[type] || 'info'
}

const getPopupTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    form: '信息收集',
    survey: '问卷调查',
    feedback: '反馈收集',
    newsletter: '订阅邮件'
  }
  return typeMap[type] || type
}

const getStatusTagType = (status: string): "success" | "info" | "warning" | "danger" | "primary" => {
  const statusMap: Record<string, "success" | "info" | "warning" | "danger" | "primary"> = {
    active: 'success',
    inactive: 'info',
    draft: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '启用',
    inactive: '禁用',
    draft: '草稿'
  }
  return statusMap[status] || status
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 组件销毁时清理弹窗
onUnmounted(() => {
  // 关闭所有可能打开的 ElMessageBox
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

.popup-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.popup-icon {
  color: var(--el-color-primary);
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
