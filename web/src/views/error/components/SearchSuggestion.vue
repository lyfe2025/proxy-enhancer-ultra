<template>
  <div class="search-suggestion">
    <p>或者尝试搜索您需要的内容：</p>
    <el-input
      v-model="keyword"
      placeholder="输入关键词搜索..."
      size="large"
      @keyup.enter="handleSearch"
    >
      <template #append>
        <el-button @click="handleSearch">
          <el-icon><Search /></el-icon>
        </el-button>
      </template>
    </el-input>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

const router = useRouter()
const keyword = ref('')

const emit = defineEmits<{
  search: [keyword: string]
}>()

const handleSearch = () => {
  if (!keyword.value.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }
  
  // 可以跳转到搜索页面或执行搜索逻辑
  router.push({
    path: '/search',
    query: { q: keyword.value.trim() }
  })
  
  emit('search', keyword.value.trim())
}
</script>

<style scoped>
.search-suggestion {
  max-width: 400px;
  margin: 0 auto 40px auto;
  color: #cccccc;
}

.search-suggestion p {
  margin: 0 0 16px 0;
  font-size: 14px;
  text-align: center;
}

.search-suggestion .el-input {
  --el-input-border-radius: 25px;
}

.search-suggestion .el-input :deep(.el-input__wrapper) {
  border-radius: 25px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.search-suggestion .el-input :deep(.el-input__wrapper):hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.search-suggestion .el-input :deep(.el-input-group__append) {
  border-radius: 0 25px 25px 0;
  border-left: none;
}

.search-suggestion .el-input :deep(.el-input-group__append .el-button) {
  border-radius: 0 23px 23px 0;
  margin: -1px;
}
</style>
