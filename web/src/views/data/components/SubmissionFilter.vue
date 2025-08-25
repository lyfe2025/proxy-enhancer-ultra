<template>
  <el-card class="filter-card">
    <el-form :model="searchForm" inline>
      <el-form-item label="弹窗">
        <el-select v-model="searchForm.popupId" placeholder="选择弹窗" clearable>
          <el-option
            v-for="popup in popupsList"
            :key="popup.id"
            :label="popup.name"
            :value="popup.id || 0"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
          <el-option label="待处理" value="pending" />
          <el-option label="已处理" value="processed" />
          <el-option label="已忽略" value="ignored" />
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
        <el-button @click="handleReset">
          <el-icon><RefreshLeft /></el-icon>
          重置
        </el-button>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { Search, RefreshLeft, Download } from '@element-plus/icons-vue'

interface PopupItem {
  id: number
  name: string
}

interface SearchForm {
  popupId: number
  status: string
  dateRange: [string, string] | null
}

defineProps<{
  popupsList: PopupItem[]
}>()

const searchForm = defineModel<SearchForm>('searchForm', {
  required: true
})

const emit = defineEmits<{
  search: []
  reset: []
  export: []
}>()

const handleSearch = () => {
  emit('search')
}

const handleReset = () => {
  searchForm.value = {
    popupId: 0,
    status: '',
    dateRange: null
  }
  emit('reset')
}

const handleExport = () => {
  emit('export')
}
</script>

<style scoped>
.filter-card {
  margin-bottom: 20px;
}

/* 深色主题适配 */
.dark .filter-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
