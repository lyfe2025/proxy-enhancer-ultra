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
        <el-button type="primary" @click="showAddDialog = true">
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
    <ProxyFilter
      v-model="searchForm"
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
      @delete="deleteProxy"
      @batch-delete="batchDelete"
      @selection-change="handleSelectionChange"
    />

    <!-- 添加/编辑代理对话框 -->
    <ProxyFormDialog
      v-model="showAddDialog"
      :editing-proxy="editingProxy"
      :saving="saving"
      @save="saveProxy"
      @close="resetForm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Connection,
  Plus,
  Refresh,
  Search,
  RefreshLeft,
  VideoPlay,
  Delete,
  Edit,
  CircleCheck,
  CircleClose,
  TrendCharts
} from '@element-plus/icons-vue'
import { proxyApi, type ProxyConfig } from '@/api/proxy'
import type { ProxyCreateRequest, ProxyUpdateRequest } from '@/types/proxy'
import { getTypeTagType, getStatusTagType, getStatusText, formatBytes, formatDate } from './utils/proxyUtils'
import { PROXY_TYPE_OPTIONS, PROXY_STATUS_OPTIONS, PAGINATION_CONFIG, DEFAULT_PROXY_FORM, DEFAULT_SEARCH_FORM, PROXY_FORM_RULES, PORT_RANGE } from './constants/proxyConstants'

// 表单验证规则
const proxyFormRules = PROXY_FORM_RULES

// 导入子组件
import ProxyStats from './components/ProxyStats.vue'
import ProxyFilter from './components/ProxyFilter.vue'
import ProxyList from './components/ProxyList.vue'
import ProxyFormDialog from './components/ProxyFormDialog.vue'

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const showAddDialog = ref(false)
const editingProxy = ref<ProxyConfig | null>(null)

// 统计数据
const stats = reactive({
  total: 0,
  active: 0,
  error: 0,
  traffic: 0
})

// 搜索表单
const searchForm = reactive({ ...DEFAULT_SEARCH_FORM })

// 代理列表
const proxyList = ref<ProxyConfig[]>([])

// 代理表单
const proxyForm = reactive({ ...DEFAULT_PROXY_FORM })

// 分页
interface Pagination {
  page: number
  size: number
  total: number
}

const pagination = reactive<Pagination>({
  page: PAGINATION_CONFIG.DEFAULT_PAGE,
  size: PAGINATION_CONFIG.DEFAULT_SIZE,
  total: 0
})

// 获取代理列表
const fetchProxyList = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      size: pagination.size,
      keyword: searchForm.keyword,
      type: searchForm.type,
      status: searchForm.status
    }
    const response = await proxyApi.getProxies(params)
    proxyList.value = (response.data as any).data || []
    pagination.total = (response.data as any).total || 0
  } catch (error) {
    ElMessage.error('获取代理列表失败')
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await proxyApi.getProxyStats()
    Object.assign(stats, response.data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  fetchProxyList()
}

// 重置搜索
const resetSearch = () => {
  Object.assign(searchForm, DEFAULT_SEARCH_FORM)
  handleSearch()
}

// 刷新数据
const refreshData = () => {
  fetchProxyList()
  fetchStats()
}

// 更新分页
const updatePagination = (newPagination: Pagination) => {
  Object.assign(pagination, newPagination)
  fetchProxyList()
}

// 选择变化
const handleSelectionChange = (selection: ProxyConfig[]) => {
  // 在这里处理选择变化，如果需要的话
}

// 测试代理
const testProxy = async (proxy: ProxyConfig) => {
  try {
    if (!proxy.id) {
      ElMessage.error('代理ID无效')
      return
    }
    await proxyApi.testProxy(proxy.id)
    ElMessage.success('代理测试成功')
    fetchProxyList()
  } catch (error) {
    ElMessage.error('代理测试失败')
  }
}

// 批量测试
const batchTest = async (proxies: ProxyConfig[]) => {
  try {
    const ids = proxies.map(p => p.id).filter(id => id !== undefined) as number[]
    for (const id of ids) {
      await proxyApi.testProxy(id)
    }
    ElMessage.success('批量测试已启动')
    fetchProxyList()
  } catch (error) {
    ElMessage.error('批量测试失败')
  }
}

// 编辑代理
const editProxy = (proxy: ProxyConfig) => {
  editingProxy.value = proxy
  showAddDialog.value = true
}

// 删除代理
const deleteProxy = async (proxy: ProxyConfig) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除代理 "${proxy.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    if (!proxy.id) {
      ElMessage.error('代理ID无效')
      return
    }
    await proxyApi.deleteProxy(proxy.id)
    ElMessage.success('删除成功')
    fetchProxyList()
    fetchStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const batchDelete = async (proxies: ProxyConfig[]) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${proxies.length} 个代理吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    const ids = proxies.map(p => p.id).filter(id => id !== undefined) as number[]
    for (const id of ids) {
      await proxyApi.deleteProxy(id)
    }
    ElMessage.success('批量删除成功')
    fetchProxyList()
    fetchStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 保存代理
const saveProxy = async (form: any) => {
  try {
    saving.value = true
    
    if (editingProxy.value) {
      // 更新代理
      const updateData: ProxyUpdateRequest = {
        name: form.name,
        type: form.type as 'http' | 'https' | 'socks5',
        host: form.host,
        port: form.port,
        username: form.username || undefined,
        password: form.password || undefined,
        description: form.description || undefined
      }
      await proxyApi.updateProxy(editingProxy.value.id!, updateData)
      ElMessage.success('更新成功')
    } else {
      // 创建代理
      const createData = {
        name: form.name,
        type: form.type,
        host: form.host,
        port: form.port,
        username: form.username || undefined,
        password: form.password || undefined,
        status: 'inactive' as const,
        description: form.description || undefined
      }
      await proxyApi.createProxy(createData)
      ElMessage.success('添加成功')
    }
    
    showAddDialog.value = false
    fetchProxyList()
    fetchStats()
  } catch (error) {
    ElMessage.error(editingProxy.value ? '更新失败' : '添加失败')
  } finally {
    saving.value = false
  }
}

// 重置表单
const resetForm = () => {
  editingProxy.value = null
  Object.assign(proxyForm, DEFAULT_PROXY_FORM)
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