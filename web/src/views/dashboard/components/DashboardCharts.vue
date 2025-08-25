<template>
  <div class="charts-grid">
    <!-- 流量趋势图 -->
    <el-card class="chart-card traffic-chart">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><TrendCharts /></el-icon>
            流量趋势
          </span>
          <el-select v-model="trafficPeriod" size="small" style="width: 120px" @change="updateTrafficChart">
            <el-option label="今日" value="today" />
            <el-option label="本周" value="week" />
            <el-option label="本月" value="month" />
          </el-select>
        </div>
      </template>
      <div class="chart-container">
        <div ref="trafficChartRef" class="chart"></div>
      </div>
    </el-card>

    <!-- 代理状态分布 -->
    <el-card class="chart-card proxy-status">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><PieChart /></el-icon>
            代理状态分布
          </span>
        </div>
      </template>
      <div class="chart-container">
        <div ref="proxyStatusChartRef" class="chart"></div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { TrendCharts, PieChart } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

interface ChartData {
  trafficData: any[]
  proxyStatusData: any[]
}

const props = defineProps<{
  chartData: ChartData
}>()

const emit = defineEmits<{
  updatePeriod: [period: string]
}>()

const trafficPeriod = ref('today')
const trafficChartRef = ref<HTMLElement>()
const proxyStatusChartRef = ref<HTMLElement>()

let trafficChart: echarts.ECharts | null = null
let proxyStatusChart: echarts.ECharts | null = null

onMounted(() => {
  initCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  destroyCharts()
  window.removeEventListener('resize', handleResize)
})

// 监听数据变化
watch(
  () => props.chartData,
  () => {
    updateCharts()
  },
  { deep: true }
)

const initCharts = () => {
  if (trafficChartRef.value) {
    trafficChart = echarts.init(trafficChartRef.value)
    updateTrafficChart()
  }
  
  if (proxyStatusChartRef.value) {
    proxyStatusChart = echarts.init(proxyStatusChartRef.value)
    updateProxyStatusChart()
  }
}

const updateTrafficChart = () => {
  if (!trafficChart) return
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: props.chartData.trafficData.map(item => item.time),
      axisTick: {
        alignWithLabel: true
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value: number) => formatBytes(value)
      }
    },
    series: [
      {
        name: '入站流量',
        type: 'line',
        data: props.chartData.trafficData.map(item => item.inbound),
        smooth: true,
        lineStyle: {
          color: '#409eff'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
            { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
          ])
        }
      },
      {
        name: '出站流量',
        type: 'line',
        data: props.chartData.trafficData.map(item => item.outbound),
        smooth: true,
        lineStyle: {
          color: '#67c23a'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
            { offset: 1, color: 'rgba(103, 194, 58, 0.1)' }
          ])
        }
      }
    ]
  }
  
  trafficChart.setOption(option)
  emit('updatePeriod', trafficPeriod.value)
}

const updateProxyStatusChart = () => {
  if (!proxyStatusChart) return
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        name: '代理状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '30',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: props.chartData.proxyStatusData
      }
    ]
  }
  
  proxyStatusChart.setOption(option)
}

const updateCharts = () => {
  updateTrafficChart()
  updateProxyStatusChart()
}

const handleResize = () => {
  trafficChart?.resize()
  proxyStatusChart?.resize()
}

const destroyCharts = () => {
  trafficChart?.dispose()
  proxyStatusChart?.dispose()
  trafficChart = null
  proxyStatusChart = null
}

// 格式化字节数
const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped>
.charts-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

.chart-card {
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.chart-container {
  height: 300px;
  width: 100%;
}

.chart {
  width: 100%;
  height: 100%;
}

/* 深色主题适配 */
.dark .chart-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
