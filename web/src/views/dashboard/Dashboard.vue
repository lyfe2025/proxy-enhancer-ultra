<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.key">
        <div class="stat-icon">
          <el-icon :size="32" :color="stat.color">
            <component :is="stat.icon" />
          </el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-change" :class="stat.changeType">
            <el-icon><component :is="stat.changeIcon" /></el-icon>
            {{ stat.change }}
          </div>
        </div>
      </div>
    </div>
    
    <!-- 图表区域 -->
    <div class="charts-grid">
      <!-- 代理状态分布 -->
      <el-card class="chart-card">
        <template #header>
          <div class="card-header">
            <span>代理状态分布</span>
            <el-button type="text" @click="refreshProxyChart">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </template>
        <div class="chart-container" ref="proxyChartRef"></div>
      </el-card>
      
      <!-- 流量趋势 -->
      <el-card class="chart-card">
        <template #header>
          <div class="card-header">
            <span>流量趋势</span>
            <el-select v-model="trafficPeriod" size="small" style="width: 100px">
              <el-option label="7天" value="7d" />
              <el-option label="30天" value="30d" />
              <el-option label="90天" value="90d" />
            </el-select>
          </div>
        </template>
        <div class="chart-container" ref="trafficChartRef"></div>
      </el-card>
    </div>
    
    <!-- 活动列表 -->
    <div class="activity-grid">
      <!-- 最近活动 -->
      <el-card class="activity-card">
        <template #header>
          <div class="card-header">
            <span>最近活动</span>
            <el-button type="text" @click="refreshActivities">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </template>
        <div class="activity-list">
          <div class="activity-item" v-for="activity in activities" :key="activity.id">
            <div class="activity-icon">
              <el-icon :color="activity.color">
                <component :is="activity.icon" />
              </el-icon>
            </div>
            <div class="activity-content">
              <div class="activity-title">{{ activity.title }}</div>
              <div class="activity-desc">{{ activity.description }}</div>
              <div class="activity-time">{{ formatTime(activity.time) }}</div>
            </div>
          </div>
          <div v-if="activities.length === 0" class="empty-state">
            <el-empty description="暂无活动记录" />
          </div>
        </div>
      </el-card>
      
      <!-- 系统状态 -->
      <el-card class="status-card">
        <template #header>
          <div class="card-header">
            <span>系统状态</span>
            <el-tag :type="systemStatus.type" size="small">
              {{ systemStatus.text }}
            </el-tag>
          </div>
        </template>
        <div class="status-list">
          <div class="status-item" v-for="item in systemMetrics" :key="item.name">
            <div class="status-label">{{ item.label }}</div>
            <div class="status-value">
              <el-progress
                :percentage="item.value"
                :color="getProgressColor(item.value)"
                :show-text="false"
                :stroke-width="6"
              />
              <span class="status-text">{{ item.value }}%</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Connection,
  Setting,
  DataAnalysis,
  User,
  TrendCharts,
  ArrowUp,
  ArrowDown,
  Refresh,
  SuccessFilled,
  WarningFilled,
  CircleCloseFilled
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { systemApi } from '@/api'
import type { DashboardStats, SystemMetric } from '@/types'

// 响应式数据
const stats = ref<DashboardStats[]>([])
const activities = ref<any[]>([])
const systemMetrics = ref<SystemMetric[]>([])
const trafficPeriod = ref('7d')
const loading = ref(false)

// 图表引用
const proxyChartRef = ref<HTMLElement>()
const trafficChartRef = ref<HTMLElement>()
let proxyChart: echarts.ECharts | null = null
let trafficChart: echarts.ECharts | null = null

// 系统状态
const systemStatus = reactive({
  type: 'success' as 'success' | 'warning' | 'danger',
  text: '正常运行'
})

// 获取仪表盘数据
const fetchDashboardData = async () => {
  try {
    loading.value = true
    const response = await systemApi.getDashboardStats()
    
    // 更新统计数据
    stats.value = [
      {
        key: 'proxies',
        label: '代理总数',
        value: response.data.proxies.total,
        change: `+${response.data.proxies.change}`,
        changeType: response.data.proxies.change >= 0 ? 'positive' : 'negative',
        changeIcon: response.data.proxies.change >= 0 ? 'ArrowUp' : 'ArrowDown',
        icon: 'Connection',
        color: '#00ff88'
      },
      {
        key: 'rules',
        label: '规则总数',
        value: response.data.rules.total,
        change: `+${response.data.rules.change}`,
        changeType: response.data.rules.change >= 0 ? 'positive' : 'negative',
        changeIcon: response.data.rules.change >= 0 ? 'ArrowUp' : 'ArrowDown',
        icon: 'Setting',
        color: '#409eff'
      },
      {
        key: 'submissions',
        label: '数据提交',
        value: response.data.submissions.total,
        change: `+${response.data.submissions.change}`,
        changeType: response.data.submissions.change >= 0 ? 'positive' : 'negative',
        changeIcon: response.data.submissions.change >= 0 ? 'ArrowUp' : 'ArrowDown',
        icon: 'DataAnalysis',
        color: '#e6a23c'
      },
      {
        key: 'users',
        label: '用户总数',
        value: response.data.users.total,
        change: `+${response.data.users.change}`,
        changeType: response.data.users.change >= 0 ? 'positive' : 'negative',
        changeIcon: response.data.users.change >= 0 ? 'ArrowUp' : 'ArrowDown',
        icon: 'User',
        color: '#f56c6c'
      }
    ]
    
  } catch (error: any) {
    ElMessage.error(error.message || '获取仪表盘数据失败')
  } finally {
    loading.value = false
  }
}

// 获取系统指标
const fetchSystemMetrics = async () => {
  try {
    const response = await systemApi.getSystemMetrics()
    systemMetrics.value = response.data
    
    // 更新系统状态
    const cpuUsage = response.data.find(m => m.name === 'cpu')?.value || 0
    const memoryUsage = response.data.find(m => m.name === 'memory')?.value || 0
    
    if (cpuUsage > 80 || memoryUsage > 80) {
      systemStatus.type = 'danger'
      systemStatus.text = '资源紧张'
    } else if (cpuUsage > 60 || memoryUsage > 60) {
      systemStatus.type = 'warning'
      systemStatus.text = '负载较高'
    } else {
      systemStatus.type = 'success'
      systemStatus.text = '正常运行'
    }
    
  } catch (error: any) {
    ElMessage.error(error.message || '获取系统指标失败')
  }
}

// 获取活动记录
const fetchActivities = async () => {
  try {
    // 模拟活动数据
    activities.value = [
      {
        id: 1,
        title: '新增代理配置',
        description: '用户 admin 创建了新的代理配置 "测试代理"',
        time: new Date(Date.now() - 1000 * 60 * 5),
        icon: 'Connection',
        color: '#00ff88'
      },
      {
        id: 2,
        title: '规则更新',
        description: '规则 "URL过滤" 已更新并生效',
        time: new Date(Date.now() - 1000 * 60 * 15),
        icon: 'Setting',
        color: '#409eff'
      },
      {
        id: 3,
        title: '数据收集',
        description: '收到 25 条新的表单提交数据',
        time: new Date(Date.now() - 1000 * 60 * 30),
        icon: 'DataAnalysis',
        color: '#e6a23c'
      },
      {
        id: 4,
        title: '用户登录',
        description: '用户 test_user 登录系统',
        time: new Date(Date.now() - 1000 * 60 * 45),
        icon: 'User',
        color: '#f56c6c'
      }
    ]
  } catch (error: any) {
    ElMessage.error(error.message || '获取活动记录失败')
  }
}

// 初始化代理状态图表
const initProxyChart = () => {
  if (!proxyChartRef.value) return
  
  proxyChart = echarts.init(proxyChartRef.value, 'dark')
  
  const option = {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(26, 26, 26, 0.9)',
      borderColor: '#333',
      textStyle: { color: '#fff' }
    },
    series: [
      {
        name: '代理状态',
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '50%'],
        data: [
          { value: 35, name: '运行中', itemStyle: { color: '#00ff88' } },
          { value: 8, name: '已停止', itemStyle: { color: '#f56c6c' } },
          { value: 3, name: '错误', itemStyle: { color: '#e6a23c' } }
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        label: {
          color: '#ccc'
        }
      }
    ]
  }
  
  proxyChart.setOption(option)
}

// 初始化流量趋势图表
const initTrafficChart = () => {
  if (!trafficChartRef.value) return
  
  trafficChart = echarts.init(trafficChartRef.value, 'dark')
  
  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(26, 26, 26, 0.9)',
      borderColor: '#333',
      textStyle: { color: '#fff' }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisLine: { lineStyle: { color: '#333' } },
      axisLabel: { color: '#ccc' }
    },
    yAxis: {
      type: 'value',
      axisLine: { lineStyle: { color: '#333' } },
      axisLabel: { color: '#ccc' },
      splitLine: { lineStyle: { color: '#333' } }
    },
    series: [
      {
        name: '请求数',
        type: 'line',
        smooth: true,
        data: [120, 132, 101, 134, 90, 230, 210],
        lineStyle: { color: '#00ff88' },
        itemStyle: { color: '#00ff88' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(0, 255, 136, 0.3)' },
            { offset: 1, color: 'rgba(0, 255, 136, 0.1)' }
          ])
        }
      }
    ]
  }
  
  trafficChart.setOption(option)
}

// 刷新图表
const refreshProxyChart = () => {
  if (proxyChart) {
    proxyChart.dispose()
    initProxyChart()
  }
}

const refreshActivities = () => {
  fetchActivities()
}

// 格式化时间
const formatTime = (time: Date) => {
  const now = new Date()
  const diff = now.getTime() - time.getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

// 获取进度条颜色
const getProgressColor = (value: number) => {
  if (value > 80) return '#f56c6c'
  if (value > 60) return '#e6a23c'
  return '#00ff88'
}

// 监听流量周期变化
watch(trafficPeriod, () => {
  // 重新获取流量数据并更新图表
  if (trafficChart) {
    // 这里可以根据周期获取不同的数据
    initTrafficChart()
  }
})

// 窗口大小变化时重新调整图表
const handleResize = () => {
  if (proxyChart) proxyChart.resize()
  if (trafficChart) trafficChart.resize()
}

// 组件挂载
onMounted(() => {
  fetchDashboardData()
  fetchSystemMetrics()
  fetchActivities()
  
  // 初始化图表
  setTimeout(() => {
    initProxyChart()
    initTrafficChart()
  }, 100)
  
  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
  
  // 定时刷新数据
  const interval = setInterval(() => {
    fetchSystemMetrics()
  }, 30000) // 30秒刷新一次
  
  // 保存定时器引用
  ;(window as any).dashboardInterval = interval
})

// 组件卸载
onUnmounted(() => {
  if (proxyChart) proxyChart.dispose()
  if (trafficChart) trafficChart.dispose()
  
  window.removeEventListener('resize', handleResize)
  
  // 清除定时器
  if ((window as any).dashboardInterval) {
    clearInterval((window as any).dashboardInterval)
  }
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%);
  border: 1px solid #333;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 255, 136, 0.15);
  border-color: #00ff88;
}

.stat-icon {
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #888;
  margin-bottom: 8px;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.stat-change.positive {
  color: #00ff88;
}

.stat-change.negative {
  color: #f56c6c;
}

/* 图表网格 */
.charts-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 20px;
  margin-bottom: 20px;
}

.chart-card {
  background-color: #1a1a1a;
  border: 1px solid #333;
}

.chart-card :deep(.el-card__header) {
  background-color: transparent;
  border-bottom: 1px solid #333;
  padding: 16px 20px;
}

.chart-card :deep(.el-card__body) {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #fff;
  font-weight: 500;
}

.chart-container {
  height: 300px;
  width: 100%;
}

/* 活动网格 */
.activity-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.activity-card,
.status-card {
  background-color: #1a1a1a;
  border: 1px solid #333;
}

.activity-card :deep(.el-card__header),
.status-card :deep(.el-card__header) {
  background-color: transparent;
  border-bottom: 1px solid #333;
  padding: 16px 20px;
}

.activity-card :deep(.el-card__body),
.status-card :deep(.el-card__body) {
  padding: 0;
}

/* 活动列表 */
.activity-list {
  max-height: 400px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #333;
  transition: background-color 0.3s ease;
}

.activity-item:hover {
  background-color: rgba(0, 255, 136, 0.05);
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  margin-bottom: 4px;
}

.activity-desc {
  font-size: 13px;
  color: #888;
  margin-bottom: 4px;
  line-height: 1.4;
}

.activity-time {
  font-size: 12px;
  color: #666;
}

.empty-state {
  padding: 40px 20px;
  text-align: center;
}

/* 系统状态 */
.status-list {
  padding: 20px;
}

.status-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.status-item:last-child {
  margin-bottom: 0;
}

.status-label {
  font-size: 14px;
  color: #ccc;
  flex-shrink: 0;
  width: 80px;
}

.status-value {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: 16px;
}

.status-text {
  font-size: 12px;
  color: #888;
  min-width: 35px;
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .activity-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .stat-card {
    padding: 16px;
  }
  
  .stat-value {
    font-size: 24px;
  }
}

/* Element Plus 组件样式覆盖 */
:deep(.el-card) {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
}

:deep(.el-card__header) {
  color: #fff;
}

:deep(.el-button--text) {
  color: #ccc;
}

:deep(.el-button--text:hover) {
  color: #00ff88;
}

:deep(.el-select) {
  --el-select-input-color: #ccc;
  --el-select-border-color-hover: #00ff88;
}

:deep(.el-tag) {
  border: none;
}

:deep(.el-tag--success) {
  background-color: rgba(0, 255, 136, 0.2);
  color: #00ff88;
}

:deep(.el-tag--warning) {
  background-color: rgba(230, 162, 60, 0.2);
  color: #e6a23c;
}

:deep(.el-tag--danger) {
  background-color: rgba(245, 108, 108, 0.2);
  color: #f56c6c;
}

:deep(.el-progress-bar__outer) {
  background-color: #333;
}

:deep(.el-empty__description) {
  color: #666;
}
</style>