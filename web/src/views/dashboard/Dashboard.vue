<template>
  <div class="dashboard">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon class="title-icon"><Odometer /></el-icon>
        仪表盘
      </h1>
      <p class="page-description">系统运行状态总览</p>
    </div>

    <!-- 统计卡片 -->
    <DashboardStats :stats="stats" />

    <!-- 图表区域 -->
    <DashboardCharts 
      :chart-data="chartData" 
      @update-period="updateTrafficPeriod"
    />

    <!-- 底部信息区域 -->
    <div class="info-grid">
      <el-row :gutter="20">
        <el-col :span="12">
          <DashboardActivity 
            :activities="recentActivities"
            @refresh="loadActivities"
          />
        </el-col>
        <el-col :span="12">
          <DashboardSystemStatus 
            :services="systemServices"
            @refresh="loadSystemStatus"
          />
        </el-col>
      </el-row>
    </div>

    <!-- 快速操作 -->
    <DashboardQuickActions
      @create-proxy="createProxy"
      @create-rule="createRule"
      @view-logs="viewLogs"
      @system-settings="systemSettings"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Odometer } from '@element-plus/icons-vue'

// 导入子组件
import DashboardStats from './components/DashboardStats.vue'
import DashboardCharts from './components/DashboardCharts.vue'
import DashboardActivity from './components/DashboardActivity.vue'
import DashboardSystemStatus from './components/DashboardSystemStatus.vue'
import DashboardQuickActions from './components/DashboardQuickActions.vue'

// 导入API
import { dashboardApi } from '@/api/dashboard'

const router = useRouter()

// 响应式数据
const stats = reactive({
  totalProxies: 0,
  proxyGrowth: 0,
  activeRules: 0,
  rulesGrowth: 0,
  totalTraffic: 0,
  trafficGrowth: 0,
  onlineUsers: 0,
  userChange: 0
})

const chartData = reactive({
  trafficData: [],
  proxyStatusData: []
})

const recentActivities = ref<any[]>([])
const systemServices = ref<any[]>([])

let refreshTimer: NodeJS.Timeout | null = null

// 初始化
onMounted(() => {
  loadDashboardData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})

// 加载所有仪表盘数据
const loadDashboardData = async () => {
  await Promise.all([
    loadStats(),
    loadChartData('today'),
    loadActivities(),
    loadSystemStatus()
  ])
}

// 加载统计数据
const loadStats = async () => {
  try {
    // 实际实现应该调用统计API
    // const response = await dashboardApi.getStats()
    
    // 模拟数据
    Object.assign(stats, {
      totalProxies: 42,
      proxyGrowth: 12.5,
      activeRules: 28,
      rulesGrowth: 8.3,
      totalTraffic: 1024 * 1024 * 1024 * 2.5, // 2.5GB
      trafficGrowth: 15.7,
      onlineUsers: 156,
      userChange: -2.1
    })
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
  }
}

// 加载图表数据
const loadChartData = async (period: string) => {
  try {
    // 实际实现应该调用图表API
    // const response = await dashboardApi.getChartData(period)
    
    // 模拟流量趋势数据
    const trafficData = []
    const now = new Date()
    for (let i = 23; i >= 0; i--) {
      const time = new Date(now.getTime() - i * 60 * 60 * 1000)
      trafficData.push({
        time: time.getHours() + ':00',
        inbound: Math.random() * 1000000 + 500000,
        outbound: Math.random() * 800000 + 400000
      })
    }
    
    // 模拟代理状态数据
    const proxyStatusData = [
      { value: 25, name: '运行中' },
      { value: 12, name: '空闲' },
      { value: 3, name: '故障' },
      { value: 2, name: '维护中' }
    ]
    
    Object.assign(chartData, {
      trafficData,
      proxyStatusData
    })
  } catch (error) {
    console.error('加载图表数据失败:', error)
    ElMessage.error('加载图表数据失败')
  }
}

// 加载活动日志
const loadActivities = async () => {
  try {
    // 实际实现应该调用活动API
    // const response = await dashboardApi.getActivities()
    
    // 模拟活动数据
    recentActivities.value = [
      {
        id: '1',
        type: 'success',
        title: '新代理服务器上线',
        description: '代理服务器 proxy-001 已成功启动并开始处理请求',
        time: new Date(Date.now() - 5 * 60 * 1000).toISOString()
      },
      {
        id: '2',
        type: 'warning',
        title: '规则匹配率下降',
        description: '路由规则 rule-003 的匹配率在过去1小时内下降了15%',
        time: new Date(Date.now() - 15 * 60 * 1000).toISOString()
      },
      {
        id: '3',
        type: 'info',
        title: '系统维护完成',
        description: '定期系统维护已完成，所有服务恢复正常运行',
        time: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString()
      }
    ]
  } catch (error) {
    console.error('加载活动日志失败:', error)
    ElMessage.error('加载活动日志失败')
  }
}

// 加载系统状态
const loadSystemStatus = async () => {
  try {
    // 实际实现应该调用系统状态API
    // const response = await dashboardApi.getSystemStatus()
    
    // 模拟系统服务状态
    systemServices.value = [
      {
        name: '代理服务',
        description: '核心代理转发服务',
        status: 'running'
      },
      {
        name: '规则引擎',
        description: '请求路由规则处理引擎',
        status: 'running'
      },
      {
        name: '监控服务',
        description: '系统监控和日志收集服务',
        status: 'running'
      },
      {
        name: '缓存服务',
        description: 'Redis缓存服务',
        status: 'warning'
      },
      {
        name: '数据库',
        description: 'PostgreSQL数据库服务',
        status: 'running'
      }
    ]
  } catch (error) {
    console.error('加载系统状态失败:', error)
    ElMessage.error('加载系统状态失败')
  }
}

// 更新流量趋势周期
const updateTrafficPeriod = (period: string) => {
  loadChartData(period)
}

// 自动刷新
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    loadStats()
    loadChartData('today')
    loadSystemStatus()
  }, 30000) // 30秒刷新一次
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 快速操作事件处理
const createProxy = () => {
  router.push('/proxy/create')
}

const createRule = () => {
  router.push('/rules/create')
}

const viewLogs = () => {
  router.push('/logs')
}

const systemSettings = () => {
  router.push('/settings')
}
</script>

<style scoped>
.dashboard {
  padding: 20px;
  min-height: 100vh;
  background: var(--el-bg-color-page);
}

.page-header {
  margin-bottom: 24px;
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

.title-icon {
  font-size: 32px;
  color: var(--el-color-primary);
}

.page-description {
  color: var(--el-text-color-regular);
  margin: 0;
  font-size: 16px;
}

.info-grid {
  margin-bottom: 24px;
}

/* 深色主题适配 */
.dark .dashboard {
  background: var(--el-bg-color-page);
}
</style>