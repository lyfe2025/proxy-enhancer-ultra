<template>
  <div class="proxy-filter">
    <el-card class="filter-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><Search /></el-icon>
            搜索筛选
          </span>
        </div>
      </template>
      
      <el-form :model="searchForm" inline class="filter-form">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索代理名称、主机地址"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="代理类型">
          <el-select
            v-model="searchForm.type"
            placeholder="选择类型"
            clearable
            style="width: 120px"
          >
            <el-option label="全部" value="" />
            <el-option
              v-for="option in PROXY_TYPE_OPTIONS"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="选择状态"
            clearable
            style="width: 120px"
          >
            <el-option
              v-for="option in PROXY_STATUS_OPTIONS"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
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
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import { Search, RefreshLeft } from '@element-plus/icons-vue'
import { PROXY_TYPE_OPTIONS, PROXY_STATUS_OPTIONS, DEFAULT_SEARCH_FORM } from '../constants/proxyConstants'

// 定义 props
interface SearchForm {
  keyword: string
  type: string
  status: string
}

const props = defineProps<{
  modelValue: SearchForm
}>()

// 定义 emits
const emit = defineEmits<{
  'update:modelValue': [value: SearchForm]
  search: []
  reset: []
}>()

// 本地搜索表单
const searchForm = reactive({ ...props.modelValue })

// 监听 props 变化
watch(
  () => props.modelValue,
  (newValue) => {
    Object.assign(searchForm, newValue)
  },
  { deep: true }
)

// 监听本地表单变化，同步到父组件
watch(
  searchForm,
  (newValue) => {
    emit('update:modelValue', { ...newValue })
  },
  { deep: true }
)

// 搜索处理
const handleSearch = () => {
  emit('search')
}

// 重置处理
const handleReset = () => {
  Object.assign(searchForm, DEFAULT_SEARCH_FORM)
  emit('reset')
}
</script>

<style scoped>
.proxy-filter {
  margin-bottom: 24px;
}

.filter-card {
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

.filter-form {
  margin: 0;
}

.filter-form .el-form-item {
  margin-bottom: 0;
  margin-right: 16px;
}

.filter-form .el-form-item:last-child {
  margin-right: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .filter-form {
    display: block;
  }
  
  .filter-form .el-form-item {
    display: block;
    margin-bottom: 16px;
    margin-right: 0;
  }
  
  .filter-form .el-form-item:last-child {
    margin-bottom: 0;
  }
  
  .filter-form .el-input,
  .filter-form .el-select {
    width: 100% !important;
  }
}

/* 深色主题适配 */
.dark .filter-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
  box-shadow: 0 2px 8px rgba(255, 255, 255, 0.1);
}
</style>