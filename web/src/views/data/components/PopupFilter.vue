<template>
  <el-card class="filter-card">
    <el-form :model="searchForm" inline>
      <el-form-item label="搜索">
        <el-input
          v-model="searchForm.keyword"
          placeholder="输入弹窗名称或描述"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="类型">
        <el-select v-model="searchForm.type" placeholder="选择弹窗类型" clearable>
          <el-option label="信息收集" value="form" />
          <el-option label="问卷调查" value="survey" />
          <el-option label="反馈收集" value="feedback" />
          <el-option label="订阅邮件" value="newsletter" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
          <el-option label="启用" value="active" />
          <el-option label="禁用" value="inactive" />
          <el-option label="草稿" value="draft" />
        </el-select>
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
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { Search, RefreshLeft } from '@element-plus/icons-vue'

interface SearchForm {
  keyword: string
  type: string
  status: string
}

const searchForm = defineModel<SearchForm>('searchForm', {
  required: true
})

const emit = defineEmits<{
  search: []
  reset: []
}>()

const handleSearch = () => {
  emit('search')
}

const handleReset = () => {
  searchForm.value = {
    keyword: '',
    type: '',
    status: ''
  }
  emit('reset')
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
