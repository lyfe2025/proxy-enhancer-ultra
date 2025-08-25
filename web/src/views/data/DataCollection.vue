<template>
  <div class="data-collection">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>数据收集</h2>
        <p>管理弹窗表单和用户提交数据</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreatePopupDialog">
          <el-icon><Plus /></el-icon>
          新增弹窗
        </el-button>
        <el-button @click="exportData" :disabled="selectedSubmissions.length === 0">
          <el-icon><Download /></el-icon>
          导出数据
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><Document /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.total_popups }}</div>
          <div class="stat-label">弹窗总数</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon active">
          <el-icon><View /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.active_popups }}</div>
          <div class="stat-label">活跃弹窗</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon submissions">
          <el-icon><Edit /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.total_submissions }}</div>
          <div class="stat-label">提交总数</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon today">
          <el-icon><Calendar /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.today_submissions }}</div>
          <div class="stat-label">今日提交</div>
        </div>
      </div>
    </div>

    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="data-tabs">
      <!-- 弹窗管理 -->
      <el-tab-pane label="弹窗管理" name="popups">
        <!-- 弹窗筛选 -->
        <div class="filter-bar">
          <div class="filter-left">
            <el-input
              v-model="popupSearchQuery"
              placeholder="搜索弹窗名称或URL"
              style="width: 300px"
              clearable
              @input="handlePopupSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            
            <el-select
              v-model="popupStatusFilter"
              placeholder="状态筛选"
              style="width: 120px"
              clearable
              @change="handlePopupFilter"
            >
              <el-option label="启用" value="enabled" />
              <el-option label="禁用" value="disabled" />
            </el-select>
          </div>
          
          <div class="filter-right">
            <el-button @click="refreshPopups">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>

        <!-- 弹窗列表 -->
        <div class="popups-table">
          <el-table
            v-loading="popupsLoading"
            :data="filteredPopups"
            row-key="id"
            stripe
            style="width: 100%"
          >
            <el-table-column prop="name" label="弹窗名称" min-width="150">
              <template #default="{ row }">
                <div class="popup-name">
                  <span class="name-text">{{ row.name }}</span>
                  <el-icon v-if="!row.enabled" color="#666"><Lock /></el-icon>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="trigger_url" label="触发URL" min-width="200">
              <template #default="{ row }">
                <span class="trigger-url">{{ row.trigger_url }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="trigger_condition" label="触发条件" width="120">
              <template #default="{ row }">
                <el-tag size="small">{{ getTriggerConditionText(row.trigger_condition) }}</el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="submission_count" label="提交次数" width="100">
              <template #default="{ row }">
                <span class="submission-count">{{ row.submission_count || 0 }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="enabled" label="状态" width="80">
              <template #default="{ row }">
                <el-switch
                  v-model="row.enabled"
                  @change="togglePopup(row)"
                  :loading="row.toggling"
                  active-color="#00ff88"
                  inactive-color="#666"
                />
              </template>
            </el-table-column>
            
            <el-table-column prop="updated_at" label="更新时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.updated_at) }}
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="text" size="small" @click="previewPopup(row)">
                    <el-icon><View /></el-icon>
                    预览
                  </el-button>
                  
                  <el-button type="text" size="small" @click="editPopup(row)">
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-button>
                  
                  <el-button type="text" size="small" @click="viewSubmissions(row)">
                    <el-icon><Document /></el-icon>
                    数据
                  </el-button>
                  
                  <el-button
                    type="text"
                    size="small"
                    @click="deletePopup(row)"
                    style="color: #f56c6c"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- 提交数据 -->
      <el-tab-pane label="提交数据" name="submissions">
        <!-- 数据筛选 -->
        <div class="filter-bar">
          <div class="filter-left">
            <el-input
              v-model="submissionSearchQuery"
              placeholder="搜索提交内容"
              style="width: 300px"
              clearable
              @input="handleSubmissionSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            
            <el-select
              v-model="submissionPopupFilter"
              placeholder="弹窗筛选"
              style="width: 200px"
              clearable
              @change="handleSubmissionFilter"
            >
              <el-option
                v-for="popup in popups"
                :key="popup.id"
                :label="popup.name"
                :value="popup.id"
              />
            </el-select>
            
            <el-date-picker
              v-model="submissionDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 240px"
              @change="handleSubmissionFilter"
            />
          </div>
          
          <div class="filter-right">
            <el-button @click="refreshSubmissions">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button @click="clearSelectedSubmissions" :disabled="selectedSubmissions.length === 0">
              清除选择
            </el-button>
          </div>
        </div>

        <!-- 提交数据列表 -->
        <div class="submissions-table">
          <el-table
            v-loading="submissionsLoading"
            :data="filteredSubmissions"
            @selection-change="handleSubmissionSelectionChange"
            row-key="id"
            stripe
            style="width: 100%"
          >
            <el-table-column type="selection" width="55" />
            
            <el-table-column prop="popup_name" label="弹窗名称" width="150">
              <template #default="{ row }">
                <el-tag size="small">{{ row.popup_name }}</el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="user_ip" label="用户IP" width="120">
              <template #default="{ row }">
                <span class="user-ip">{{ row.user_ip }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="user_agent" label="用户代理" min-width="200">
              <template #default="{ row }">
                <el-tooltip :content="row.user_agent" placement="top">
                  <span class="user-agent">{{ truncateText(row.user_agent, 30) }}</span>
                </el-tooltip>
              </template>
            </el-table-column>
            
            <el-table-column prop="form_data" label="提交数据" min-width="250">
              <template #default="{ row }">
                <div class="form-data">
                  <div v-for="(value, key) in parseFormData(row.form_data)" :key="key" class="data-item">
                    <span class="data-key">{{ key }}:</span>
                    <span class="data-value">{{ truncateText(value, 20) }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="submitted_at" label="提交时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.submitted_at) }}
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="text" size="small" @click="viewSubmissionDetail(row)">
                    <el-icon><View /></el-icon>
                    详情
                  </el-button>
                  
                  <el-button
                    type="text"
                    size="small"
                    @click="deleteSubmission(row)"
                    style="color: #f56c6c"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
          
          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="submissionCurrentPage"
              v-model:page-size="submissionPageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="submissionTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSubmissionSizeChange"
              @current-change="handleSubmissionCurrentChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 创建/编辑弹窗对话框 -->
    <el-dialog
      v-model="popupDialogVisible"
      :title="isEditPopup ? '编辑弹窗' : '新增弹窗'"
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="popupFormRef"
        :model="popupForm"
        :rules="popupFormRules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="弹窗名称" prop="name">
              <el-input v-model="popupForm.name" placeholder="请输入弹窗名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="触发条件" prop="trigger_condition">
              <el-select v-model="popupForm.trigger_condition" style="width: 100%">
                <el-option label="页面加载" value="page_load" />
                <el-option label="点击事件" value="click" />
                <el-option label="滚动事件" value="scroll" />
                <el-option label="时间延迟" value="delay" />
                <el-option label="退出意图" value="exit_intent" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="触发URL" prop="trigger_url">
          <el-input
            v-model="popupForm.trigger_url"
            placeholder="例如: https://example.com/* 或 *.example.com"
          />
          <div class="form-tip">支持通配符 * 匹配多个页面</div>
        </el-form-item>
        
        <el-form-item label="弹窗标题" prop="title">
          <el-input v-model="popupForm.title" placeholder="请输入弹窗标题" />
        </el-form-item>
        
        <el-form-item label="弹窗内容" prop="content">
          <el-input
            v-model="popupForm.content"
            type="textarea"
            :rows="3"
            placeholder="请输入弹窗内容描述"
          />
        </el-form-item>
        
        <el-form-item label="表单字段" prop="form_fields">
          <div class="form-fields-editor">
            <div v-for="(field, index) in popupForm.form_fields" :key="index" class="field-item">
              <el-row :gutter="10">
                <el-col :span="6">
                  <el-input v-model="field.name" placeholder="字段名" />
                </el-col>
                <el-col :span="4">
                  <el-select v-model="field.type" placeholder="类型">
                    <el-option label="文本" value="text" />
                    <el-option label="邮箱" value="email" />
                    <el-option label="电话" value="phone" />
                    <el-option label="数字" value="number" />
                    <el-option label="多行文本" value="textarea" />
                    <el-option label="选择框" value="select" />
                    <el-option label="复选框" value="checkbox" />
                  </el-select>
                </el-col>
                <el-col :span="6">
                  <el-input v-model="field.label" placeholder="显示标签" />
                </el-col>
                <el-col :span="6">
                  <el-input v-model="field.placeholder" placeholder="提示文本" />
                </el-col>
                <el-col :span="2">
                  <el-button
                    type="text"
                    @click="removeFormField(index)"
                    style="color: #f56c6c"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-col>
              </el-row>
              
              <el-row v-if="field.type === 'select'" :gutter="10" style="margin-top: 8px">
                <el-col :span="22">
                  <el-input
                    v-model="field.options"
                    placeholder="选项值，用逗号分隔"
                  />
                </el-col>
              </el-row>
              
              <el-row :gutter="10" style="margin-top: 8px">
                <el-col :span="12">
                  <el-checkbox v-model="field.required">必填</el-checkbox>
                </el-col>
              </el-row>
            </div>
            
            <el-button @click="addFormField" type="dashed" style="width: 100%; margin-top: 10px">
              <el-icon><Plus /></el-icon>
              添加字段
            </el-button>
          </div>
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="显示延迟">
              <el-input-number
                v-model="popupForm.delay_seconds"
                :min="0"
                :max="60"
                style="width: 100%"
              />
              <div class="form-tip">秒，0表示立即显示</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item>
              <el-checkbox v-model="popupForm.enabled">启用弹窗</el-checkbox>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="popupDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPopupForm" :loading="popupSubmitting">
            {{ isEditPopup ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 弹窗预览对话框 -->
    <el-dialog
      v-model="previewDialogVisible"
      title="弹窗预览"
      width="600px"
    >
      <div class="popup-preview">
        <div class="preview-popup">
          <div class="popup-header">
            <h3>{{ previewPopup.title }}</h3>
            <el-button type="text" @click="previewDialogVisible = false">
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
          
          <div class="popup-content">
            <p>{{ previewPopup.content }}</p>
            
            <div class="popup-form">
              <div v-for="field in previewPopup.form_fields" :key="field.name" class="form-field">
                <label>{{ field.label }}</label>
                <el-input
                  v-if="field.type === 'text' || field.type === 'email' || field.type === 'phone'"
                  :placeholder="field.placeholder"
                  disabled
                />
                <el-input
                  v-else-if="field.type === 'textarea'"
                  type="textarea"
                  :placeholder="field.placeholder"
                  disabled
                />
                <el-select
                  v-else-if="field.type === 'select'"
                  :placeholder="field.placeholder"
                  disabled
                >
                  <el-option
                    v-for="option in field.options?.split(',')"
                    :key="option"
                    :label="option.trim()"
                    :value="option.trim()"
                  />
                </el-select>
                <el-checkbox v-else-if="field.type === 'checkbox'" disabled>
                  {{ field.label }}
                </el-checkbox>
              </div>
              
              <div class="form-actions">
                <el-button type="primary" disabled>提交</el-button>
                <el-button disabled>取消</el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 提交详情对话框 -->
    <el-dialog
      v-model="submissionDetailDialogVisible"
      title="提交详情"
      width="700px"
    >
      <div v-if="selectedSubmissionDetail" class="submission-detail">
        <div class="detail-section">
          <h4>基本信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="弹窗名称">{{ selectedSubmissionDetail.popup_name }}</el-descriptions-item>
            <el-descriptions-item label="提交时间">{{ formatDate(selectedSubmissionDetail.submitted_at) }}</el-descriptions-item>
            <el-descriptions-item label="用户IP">{{ selectedSubmissionDetail.user_ip }}</el-descriptions-item>
            <el-descriptions-item label="用户代理">{{ selectedSubmissionDetail.user_agent }}</el-descriptions-item>
          </el-descriptions>
        </div>
        
        <div class="detail-section">
          <h4>提交数据</h4>
          <div class="form-data-detail">
            <div v-for="(value, key) in parseFormData(selectedSubmissionDetail.form_data)" :key="key" class="data-row">
              <div class="data-label">{{ key }}:</div>
              <div class="data-content">{{ value }}</div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Delete,
  Edit,
  View,
  Document,
  Calendar,
  Lock,
  Download,
  Upload,
  Close
} from '@element-plus/icons-vue'
import { popupsApi, submissionsApi } from '@/api'
import type { Popup, Submission } from '@/types'

// 响应式数据
const activeTab = ref('popups')
const popupsLoading = ref(false)
const submissionsLoading = ref(false)
const popupSubmitting = ref(false)
const popupDialogVisible = ref(false)
const previewDialogVisible = ref(false)
const submissionDetailDialogVisible = ref(false)
const isEditPopup = ref(false)

// 搜索和筛选
const popupSearchQuery = ref('')
const popupStatusFilter = ref('')
const submissionSearchQuery = ref('')
const submissionPopupFilter = ref('')
const submissionDateRange = ref<[Date, Date] | null>(null)

// 分页
const submissionCurrentPage = ref(1)
const submissionPageSize = ref(20)
const submissionTotal = ref(0)

// 数据
const popups = ref<Popup[]>([])
const submissions = ref<Submission[]>([])
const selectedSubmissions = ref<Submission[]>([])
const selectedSubmissionDetail = ref<Submission | null>(null)
const previewPopup = ref<Popup>({} as Popup)

// 统计数据
const stats = ref({
  total_popups: 0,
  active_popups: 0,
  total_submissions: 0,
  today_submissions: 0
})

// 表单引用和数据
const popupFormRef = ref()
const popupForm = reactive({
  id: '',
  name: '',
  trigger_url: '',
  trigger_condition: 'page_load',
  title: '',
  content: '',
  form_fields: [] as any[],
  delay_seconds: 0,
  enabled: true
})

// 表单验证规则
const popupFormRules = {
  name: [
    { required: true, message: '请输入弹窗名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  trigger_url: [
    { required: true, message: '请输入触发URL', trigger: 'blur' }
  ],
  trigger_condition: [
    { required: true, message: '请选择触发条件', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入弹窗标题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入弹窗内容', trigger: 'blur' }
  ]
}

// 计算属性
const filteredPopups = computed(() => {
  let result = popups.value
  
  // 搜索过滤
  if (popupSearchQuery.value) {
    const query = popupSearchQuery.value.toLowerCase()
    result = result.filter(popup => 
      popup.name.toLowerCase().includes(query) ||
      popup.trigger_url.toLowerCase().includes(query)
    )
  }
  
  // 状态过滤
  if (popupStatusFilter.value) {
    const enabled = popupStatusFilter.value === 'enabled'
    result = result.filter(popup => popup.enabled === enabled)
  }
  
  return result
})

const filteredSubmissions = computed(() => {
  let result = submissions.value
  
  // 搜索过滤
  if (submissionSearchQuery.value) {
    const query = submissionSearchQuery.value.toLowerCase()
    result = result.filter(submission => {
      const formDataStr = JSON.stringify(submission.form_data).toLowerCase()
      return formDataStr.includes(query)
    })
  }
  
  // 弹窗过滤
  if (submissionPopupFilter.value) {
    result = result.filter(submission => submission.popup_id === submissionPopupFilter.value)
  }
  
  // 日期过滤
  if (submissionDateRange.value) {
    const [startDate, endDate] = submissionDateRange.value
    result = result.filter(submission => {
      const submissionDate = new Date(submission.submitted_at)
      return submissionDate >= startDate && submissionDate <= endDate
    })
  }
  
  return result.sort((a, b) => new Date(b.submitted_at).getTime() - new Date(a.submitted_at).getTime())
})

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await popupsApi.getStats()
    stats.value = response.data
  } catch (error: any) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取弹窗列表
const fetchPopups = async () => {
  try {
    popupsLoading.value = true
    const response = await popupsApi.getPopups()
    
    popups.value = response.data.map(popup => ({
      ...popup,
      toggling: false
    }))
    
  } catch (error: any) {
    ElMessage.error(error.message || '获取弹窗列表失败')
  } finally {
    popupsLoading.value = false
  }
}

// 获取提交数据列表
const fetchSubmissions = async () => {
  try {
    submissionsLoading.value = true
    const response = await submissionsApi.getSubmissions({
      page: submissionCurrentPage.value,
      limit: submissionPageSize.value
    })
    
    submissions.value = response.data.items
    submissionTotal.value = response.data.total
    
  } catch (error: any) {
    ElMessage.error(error.message || '获取提交数据失败')
  } finally {
    submissionsLoading.value = false
  }
}

// 搜索处理
const handlePopupSearch = () => {
  // 实时搜索，无需额外处理
}

const handleSubmissionSearch = () => {
  submissionCurrentPage.value = 1
}

// 筛选处理
const handlePopupFilter = () => {
  // 实时筛选，无需额外处理
}

const handleSubmissionFilter = () => {
  submissionCurrentPage.value = 1
}

// 刷新数据
const refreshPopups = () => {
  fetchPopups()
  fetchStats()
}

const refreshSubmissions = () => {
  fetchSubmissions()
}

// 分页处理
const handleSubmissionSizeChange = (size: number) => {
  submissionPageSize.value = size
  fetchSubmissions()
}

const handleSubmissionCurrentChange = (page: number) => {
  submissionCurrentPage.value = page
  fetchSubmissions()
}

// 选择变化处理
const handleSubmissionSelectionChange = (selection: Submission[]) => {
  selectedSubmissions.value = selection
}

const clearSelectedSubmissions = () => {
  selectedSubmissions.value = []
}

// 显示创建弹窗对话框
const showCreatePopupDialog = () => {
  isEditPopup.value = false
  resetPopupForm()
  popupDialogVisible.value = true
}

// 编辑弹窗
const editPopup = (popup: Popup) => {
  isEditPopup.value = true
  Object.assign(popupForm, {
    id: popup.id,
    name: popup.name,
    trigger_url: popup.trigger_url,
    trigger_condition: popup.trigger_condition,
    title: popup.title,
    content: popup.content,
    form_fields: popup.form_fields || [],
    delay_seconds: popup.delay_seconds || 0,
    enabled: popup.enabled
  })
  popupDialogVisible.value = true
}

// 预览弹窗
const previewPopup = (popup: Popup) => {
  previewPopup.value = popup
  previewDialogVisible.value = true
}

// 查看提交数据
const viewSubmissions = (popup: Popup) => {
  activeTab.value = 'submissions'
  submissionPopupFilter.value = popup.id
  handleSubmissionFilter()
}

// 查看提交详情
const viewSubmissionDetail = (submission: Submission) => {
  selectedSubmissionDetail.value = submission
  submissionDetailDialogVisible.value = true
}

// 重置弹窗表单
const resetPopupForm = () => {
  Object.assign(popupForm, {
    id: '',
    name: '',
    trigger_url: '',
    trigger_condition: 'page_load',
    title: '',
    content: '',
    form_fields: [],
    delay_seconds: 0,
    enabled: true
  })
  
  if (popupFormRef.value) {
    popupFormRef.value.clearValidate()
  }
}

// 提交弹窗表单
const submitPopupForm = async () => {
  if (!popupFormRef.value) return
  
  try {
    await popupFormRef.value.validate()
    popupSubmitting.value = true
    
    const data = { ...popupForm }
    delete data.id
    
    if (isEditPopup.value) {
      await popupsApi.updatePopup(popupForm.id, data)
      ElMessage.success('弹窗更新成功')
    } else {
      await popupsApi.createPopup(data)
      ElMessage.success('弹窗创建成功')
    }
    
    popupDialogVisible.value = false
    fetchPopups()
    fetchStats()
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    popupSubmitting.value = false
  }
}

// 切换弹窗状态
const togglePopup = async (popup: Popup) => {
  try {
    popup.toggling = true
    
    if (popup.enabled) {
      await popupsApi.enablePopup(popup.id)
      ElMessage.success('弹窗已启用')
    } else {
      await popupsApi.disablePopup(popup.id)
      ElMessage.success('弹窗已禁用')
    }
    
    fetchStats()
    
  } catch (error: any) {
    popup.enabled = !popup.enabled // 回滚状态
    ElMessage.error(error.message || '操作失败')
  } finally {
    popup.toggling = false
  }
}

// 删除弹窗
const deletePopup = async (popup: Popup) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除弹窗 "${popup.name}" 吗？此操作将同时删除相关的提交数据。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await popupsApi.deletePopup(popup.id)
    ElMessage.success('弹窗删除成功')
    fetchPopups()
    fetchStats()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 删除提交数据
const deleteSubmission = async (submission: Submission) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这条提交数据吗？此操作不可恢复。',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await submissionsApi.deleteSubmission(submission.id)
    ElMessage.success('提交数据删除成功')
    fetchSubmissions()
    fetchStats()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 导出数据
const exportData = async () => {
  if (selectedSubmissions.value.length === 0) {
    ElMessage.warning('请选择要导出的数据')
    return
  }
  
  try {
    const data = selectedSubmissions.value.map(submission => ({
      弹窗名称: submission.popup_name,
      用户IP: submission.user_ip,
      用户代理: submission.user_agent,
      提交数据: JSON.stringify(submission.form_data),
      提交时间: formatDate(submission.submitted_at)
    }))
    
    const csvContent = convertToCSV(data)
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `submissions_${new Date().toISOString().split('T')[0]}.csv`
    link.click()
    
    URL.revokeObjectURL(url)
    ElMessage.success('数据导出成功')
    
  } catch (error: any) {
    ElMessage.error(error.message || '导出失败')
  }
}

// 表单字段管理
const addFormField = () => {
  popupForm.form_fields.push({
    name: '',
    type: 'text',
    label: '',
    placeholder: '',
    required: false,
    options: ''
  })
}

const removeFormField = (index: number) => {
  popupForm.form_fields.splice(index, 1)
}

// 工具函数
const getTriggerConditionText = (condition: string) => {
  const conditionMap: Record<string, string> = {
    page_load: '页面加载',
    click: '点击事件',
    scroll: '滚动事件',
    delay: '时间延迟',
    exit_intent: '退出意图'
  }
  return conditionMap[condition] || condition
}

const parseFormData = (formData: any) => {
  if (typeof formData === 'string') {
    try {
      return JSON.parse(formData)
    } catch {
      return { data: formData }
    }
  }
  return formData || {}
}

const truncateText = (text: string, maxLength: number) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

const convertToCSV = (data: any[]) => {
  if (data.length === 0) return ''
  
  const headers = Object.keys(data[0])
  const csvRows = []
  
  // 添加标题行
  csvRows.push(headers.join(','))
  
  // 添加数据行
  for (const row of data) {
    const values = headers.map(header => {
      const value = row[header]
      return `"${String(value).replace(/"/g, '""')}"`
    })
    csvRows.push(values.join(','))
  }
  
  return csvRows.join('\n')
}

// 组件挂载
onMounted(() => {
  fetchStats()
  fetchPopups()
  fetchSubmissions()
})
</script>

<style scoped>
.data-collection {
  padding: 0;
}

/* 页面头部 */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding: 24px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%);
  border: 1px solid #333;
  border-radius: 12px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  color: #fff;
  font-size: 24px;
  font-weight: 600;
}

.header-left p {
  margin: 0;
  color: #888;
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%);
  border: 1px solid #333;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  border-color: #00ff88;
  box-shadow: 0 4px 12px rgba(0, 255, 136, 0.1);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background-color: #333;
  border-radius: 8px;
  margin-right: 16px;
  color: #666;
}

.stat-icon.active {
  background-color: rgba(0, 255, 136, 0.1);
  color: #00ff88;
}

.stat-icon.submissions {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.stat-icon.today {
  background-color: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #888;
}

/* 标签页 */
.data-tabs {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  overflow: hidden;
}

/* 筛选栏 */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 16px 20px;
  background-color: #2a2a2a;
  border-bottom: 1px solid #333;
}

.filter-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 表格 */
.popups-table,
.submissions-table {
  padding: 0 20px 20px;
}

.popup-name,
.submission-detail {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name-text {
  font-weight: 500;
  color: #fff;
}

.trigger-url {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  color: #00ff88;
  font-size: 13px;
}

.submission-count {
  color: #ccc;
  font-weight: 500;
}

.user-ip {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  color: #409eff;
  font-size: 13px;
}

.user-agent {
  color: #ccc;
  font-size: 12px;
}

.form-data {
  max-width: 250px;
}

.data-item {
  display: flex;
  margin-bottom: 4px;
  font-size: 12px;
}

.data-key {
  color: #888;
  margin-right: 4px;
  min-width: 60px;
}

.data-value {
  color: #ccc;
  flex: 1;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  padding: 20px;
  border-top: 1px solid #333;
}

/* 对话框 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.form-tip {
  font-size: 12px;
  color: #888;
  margin-top: 4px;
}

/* 表单字段编辑器 */
.form-fields-editor {
  border: 1px solid #333;
  border-radius: 8px;
  padding: 16px;
  background-color: #2a2a2a;
}

.field-item {
  margin-bottom: 16px;
  padding: 12px;
  background-color: #1a1a1a;
  border-radius: 6px;
  border: 1px solid #333;
}

.field-item:last-child {
  margin-bottom: 0;
}

/* 弹窗预览 */
.popup-preview {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.preview-popup {
  width: 400px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
}

.popup-header h3 {
  margin: 0;
  color: #333;
  font-size: 16px;
}

.popup-content {
  padding: 20px;
}

.popup-content p {
  margin: 0 0 20px 0;
  color: #666;
  line-height: 1.5;
}

.popup-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-field label {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

/* 提交详情 */
.detail-section {
  margin-bottom: 24px;
}

.detail-section h4 {
  margin: 0 0 12px 0;
  color: #fff;
  font-size: 16px;
}

.form-data-detail {
  background-color: #2a2a2a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 16px;
}

.data-row {
  display: flex;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #333;
}

.data-row:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.data-label {
  color: #888;
  min-width: 100px;
  font-weight: 500;
}

.data-content {
  color: #ccc;
  flex: 1;
  word-break: break-all;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .header-right {
    width: 100%;
    justify-content: flex-start;
  }
  
  .stats-cards {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }
  
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .filter-left,
  .filter-right {
    justify-content: flex-start;
    flex-wrap: wrap;
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: stretch;
  }
  
  .preview-popup {
    width: 100%;
    max-width: 350px;
  }
}

/* Element Plus 组件样式覆盖 */
:deep(.el-tabs__header) {
  background-color: #2a2a2a;
  margin: 0;
  border-bottom: 1px solid #333;
}

:deep(.el-tabs__nav-wrap) {
  padding: 0 20px;
}

:deep(.el-tabs__item) {
  color: #888;
  border-bottom: 2px solid transparent;
}

:deep(.el-tabs__item.is-active) {
  color: #00ff88;
  border-bottom-color: #00ff88;
}

:deep(.el-tabs__item:hover) {
  color: #00ff88;
}

:deep(.el-tabs__content) {
  padding: 0;
}

:deep(.el-table) {
  background-color: transparent;
  color: #ccc;
}

:deep(.el-table__header) {
  background-color: #2a2a2a;
}

:deep(.el-table th) {
  background-color: #2a2a2a;
  color: #fff;
  border-bottom: 1px solid #333;
}

:deep(.el-table td) {
  border-bottom: 1px solid #333;
}

:deep(.el-table__row) {
  background-color: transparent;
}

:deep(.el-table__row:hover) {
  background-color: rgba(0, 255, 136, 0.05);
}

:deep(.el-table__row--striped) {
  background-color: rgba(255, 255, 255, 0.02);
}

:deep(.el-table__row--striped:hover) {
  background-color: rgba(0, 255, 136, 0.05);
}

:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #00ff88;
}

:deep(.el-button--text) {
  color: #ccc;
}

:deep(.el-button--text:hover) {
  color: #00ff88;
}

:deep(.el-input__wrapper) {
  background-color: #2a2a2a;
  border: 1px solid #333;
}

:deep(.el-input__wrapper:hover) {
  border-color: #00ff88;
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #00ff88;
  box-shadow: 0 0 0 2px rgba(0, 255, 136, 0.2);
}

:deep(.el-input__inner) {
  color: #ccc;
}

:deep(.el-select) {
  --el-select-input-color: #ccc;
  --el-select-border-color-hover: #00ff88;
}

:deep(.el-pagination) {
  --el-pagination-text-color: #ccc;
  --el-pagination-bg-color: #2a2a2a;
  --el-pagination-border-color: #333;
}

:deep(.el-dialog) {
  background-color: #1a1a1a;
  border: 1px solid #333;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #333;
}

:deep(.el-dialog__title) {
  color: #fff;
}

:deep(.el-form-item__label) {
  color: #ccc;
}

:deep(.el-checkbox__label) {
  color: #ccc;
}

:deep(.el-radio__label) {
  color: #ccc;
}

:deep(.el-textarea__inner) {
  background-color: #2a2a2a;
  border: 1px solid #333;
  color: #ccc;
}

:deep(.el-textarea__inner:hover) {
  border-color: #00ff88;
}

:deep(.el-textarea__inner:focus) {
  border-color: #00ff88;
  box-shadow: 0 0 0 2px rgba(0, 255, 136, 0.2);
}

:deep(.el-descriptions) {
  --el-descriptions-table-border: 1px solid #333;
  --el-descriptions-item-bordered-label-background: #2a2a2a;
}

:deep(.el-descriptions__label) {
  color: #888;
}

:deep(.el-descriptions__content) {
  color: #ccc;
}

:deep(.el-date-editor) {
  --el-date-editor-width: 240px;
}

:deep(.el-date-editor .el-input__wrapper) {
  background-color: #2a2a2a;
  border: 1px solid #333;
}

:deep(.el-button--dashed) {
  border: 1px dashed #333;
  color: #888;
  background-color: transparent;
}

:deep(.el-button--dashed:hover) {
  border-color: #00ff88;
  color: #00ff88;
}
</style>