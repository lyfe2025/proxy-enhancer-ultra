<template>
  <div class="proxy-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">
          <el-icon><Connection /></el-icon>
          代理管理
        </h1>
        <p class="page-description">管理和配置代理服务器，监控代理状态和性能</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="() => { editingProxy = null; showFormDialog = true }">
          <el-icon><Plus /></el-icon>
          添加代理
        </el-button>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <ProxyStats :stats="stats" />

    <!-- 搜索和筛选 -->
    <ProxySearchFilter
      :model-value="searchForm"
      @update:model-value="(value) => Object.assign(searchForm, value)"
      @search="handleSearch"
      @reset="resetSearch"
    />

    <!-- 代理列表 -->
    <ProxyList
      :proxy-list="proxyList"
      :loading="loading"
      :pagination="pagination"
      @update:pagination="updatePagination"
      @test="testProxy"
      @batch-test="batchTest"
      @edit="editProxy"
      @delete="handleDeleteProxy"
      @batch-delete="handleBatchDelete"
      @selection-change="handleSelectionChange"
    />

    <!-- 添加/编辑代理对话框 -->
    <ProxyFormDialog
      v-model:visible="showFormDialog"
      :edit-data="editingProxy"
      :submitting="saving"
      @submit="handleSaveProxy"
      @close="resetForm"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { Connection, Plus, Refresh } from '@element-plus/icons-vue'

// 导入组合式函数
import { useProxyManagement } from './composables/useProxyManagement'
import { useProxyStats } from './composables/useProxyStats'
import { useProxyForm } from './composables/useProxyForm'

// 导入子组件
import ProxyStats from './components/ProxyStats.vue'
import ProxySearchFilter from './components/ProxySearchFilter.vue'
import ProxyList from './components/ProxyList.vue'
import ProxyFormDialog from './components/ProxyFormDialog.vue'

// 使用组合式函数
const {
  loading,
  proxyList,
  searchForm,
  pagination,
  fetchProxyList,
  handleSearch,
  resetSearch,
  updatePagination,
  testProxy,
  batchTest,
  deleteProxy,
  batchDelete,
  handleSelectionChange
} = useProxyManagement()

const {
  stats,
  fetchStats
} = useProxyStats()

const {
  saving,
  showFormDialog,
  editingProxy,
  editProxy,
  saveProxy,
  resetForm
} = useProxyForm()

// 刷新数据
const refreshData = () => {
  fetchProxyList()
  fetchStats()
}

// 保存代理处理（需要在保存后刷新数据）
const handleSaveProxy = async (form: any) => {
  const success = await saveProxy(form)
  if (success) {
    fetchProxyList()
    fetchStats()
  }
}

// 删除代理处理（需要在删除后刷新统计）
const handleDeleteProxy = async (proxy: any) => {
  await deleteProxy(proxy)
  fetchStats()
}

// 批量删除处理（需要在删除后刷新统计）
const handleBatchDelete = async (proxies: any[]) => {
  await batchDelete(proxies)
  fetchStats()
}

// 初始化
onMounted(() => {
  fetchProxyList()
  fetchStats()
})
</script>

<style scoped>
.proxy-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.header-left {
  flex: 1;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.page-description {
  color: var(--el-text-color-regular);
  margin: 0;
  font-size: 16px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* 深色主题适配 */
.dark .proxy-management {
  background: var(--el-bg-color-page);
}
</style>