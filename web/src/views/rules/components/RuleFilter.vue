<template>
  <el-card class="filter-card">
    <el-form :model="searchForm" inline>
      <el-form-item label="搜索">
        <el-input
          v-model="searchForm.keyword"
          placeholder="输入规则名称或描述"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="类型">
        <el-select v-model="searchForm.type" placeholder="选择规则类型" clearable>
          <el-option label="路由规则" value="route" />
          <el-option label="过滤规则" value="filter" />
          <el-option label="重写规则" value="rewrite" />
          <el-option label="限流规则" value="rate_limit" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
          <el-option label="启用" value="enabled" />
          <el-option label="禁用" value="disabled" />
        </el-select>
      </el-form-item>
      <el-form-item label="优先级">
        <el-select v-model="searchForm.priority" placeholder="选择优先级" clearable>
          <el-option label="高" value="high" />
          <el-option label="中" value="medium" />
          <el-option label="低" value="low" />
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

interface RuleSearchForm {
  keyword: string
  type: string
  status: string
  priority: string
}

const searchForm = defineModel<RuleSearchForm>('searchForm', {
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
    status: '',
    priority: ''
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
