import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { proxyApi, type ProxyConfig } from '@/api/proxy'
import type { ProxyCreateRequest, ProxyUpdateRequest } from '@/types/proxy'
import { PAGINATION_CONFIG, DEFAULT_SEARCH_FORM } from '../constants/proxyConstants'

// 分页接口
interface Pagination {
  page: number
  size: number
  total: number
}

// 搜索表单接口
interface SearchForm {
  keyword: string
  type: string
  status: string
}

export function useProxyManagement() {
  // 响应式数据
  const loading = ref(false)
  const proxyList = ref<ProxyConfig[]>([])
  
  // 搜索表单
  const searchForm = reactive<SearchForm>({ ...DEFAULT_SEARCH_FORM })
  
  // 分页
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

  // 更新分页
  const updatePagination = (newPagination: Pagination) => {
    Object.assign(pagination, newPagination)
    fetchProxyList()
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
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('批量删除失败')
      }
    }
  }

  // 选择变化处理
  const handleSelectionChange = (selection: ProxyConfig[]) => {
    // 在这里处理选择变化，如果需要的话
  }

  return {
    // 响应式数据
    loading,
    proxyList,
    searchForm,
    pagination,
    
    // 方法
    fetchProxyList,
    handleSearch,
    resetSearch,
    updatePagination,
    testProxy,
    batchTest,
    deleteProxy,
    batchDelete,
    handleSelectionChange
  }
}