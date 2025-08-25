<template>
  <div class="data-table">
    <!-- 筛选栏 -->
    <div v-if="showFilter" class="filter-bar">
      <div class="filter-left">
        <slot name="filter-left">
          <el-input
            v-if="searchable"
            v-model="searchQuery"
            :placeholder="searchPlaceholder"
            prefix-icon="Search"
            clearable
            style="width: 300px; margin-right: 16px"
            @input="handleSearch"
          />
        </slot>
      </div>
      <div class="filter-right">
        <slot name="filter-right">
          <el-button v-if="showRefresh" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </slot>
      </div>
    </div>

    <!-- 数据表格 -->
    <el-table
      ref="tableRef"
      :data="tableData"
      v-loading="loading"
      :stripe="stripe"
      :border="border"
      :selection="selection"
      style="width: 100%"
      @selection-change="handleSelectionChange"
      v-bind="$attrs"
    >
      <!-- 选择列 -->
      <el-table-column v-if="selection" type="selection" width="55" />
      
      <!-- 动态列 -->
      <el-table-column
        v-for="column in columns"
        :key="column.prop"
        :prop="column.prop"
        :label="column.label"
        :width="column.width"
        :min-width="column.minWidth"
        :fixed="column.fixed"
        :sortable="column.sortable"
        :formatter="column.formatter"
      >
        <template v-if="column.slot" #default="scope">
          <slot :name="column.slot" :row="scope.row" :column="column" :$index="scope.$index" />
        </template>
      </el-table-column>

      <!-- 操作列 -->
      <el-table-column v-if="showActions" :label="actionsLabel" :width="actionsWidth" :fixed="actionsFixed">
        <template #default="scope">
          <slot name="actions" :row="scope.row" :$index="scope.$index">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </slot>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="showPagination" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="pageSizes"
        :total="total"
        :layout="paginationLayout"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Refresh, Search } from '@element-plus/icons-vue'

// 接口定义
export interface TableColumn {
  prop: string
  label: string
  width?: number | string
  minWidth?: number | string
  fixed?: boolean | string
  sortable?: boolean
  formatter?: Function
  slot?: string // 自定义插槽名称
}

export interface PaginationData {
  page: number
  size: number
  total: number
}

// Props
interface Props {
  data: any[]
  columns: TableColumn[]
  loading?: boolean
  stripe?: boolean
  border?: boolean
  selection?: boolean
  
  // 筛选相关
  showFilter?: boolean
  searchable?: boolean
  searchPlaceholder?: string
  showRefresh?: boolean
  
  // 操作列相关
  showActions?: boolean
  actionsLabel?: string
  actionsWidth?: number | string
  actionsFixed?: boolean | string
  
  // 分页相关
  showPagination?: boolean
  pagination?: PaginationData
  pageSizes?: number[]
  paginationLayout?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  stripe: true,
  border: false,
  selection: false,
  showFilter: true,
  searchable: true,
  searchPlaceholder: '搜索...',
  showRefresh: true,
  showActions: true,
  actionsLabel: '操作',
  actionsWidth: 150,
  actionsFixed: 'right',
  showPagination: true,
  pageSizes: () => [10, 20, 50, 100],
  paginationLayout: 'total, sizes, prev, pager, next, jumper'
})

// Emits
interface Emits {
  'selection-change': [selection: any[]]
  'edit': [row: any]
  'delete': [row: any]
  'refresh': []
  'search': [query: string]
  'page-change': [page: number]
  'size-change': [size: number]
}

const emit = defineEmits<Emits>()

// 响应式数据
const tableRef = ref()
const searchQuery = ref('')
const selectedRows = ref<any[]>([])

// 分页数据
const currentPage = ref(props.pagination?.page || 1)
const pageSize = ref(props.pagination?.size || 20)
const total = ref(props.pagination?.total || 0)

// 计算属性
const tableData = computed(() => {
  if (!props.searchable || !searchQuery.value) {
    return props.data
  }
  
  // 简单的搜索逻辑，可以根据需要扩展
  return props.data.filter(item => {
    return Object.values(item).some(value => 
      String(value).toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  })
})

// 监听分页数据变化
watch(() => props.pagination, (newPagination) => {
  if (newPagination) {
    currentPage.value = newPagination.page
    pageSize.value = newPagination.size
    total.value = newPagination.total
  }
}, { immediate: true, deep: true })

// 方法
const handleSelectionChange = (selection: any[]) => {
  selectedRows.value = selection
  emit('selection-change', selection)
}

const handleEdit = (row: any) => {
  emit('edit', row)
}

const handleDelete = (row: any) => {
  emit('delete', row)
}

const handleRefresh = () => {
  emit('refresh')
}

const handleSearch = () => {
  emit('search', searchQuery.value)
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  emit('size-change', size)
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  emit('page-change', page)
}

// 暴露方法
const clearSelection = () => {
  tableRef.value?.clearSelection()
}

const toggleRowSelection = (row: any, selected?: boolean) => {
  tableRef.value?.toggleRowSelection(row, selected)
}

defineExpose({
  clearSelection,
  toggleRowSelection,
  selectedRows
})
</script>

<style scoped>
.data-table {
  width: 100%;
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 16px;
  background: var(--el-bg-color-page);
  border-radius: 8px;
}

.filter-left {
  display: flex;
  align-items: center;
}

.filter-right {
  display: flex;
  gap: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

/* 深色主题适配 */
.dark .filter-bar {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
