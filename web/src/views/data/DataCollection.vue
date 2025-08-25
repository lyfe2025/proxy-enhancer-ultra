<template>
  <div class="data-collection">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">
          <el-icon><DataAnalysis /></el-icon>
          数据收集
        </h1>
        <p class="page-description">管理弹窗表单和收集的用户数据</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreatePopupDialog = true">
          <el-icon><Plus /></el-icon>
          创建弹窗
        </el-button>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <DataCollectionStats :stats="stats" />

    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="main-tabs">
      <!-- 弹窗管理 -->
      <el-tab-pane label="弹窗管理" name="popups">
        <PopupManagement
          :popups-list="popupsList"
          :loading="popupLoading"
          :pagination="popupPagination"
          @search="handlePopupSearch"
          @reset-search="resetPopupSearch"
          @batch-enable="batchEnablePopups"
          @batch-disable="batchDisablePopups" 
          @batch-delete="batchDeletePopups"
          @preview="previewPopup"
          @edit="editPopup"
          @delete="deletePopup"
          @size-change="handlePopupSizeChange"
          @current-change="handlePopupCurrentChange"
        />
      </el-tab-pane>

      <!-- 提交数据 -->
      <el-tab-pane label="提交数据" name="submissions">
        <SubmissionManagement
          :submissions-list="submissionsList"
          :popups-list="popupsList"
          :loading="submissionLoading"
          :pagination="submissionPagination"
          @search="handleSubmissionSearch"
          @reset-search="resetSubmissionSearch"
          @export="exportSubmissions"
          @batch-process="batchProcessSubmissions"
          @batch-delete="batchDeleteSubmissions"
          @view-detail="viewSubmissionDetail"
          @update-status="updateSubmissionStatus"
          @delete="deleteSubmission"
          @size-change="handleSubmissionSizeChange"
          @current-change="handleSubmissionCurrentChange"
        />
      </el-tab-pane>
    </el-tabs>

    <!-- 创建/编辑弹窗对话框 -->
    <PopupFormDialog
      v-model="showCreatePopupDialog"
      :editing-popup="editingPopup"
      @submit="handlePopupSubmit"
      @close="resetPopupForm"
    />

    <!-- 提交详情对话框 -->
    <SubmissionDetailDialog
      v-model="showSubmissionDetailDialog"
      :submission="selectedSubmission"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { DataAnalysis, Plus, Refresh } from '@element-plus/icons-vue'

// 导入子组件
import DataCollectionStats from './components/DataCollectionStats.vue'
import PopupManagement from './components/PopupManagement.vue'
import SubmissionManagement from './components/SubmissionManagement.vue'
import PopupFormDialog from './components/PopupFormDialog.vue'
import SubmissionDetailDialog from './components/SubmissionDetailDialog.vue'

// 导入API
import { popupApi } from '@/api/popup'
// import { submissionApi } from '@/api/submission'

// 响应式数据
const activeTab = ref('popups')
const showCreatePopupDialog = ref(false)
const showSubmissionDetailDialog = ref(false)
const editingPopup = ref(null)
const selectedSubmission = ref(null)

// 统计数据
const stats = reactive({
  totalPopups: 0,
  activePopups: 0,
  totalSubmissions: 0,
  conversionRate: 0
})

// 弹窗数据
const popupsList = ref<any[]>([])
const popupLoading = ref(false)
const popupPagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 提交数据
const submissionsList = ref<any[]>([])
const submissionLoading = ref(false)
const submissionPagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 初始化
onMounted(() => {
  loadStats()
  loadPopups()
  loadSubmissions()
})

// 加载统计数据
const loadStats = async () => {
  try {
    // 实际实现应该调用统计API
    stats.totalPopups = 15
    stats.activePopups = 12
    stats.totalSubmissions = 1247
    stats.conversionRate = 0.245
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载弹窗列表
const loadPopups = async () => {
  popupLoading.value = true
  try {
    // 实际实现应该调用弹窗API
    const response = await popupApi.getPopups({
      page: popupPagination.page,
      page_size: popupPagination.size
    })
    popupsList.value = response.data.data?.items || []
    popupPagination.total = response.data.data?.total || 0
  } catch (error) {
    console.error('加载弹窗列表失败:', error)
    ElMessage.error('加载弹窗列表失败')
  } finally {
    popupLoading.value = false
  }
}

// 加载提交数据
const loadSubmissions = async () => {
  submissionLoading.value = true
  try {
    // 实际实现应该调用提交API
    // const response = await submissionApi.getSubmissionList({
    //   page: submissionPagination.page,
    //   pageSize: submissionPagination.size
    // })
    // submissionsList.value = response.data.data
    // submissionPagination.total = response.data.total
    
    // 临时模拟数据
    submissionsList.value = []
    submissionPagination.total = 0
  } catch (error) {
    console.error('加载提交数据失败:', error)
    ElMessage.error('加载提交数据失败')
  } finally {
    submissionLoading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  loadStats()
  loadPopups()
  loadSubmissions()
  ElMessage.success('数据已刷新')
}

// 弹窗管理事件处理
const handlePopupSearch = (searchForm: any) => {
  console.log('搜索弹窗:', searchForm)
  loadPopups()
}

const resetPopupSearch = () => {
  console.log('重置弹窗搜索')
  loadPopups()
}

const batchEnablePopups = (ids: number[]) => {
  console.log('批量启用弹窗:', ids)
  ElMessage.success('批量启用成功')
  loadPopups()
}

const batchDisablePopups = (ids: number[]) => {
  console.log('批量禁用弹窗:', ids)
  ElMessage.success('批量禁用成功')
  loadPopups()
}

const batchDeletePopups = (ids: number[]) => {
  console.log('批量删除弹窗:', ids)
  ElMessage.success('批量删除成功')
  loadPopups()
}

const previewPopup = (popup: any) => {
  console.log('预览弹窗:', popup)
  // 实现预览逻辑
}

const editPopup = (popup: any) => {
  editingPopup.value = popup
  showCreatePopupDialog.value = true
}

const deletePopup = (popup: any) => {
  console.log('删除弹窗:', popup)
  ElMessage.success('删除成功')
  loadPopups()
}

const handlePopupSizeChange = (size: number) => {
  popupPagination.size = size
  loadPopups()
}

const handlePopupCurrentChange = (page: number) => {
  popupPagination.page = page
  loadPopups()
}

// 提交数据事件处理
const handleSubmissionSearch = (searchForm: any) => {
  console.log('搜索提交数据:', searchForm)
  loadSubmissions()
}

const resetSubmissionSearch = () => {
  console.log('重置提交数据搜索')
  loadSubmissions()
}

const exportSubmissions = () => {
  console.log('导出提交数据')
  ElMessage.success('导出功能开发中')
}

const batchProcessSubmissions = (ids: number[]) => {
  console.log('批量处理提交数据:', ids)
  ElMessage.success('批量处理成功')
  loadSubmissions()
}

const batchDeleteSubmissions = (ids: number[]) => {
  console.log('批量删除提交数据:', ids)
  ElMessage.success('批量删除成功')
  loadSubmissions()
}

const viewSubmissionDetail = (submission: any) => {
  selectedSubmission.value = submission
  showSubmissionDetailDialog.value = true
}

const updateSubmissionStatus = (submission: any) => {
  console.log('更新提交状态:', submission)
  ElMessage.success('状态更新成功')
  loadSubmissions()
}

const deleteSubmission = (submission: any) => {
  console.log('删除提交数据:', submission)
  ElMessage.success('删除成功')
  loadSubmissions()
}

const handleSubmissionSizeChange = (size: number) => {
  submissionPagination.size = size
  loadSubmissions()
}

const handleSubmissionCurrentChange = (page: number) => {
  submissionPagination.page = page
  loadSubmissions()
}

// 弹窗表单处理
const handlePopupSubmit = (formData: any) => {
  console.log('提交弹窗表单:', formData)
  ElMessage.success(editingPopup.value ? '编辑成功' : '创建成功')
  showCreatePopupDialog.value = false
  resetPopupForm()
  loadPopups()
}

const resetPopupForm = () => {
  editingPopup.value = null
}
</script>

<style scoped>
.data-collection {
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

.main-tabs {
  margin-top: 20px;
}

/* 深色主题适配 */
.dark .data-collection {
  background: var(--el-bg-color-page);
}
</style>