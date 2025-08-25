<template>
  <div class="rules-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">
          <el-icon><SetUp /></el-icon>
          规则配置
        </h1>
        <p class="page-description">配置和管理代理规则，控制流量路由和处理逻辑</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showAddDialog = true">
          <el-icon><Plus /></el-icon>
          添加规则
        </el-button>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <RuleStats :stats="stats" />

    <!-- 搜索筛选 -->
    <RuleFilter
      v-model:search-form="searchForm"
      @search="handleSearch"
      @reset="resetSearch"
    />

    <!-- 规则列表 -->
    <RuleList
      :rules-list="rulesList"
      :loading="loading"
      :pagination="pagination"
      @batch-enable="batchEnableRules"
      @batch-disable="batchDisableRules"
      @batch-delete="batchDeleteRules"
      @toggle-status="toggleRuleStatus"
      @move-up="moveRuleUp"
      @move-down="moveRuleDown"
      @edit="editRule"
      @delete="deleteRule"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />

    <!-- 添加/编辑规则对话框 -->
    <RuleFormDialog
      v-model="showAddDialog"
      :editing-rule="editingRule"
      @submit="handleRuleSubmit"
      @close="resetForm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { SetUp, Plus, Refresh } from '@element-plus/icons-vue'

// 导入子组件
import RuleStats from './components/RuleStats.vue'
import RuleFilter from './components/RuleFilter.vue'
import RuleList from './components/RuleList.vue'
import RuleFormDialog from './components/RuleFormDialog.vue'

// 导入API
import { rulesApi } from '@/api/rules'

// 响应式数据
const showAddDialog = ref(false)
const editingRule = ref(null)
const loading = ref(false)

// 统计数据
const stats = reactive({
  total: 0,
  active: 0,
  executed: 0,
  successRate: 0
})

// 搜索表单
const searchForm = ref({
  keyword: '',
  type: '',
  status: '',
  priority: ''
})

// 规则列表
const rulesList = ref<any[]>([])
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 初始化
onMounted(() => {
  loadStats()
  loadRules()
})

// 加载统计数据
const loadStats = async () => {
  try {
    // 实际实现应该调用统计API
    stats.total = 25
    stats.active = 18
    stats.executed = 1562
    stats.successRate = 0.947
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载规则列表
const loadRules = async () => {
  loading.value = true
  try {
    // 实际实现应该调用规则API
    const response = await rulesApi.getRules({
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm.value
    })
    rulesList.value = response.data.data?.list || []
    pagination.total = response.data.data?.total || 0
  } catch (error) {
    console.error('加载规则列表失败:', error)
    ElMessage.error('加载规则列表失败')
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  loadStats()
  loadRules()
  ElMessage.success('数据已刷新')
}

// 搜索事件处理
const handleSearch = () => {
  pagination.page = 1
  loadRules()
}

const resetSearch = () => {
  searchForm.value = {
    keyword: '',
    type: '',
    status: '',
    priority: ''
  }
  pagination.page = 1
  loadRules()
}

// 批量操作
const batchEnableRules = (ids: string[]) => {
  console.log('批量启用规则:', ids)
  ElMessage.success('批量启用成功')
  loadRules()
}

const batchDisableRules = (ids: string[]) => {
  console.log('批量禁用规则:', ids)
  ElMessage.success('批量禁用成功')
  loadRules()
}

const batchDeleteRules = (ids: string[]) => {
  console.log('批量删除规则:', ids)
  ElMessage.success('批量删除成功')
  loadRules()
}

// 单个规则操作
const toggleRuleStatus = (rule: any) => {
  console.log('切换规则状态:', rule)
  ElMessage.success('状态切换成功')
  loadRules()
}

const moveRuleUp = (rule: any) => {
  console.log('上移规则:', rule)
  ElMessage.success('规则优先级已调整')
  loadRules()
}

const moveRuleDown = (rule: any) => {
  console.log('下移规则:', rule)
  ElMessage.success('规则优先级已调整')
  loadRules()
}

const editRule = (rule: any) => {
  editingRule.value = rule
  showAddDialog.value = true
}

const deleteRule = (rule: any) => {
  console.log('删除规则:', rule)
  ElMessage.success('删除成功')
  loadRules()
}

// 分页处理
const handleSizeChange = (size: number) => {
  pagination.size = size
  loadRules()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadRules()
}

// 表单处理
const handleRuleSubmit = (formData: any) => {
  console.log('提交规则表单:', formData)
  ElMessage.success(editingRule.value ? '规则更新成功' : '规则创建成功')
  showAddDialog.value = false
  resetForm()
  loadRules()
}

const resetForm = () => {
  editingRule.value = null
}
</script>

<style scoped>
.rules-management {
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
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.page-description {
  color: var(--el-text-color-regular);
  margin: 0;
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* 深色主题适配 */
.dark .rules-management {
  background: var(--el-bg-color-page);
}
</style>