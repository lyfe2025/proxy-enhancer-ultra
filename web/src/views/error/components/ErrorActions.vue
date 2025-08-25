<template>
  <div class="error-actions">
    <el-button type="primary" size="large" @click="goHome">
      <el-icon><HomeFilled /></el-icon>
      返回首页
    </el-button>
    <el-button size="large" @click="goBack">
      <el-icon><ArrowLeft /></el-icon>
      返回上页
    </el-button>
    <el-button size="large" @click="refresh">
      <el-icon><Refresh /></el-icon>
      刷新页面
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HomeFilled, ArrowLeft, Refresh } from '@element-plus/icons-vue'

const router = useRouter()

const emit = defineEmits<{
  goHome: []
  goBack: []
  refresh: []
}>()

const goHome = () => {
  router.push('/')
  emit('goHome')
}

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
  emit('goBack')
}

const refresh = () => {
  location.reload()
  emit('refresh')
}
</script>

<style scoped>
.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 40px;
}

.error-actions .el-button {
  min-width: 120px;
  border-radius: 25px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.error-actions .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .error-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .error-actions .el-button {
    width: 200px;
  }
}
</style>
