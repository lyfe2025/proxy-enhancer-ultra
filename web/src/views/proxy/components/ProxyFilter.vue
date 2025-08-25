<template>
  <el-card class="filter-card">
    <el-form :model="form" inline>
      <el-form-item label="搜索">
        <el-input
          v-model="form.keyword"
          placeholder="输入代理名称或地址"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="类型">
        <el-select v-model="form.type" placeholder="选择代理类型" clearable>
          <el-option label="HTTP" value="http" />
          <el-option label="HTTPS" value="https" />
          <el-option label="SOCKS5" value="socks5" />
          <el-option label="SOCKS4" value="socks4" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="form.status" placeholder="选择状态" clearable>
          <el-option label="活跃" value="active" />
          <el-option label="停用" value="inactive" />
          <el-option label="异常" value="error" />
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
import { reactive, watch } from 'vue'
import { Search, RefreshLeft } from '@element-plus/icons-vue'

interface FilterForm {
  keyword: string
  type: string
  status: string
}

const props = defineProps<{
  modelValue: FilterForm
}>()

const emit = defineEmits<{
  'update:modelValue': [value: FilterForm]
  search: []
  reset: []
}>()

const form = reactive<FilterForm>({
  keyword: props.modelValue.keyword,
  type: props.modelValue.type,
  status: props.modelValue.status
})

// 监听表单变化，同步到父组件
watch(
  form,
  (newValue) => {
    emit('update:modelValue', { ...newValue })
  },
  { deep: true }
)

// 监听父组件传入的值变化
watch(
  () => props.modelValue,
  (newValue) => {
    Object.assign(form, newValue)
  },
  { deep: true }
)

const handleSearch = () => {
  emit('search')
}

const handleReset = () => {
  Object.assign(form, {
    keyword: '',
    type: '',
    status: ''
  })
  emit('reset')
}
</script>

<style scoped>
.filter-card {
  margin-bottom: 20px;
  border: 1px solid var(--el-border-color-light);
}

/* 深色主题适配 */
.dark .filter-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
