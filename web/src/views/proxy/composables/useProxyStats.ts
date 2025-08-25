import { reactive } from 'vue'
import { proxyApi } from '@/api/proxy'

// 统计数据接口
interface ProxyStats {
  total: number
  active: number
  error: number
  traffic: number
}

export function useProxyStats() {
  // 统计数据
  const stats = reactive<ProxyStats>({
    total: 0,
    active: 0,
    error: 0,
    traffic: 0
  })

  // 获取统计数据
  const fetchStats = async () => {
    try {
      const response = await proxyApi.getProxyStats()
      Object.assign(stats, response.data)
    } catch (error) {
      console.error('获取统计数据失败:', error)
    }
  }

  // 刷新统计数据
  const refreshStats = () => {
    fetchStats()
  }

  return {
    // 响应式数据
    stats,
    
    // 方法
    fetchStats,
    refreshStats
  }
}