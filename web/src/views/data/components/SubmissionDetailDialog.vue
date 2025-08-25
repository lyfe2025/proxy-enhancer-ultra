<template>
  <el-dialog
    v-model="dialogVisible"
    title="提交详情"
    width="600px"
    @close="handleClose"
  >
    <div class="submission-detail" v-if="submission">
      <!-- 基本信息 -->
      <el-descriptions title="基本信息" :column="2" border>
        <el-descriptions-item label="弹窗名称">
          {{ submission.popupName }}
        </el-descriptions-item>
        <el-descriptions-item label="提交状态">
          <el-tag :type="getStatusTagType(submission.status)" size="small">
            {{ getStatusText(submission.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="提交时间">
          {{ formatDate(submission.submittedAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          {{ submission.ip }}
        </el-descriptions-item>
        <el-descriptions-item label="用户代理" :span="2">
          <el-tooltip :content="submission.userAgent" placement="top">
            <span class="user-agent">{{ submission.userAgent }}</span>
          </el-tooltip>
        </el-descriptions-item>
        <el-descriptions-item label="来源页面" :span="2" v-if="submission.referrer">
          <a :href="submission.referrer" target="_blank" class="referrer-link">
            {{ submission.referrer }}
          </a>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 提交数据 -->
      <div class="submission-data-section">
        <h4>提交数据</h4>
        <el-table :data="dataTableData" border stripe>
          <el-table-column prop="key" label="字段名称" width="150" />
          <el-table-column prop="value" label="字段值" min-width="200">
            <template #default="{ row }">
              <span class="field-value">{{ row.value }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="字段类型" width="100">
            <template #default="{ row }">
              <el-tag size="small" type="info">{{ row.type }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 处理历史 -->
      <div class="processing-history" v-if="submission.history && submission.history.length">
        <h4>处理历史</h4>
        <el-timeline>
          <el-timeline-item
            v-for="(item, index) in submission.history"
            :key="index"
            :timestamp="formatDate(item.createdAt)"
            :type="getHistoryType(item.action)"
          >
            <div class="history-item">
              <div class="history-action">{{ item.action }}</div>
              <div class="history-operator">操作人：{{ item.operator }}</div>
              <div class="history-comment" v-if="item.comment">
                备注：{{ item.comment }}
              </div>
            </div>
          </el-timeline-item>
        </el-timeline>
      </div>

      <!-- 状态更新 -->
      <div class="status-update">
        <h4>状态处理</h4>
        <el-form :model="statusForm" inline>
          <el-form-item label="处理状态">
            <el-select v-model="statusForm.status" placeholder="选择处理状态">
              <el-option label="待处理" value="pending" />
              <el-option label="已处理" value="processed" />
              <el-option label="已忽略" value="ignored" />
            </el-select>
          </el-form-item>
          <el-form-item label="处理备注">
            <el-input
              v-model="statusForm.comment"
              placeholder="请输入处理备注"
              style="width: 200px;"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updateStatus" :loading="updateLoading">
              更新状态
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="exportData">
          <el-icon><Download /></el-icon>
          导出数据
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Download } from '@element-plus/icons-vue'

interface SubmissionData {
  id: number
  popupName: string
  data: Record<string, any>
  status: string
  ip: string
  userAgent: string
  submittedAt: string
  referrer?: string
  history?: Array<{
    action: string
    operator: string
    comment?: string
    createdAt: string
  }>
}

interface StatusForm {
  status: string
  comment: string
}

const props = defineProps<{
  submission?: SubmissionData | null
}>()

const emit = defineEmits<{
  statusUpdate: [id: number, status: string, comment: string]
  export: [id: number]
}>()

const dialogVisible = defineModel<boolean>({ required: true })
const updateLoading = ref(false)

const statusForm = reactive<StatusForm>({
  status: '',
  comment: ''
})

// 将提交数据转换为表格数据
const dataTableData = computed(() => {
  if (!props.submission?.data) return []
  
  return Object.entries(props.submission.data).map(([key, value]) => ({
    key,
    value: Array.isArray(value) ? value.join(', ') : String(value),
    type: Array.isArray(value) ? 'array' : typeof value
  }))
})

// 监听提交数据变化，更新状态表单
watch(
  () => props.submission,
  (newValue) => {
    if (newValue) {
      statusForm.status = newValue.status
      statusForm.comment = ''
    }
  },
  { immediate: true }
)

const updateStatus = async () => {
  if (!props.submission) return
  
  if (!statusForm.status) {
    ElMessage.warning('请选择处理状态')
    return
  }
  
  updateLoading.value = true
  try {
    emit('statusUpdate', props.submission.id, statusForm.status, statusForm.comment)
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('更新状态失败:', error)
    ElMessage.error('更新状态失败')
  } finally {
    updateLoading.value = false
  }
}

const exportData = () => {
  if (!props.submission) return
  emit('export', props.submission.id)
}

const handleClose = () => {
  statusForm.status = ''
  statusForm.comment = ''
}

// 工具函数
const getStatusTagType = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    processed: 'success',
    ignored: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    processed: '已处理',
    ignored: '已忽略'
  }
  return statusMap[status] || status
}

const getHistoryType = (action: string) => {
  const typeMap: Record<string, string> = {
    '创建': 'primary',
    '处理': 'success',
    '忽略': 'warning',
    '删除': 'danger'
  }
  return typeMap[action] || 'primary'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.submission-detail {
  padding: 20px 0;
}

.submission-data-section,
.processing-history,
.status-update {
  margin-top: 24px;
}

.submission-data-section h4,
.processing-history h4,
.status-update h4 {
  margin: 0 0 16px 0;
  color: var(--el-text-color-primary);
  font-size: 16px;
  font-weight: 600;
}

.user-agent {
  display: inline-block;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: top;
}

.referrer-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.referrer-link:hover {
  text-decoration: underline;
}

.field-value {
  word-break: break-all;
  line-height: 1.4;
}

.history-item {
  padding: 8px 0;
}

.history-action {
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 4px;
}

.history-operator {
  font-size: 12px;
  color: var(--el-text-color-regular);
  margin-bottom: 4px;
}

.history-comment {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  background: var(--el-bg-color-page);
  padding: 4px 8px;
  border-radius: 4px;
  border-left: 3px solid var(--el-color-primary);
}

.dialog-footer {
  text-align: right;
}

/* 深色主题适配 */
.dark .history-comment {
  background: var(--el-bg-color);
}
</style>
