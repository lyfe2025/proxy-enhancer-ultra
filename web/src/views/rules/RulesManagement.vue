<template>
  <div class="rules-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>规则配置</h2>
        <p>管理代理规则和过滤条件</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新增规则
        </el-button>
        <el-button @click="showBatchPriorityDialog" :disabled="selectedRules.length === 0">
          <el-icon><Sort /></el-icon>
          批量调整优先级
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filter-bar">
      <div class="filter-left">
        <el-input
          v-model="searchQuery"
          placeholder="搜索规则名称或条件"
          style="width: 300px"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select
          v-model="typeFilter"
          placeholder="类型筛选"
          style="width: 120px"
          clearable
          @change="handleFilter"
        >
          <el-option label="URL匹配" value="url" />
          <el-option label="域名匹配" value="domain" />
          <el-option label="IP匹配" value="ip" />
          <el-option label="用户代理" value="user_agent" />
        </el-select>
        
        <el-select
          v-model="actionFilter"
          placeholder="动作筛选"
          style="width: 120px"
          clearable
          @change="handleFilter"
        >
          <el-option label="允许" value="allow" />
          <el-option label="阻止" value="block" />
          <el-option label="重定向" value="redirect" />
          <el-option label="修改" value="modify" />
        </el-select>
        
        <el-select
          v-model="statusFilter"
          placeholder="状态筛选"
          style="width: 120px"
          clearable
          @change="handleFilter"
        >
          <el-option label="启用" value="enabled" />
          <el-option label="禁用" value="disabled" />
        </el-select>
      </div>
      
      <div class="filter-right">
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="exportRules">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
        <el-button @click="showImportDialog">
          <el-icon><Upload /></el-icon>
          导入
        </el-button>
      </div>
    </div>

    <!-- 规则列表 -->
    <div class="rules-table">
      <el-table
        v-loading="loading"
        :data="filteredRules"
        @selection-change="handleSelectionChange"
        row-key="id"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="priority" label="优先级" width="80" sortable>
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.priority }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="name" label="规则名称" min-width="150">
          <template #default="{ row }">
            <div class="rule-name">
              <span class="name-text">{{ row.name }}</span>
              <el-icon v-if="!row.enabled" color="#666"><Lock /></el-icon>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="condition" label="匹配条件" min-width="200">
          <template #default="{ row }">
            <span class="rule-condition">{{ row.condition }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="action" label="动作" width="100">
          <template #default="{ row }">
            <el-tag :type="getActionTagType(row.action)" size="small">
              {{ getActionText(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="match_count" label="匹配次数" width="100">
          <template #default="{ row }">
            <span class="match-count">{{ row.match_count || 0 }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              @change="toggleRule(row)"
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
              <el-button
                type="text"
                size="small"
                @click="testRule(row)"
                :loading="row.testing"
              >
                <el-icon><Connection /></el-icon>
                测试
              </el-button>
              
              <el-button type="text" size="small" @click="editRule(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              
              <el-button type="text" size="small" @click="duplicateRule(row)">
                <el-icon><CopyDocument /></el-icon>
                复制
              </el-button>
              
              <el-button
                type="text"
                size="small"
                @click="deleteRule(row)"
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
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 创建/编辑规则对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑规则' : '新增规则'"
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="ruleForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="ruleForm.name" placeholder="请输入规则名称" />
        </el-form-item>
        
        <el-form-item label="匹配类型" prop="type">
          <el-select v-model="ruleForm.type" style="width: 100%" @change="handleTypeChange">
            <el-option label="URL匹配" value="url" />
            <el-option label="域名匹配" value="domain" />
            <el-option label="IP匹配" value="ip" />
            <el-option label="用户代理" value="user_agent" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="匹配条件" prop="condition">
          <el-input
            v-model="ruleForm.condition"
            :placeholder="getConditionPlaceholder(ruleForm.type)"
            type="textarea"
            :rows="2"
          />
          <div class="form-tip">
            {{ getConditionTip(ruleForm.type) }}
          </div>
        </el-form-item>
        
        <el-form-item label="执行动作" prop="action">
          <el-select v-model="ruleForm.action" style="width: 100%" @change="handleActionChange">
            <el-option label="允许" value="allow" />
            <el-option label="阻止" value="block" />
            <el-option label="重定向" value="redirect" />
            <el-option label="修改" value="modify" />
          </el-select>
        </el-form-item>
        
        <el-form-item
          v-if="ruleForm.action === 'redirect'"
          label="重定向URL"
          prop="redirect_url"
        >
          <el-input v-model="ruleForm.redirect_url" placeholder="请输入重定向URL" />
        </el-form-item>
        
        <el-form-item
          v-if="ruleForm.action === 'modify'"
          label="修改配置"
          prop="modify_config"
        >
          <el-input
            v-model="ruleForm.modify_config"
            type="textarea"
            :rows="3"
            placeholder="请输入修改配置（JSON格式）"
          />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="优先级" prop="priority">
              <el-input-number
                v-model="ruleForm.priority"
                :min="1"
                :max="999"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item>
              <el-checkbox v-model="ruleForm.enabled">启用规则</el-checkbox>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="描述">
          <el-input
            v-model="ruleForm.description"
            type="textarea"
            :rows="2"
            placeholder="请输入规则描述（可选）"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 批量优先级调整对话框 -->
    <el-dialog
      v-model="priorityDialogVisible"
      title="批量调整优先级"
      width="500px"
    >
      <div class="priority-form">
        <p>已选择 {{ selectedRules.length }} 个规则</p>
        <el-form label-width="100px">
          <el-form-item label="调整方式">
            <el-radio-group v-model="priorityAdjustType">
              <el-radio label="set">设置为</el-radio>
              <el-radio label="increase">增加</el-radio>
              <el-radio label="decrease">减少</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item label="数值">
            <el-input-number
              v-model="priorityValue"
              :min="priorityAdjustType === 'set' ? 1 : -999"
              :max="999"
              style="width: 100%"
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="priorityDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPriorityAdjust" :loading="adjusting">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 导入规则对话框 -->
    <el-dialog
      v-model="importDialogVisible"
      title="导入规则"
      width="600px"
    >
      <div class="import-form">
        <el-upload
          ref="uploadRef"
          :auto-upload="false"
          :show-file-list="true"
          :limit="1"
          accept=".json"
          @change="handleFileChange"
        >
          <el-button type="primary">
            <el-icon><Upload /></el-icon>
            选择文件
          </el-button>
          <template #tip>
            <div class="el-upload__tip">
              只能上传 JSON 格式的规则文件
            </div>
          </template>
        </el-upload>
        
        <div v-if="importPreview.length > 0" class="import-preview">
          <h4>预览导入规则 ({{ importPreview.length }} 条)</h4>
          <div class="preview-list">
            <div v-for="rule in importPreview.slice(0, 5)" :key="rule.name" class="preview-item">
              <span class="rule-name">{{ rule.name }}</span>
              <el-tag size="small">{{ getTypeText(rule.type) }}</el-tag>
              <el-tag size="small" :type="getActionTagType(rule.action)">
                {{ getActionText(rule.action) }}
              </el-tag>
            </div>
            <div v-if="importPreview.length > 5" class="more-items">
              还有 {{ importPreview.length - 5 }} 条规则...
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="importDialogVisible = false">取消</el-button>
          <el-button
            type="primary"
            @click="submitImport"
            :loading="importing"
            :disabled="importPreview.length === 0"
          >
            导入 ({{ importPreview.length }})
          </el-button>
        </div>
      </template>
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
  Connection,
  Lock,
  Sort,
  Download,
  Upload,
  CopyDocument
} from '@element-plus/icons-vue'
import { rulesApi } from '@/api'
import type { Rule } from '@/types'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const adjusting = ref(false)
const importing = ref(false)
const dialogVisible = ref(false)
const priorityDialogVisible = ref(false)
const importDialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const typeFilter = ref('')
const actionFilter = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedRules = ref<Rule[]>([])
const rules = ref<Rule[]>([])
const importPreview = ref<Rule[]>([])

// 优先级调整
const priorityAdjustType = ref('set')
const priorityValue = ref(1)

// 表单引用和数据
const formRef = ref()
const uploadRef = ref()
const ruleForm = reactive({
  id: '',
  name: '',
  type: 'url',
  condition: '',
  action: 'allow',
  redirect_url: '',
  modify_config: '',
  priority: 100,
  enabled: true,
  description: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入规则名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择匹配类型', trigger: 'change' }
  ],
  condition: [
    { required: true, message: '请输入匹配条件', trigger: 'blur' }
  ],
  action: [
    { required: true, message: '请选择执行动作', trigger: 'change' }
  ],
  redirect_url: [
    { required: true, message: '请输入重定向URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ],
  priority: [
    { required: true, message: '请输入优先级', trigger: 'blur' },
    { type: 'number', min: 1, max: 999, message: '优先级范围 1-999', trigger: 'blur' }
  ]
}

// 计算属性
const filteredRules = computed(() => {
  let result = rules.value
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(rule => 
      rule.name.toLowerCase().includes(query) ||
      rule.condition.toLowerCase().includes(query)
    )
  }
  
  // 类型过滤
  if (typeFilter.value) {
    result = result.filter(rule => rule.type === typeFilter.value)
  }
  
  // 动作过滤
  if (actionFilter.value) {
    result = result.filter(rule => rule.action === actionFilter.value)
  }
  
  // 状态过滤
  if (statusFilter.value) {
    const enabled = statusFilter.value === 'enabled'
    result = result.filter(rule => rule.enabled === enabled)
  }
  
  return result.sort((a, b) => a.priority - b.priority)
})

// 获取规则列表
const fetchRules = async () => {
  try {
    loading.value = true
    const response = await rulesApi.getRules({
      page: currentPage.value,
      limit: pageSize.value
    })
    
    rules.value = response.data.items.map(rule => ({
      ...rule,
      testing: false,
      toggling: false
    }))
    total.value = response.data.total
    
  } catch (error: any) {
    ElMessage.error(error.message || '获取规则列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
}

// 筛选处理
const handleFilter = () => {
  currentPage.value = 1
}

// 刷新数据
const refreshData = () => {
  fetchRules()
}

// 选择变化处理
const handleSelectionChange = (selection: Rule[]) => {
  selectedRules.value = selection
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchRules()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchRules()
}

// 显示创建对话框
const showCreateDialog = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 编辑规则
const editRule = (rule: Rule) => {
  isEdit.value = true
  Object.assign(ruleForm, {
    id: rule.id,
    name: rule.name,
    type: rule.type,
    condition: rule.condition,
    action: rule.action,
    redirect_url: rule.redirect_url || '',
    modify_config: rule.modify_config || '',
    priority: rule.priority,
    enabled: rule.enabled,
    description: rule.description || ''
  })
  dialogVisible.value = true
}

// 复制规则
const duplicateRule = (rule: Rule) => {
  isEdit.value = false
  Object.assign(ruleForm, {
    id: '',
    name: `${rule.name} (副本)`,
    type: rule.type,
    condition: rule.condition,
    action: rule.action,
    redirect_url: rule.redirect_url || '',
    modify_config: rule.modify_config || '',
    priority: rule.priority + 1,
    enabled: false,
    description: rule.description || ''
  })
  dialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  Object.assign(ruleForm, {
    id: '',
    name: '',
    type: 'url',
    condition: '',
    action: 'allow',
    redirect_url: '',
    modify_config: '',
    priority: 100,
    enabled: true,
    description: ''
  })
  
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    const data = { ...ruleForm }
    delete data.id
    
    if (isEdit.value) {
      await rulesApi.updateRule(ruleForm.id, data)
      ElMessage.success('规则更新成功')
    } else {
      await rulesApi.createRule(data)
      ElMessage.success('规则创建成功')
    }
    
    dialogVisible.value = false
    fetchRules()
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 测试规则
const testRule = async (rule: Rule) => {
  try {
    rule.testing = true
    const response = await rulesApi.testRule(rule.id, {
      test_url: 'https://example.com'
    })
    
    if (response.data.matched) {
      ElMessage.success(`规则匹配成功: ${response.data.result}`)
    } else {
      ElMessage.info('规则未匹配')
    }
    
  } catch (error: any) {
    ElMessage.error(error.message || '规则测试失败')
  } finally {
    rule.testing = false
  }
}

// 切换规则状态
const toggleRule = async (rule: Rule) => {
  try {
    rule.toggling = true
    
    if (rule.enabled) {
      await rulesApi.enableRule(rule.id)
      ElMessage.success('规则已启用')
    } else {
      await rulesApi.disableRule(rule.id)
      ElMessage.success('规则已禁用')
    }
    
  } catch (error: any) {
    rule.enabled = !rule.enabled // 回滚状态
    ElMessage.error(error.message || '操作失败')
  } finally {
    rule.toggling = false
  }
}

// 删除规则
const deleteRule = async (rule: Rule) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除规则 "${rule.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await rulesApi.deleteRule(rule.id)
    ElMessage.success('规则删除成功')
    fetchRules()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 显示批量优先级调整对话框
const showBatchPriorityDialog = () => {
  if (selectedRules.value.length === 0) {
    ElMessage.warning('请选择要调整的规则')
    return
  }
  priorityDialogVisible.value = true
}

// 提交优先级调整
const submitPriorityAdjust = async () => {
  try {
    adjusting.value = true
    
    const updates = selectedRules.value.map(rule => {
      let newPriority = rule.priority
      
      switch (priorityAdjustType.value) {
        case 'set':
          newPriority = priorityValue.value
          break
        case 'increase':
          newPriority = rule.priority + priorityValue.value
          break
        case 'decrease':
          newPriority = rule.priority - priorityValue.value
          break
      }
      
      return {
        id: rule.id,
        priority: Math.max(1, Math.min(999, newPriority))
      }
    })
    
    await rulesApi.batchUpdatePriority(updates)
    ElMessage.success('优先级调整成功')
    
    priorityDialogVisible.value = false
    selectedRules.value = []
    fetchRules()
    
  } catch (error: any) {
    ElMessage.error(error.message || '优先级调整失败')
  } finally {
    adjusting.value = false
  }
}

// 导出规则
const exportRules = async () => {
  try {
    const response = await rulesApi.getRules({ limit: 1000 })
    const data = JSON.stringify(response.data.items, null, 2)
    
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `rules_${new Date().toISOString().split('T')[0]}.json`
    link.click()
    
    URL.revokeObjectURL(url)
    ElMessage.success('规则导出成功')
    
  } catch (error: any) {
    ElMessage.error(error.message || '导出失败')
  }
}

// 显示导入对话框
const showImportDialog = () => {
  importPreview.value = []
  importDialogVisible.value = true
}

// 处理文件变化
const handleFileChange = (file: any) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const content = e.target?.result as string
      const data = JSON.parse(content)
      
      if (Array.isArray(data)) {
        importPreview.value = data
      } else {
        ElMessage.error('文件格式错误，应为规则数组')
      }
    } catch (error) {
      ElMessage.error('文件解析失败，请检查JSON格式')
    }
  }
  reader.readAsText(file.raw)
}

// 提交导入
const submitImport = async () => {
  try {
    importing.value = true
    
    for (const rule of importPreview.value) {
      await rulesApi.createRule(rule)
    }
    
    ElMessage.success(`成功导入 ${importPreview.value.length} 条规则`)
    importDialogVisible.value = false
    importPreview.value = []
    fetchRules()
    
  } catch (error: any) {
    ElMessage.error(error.message || '导入失败')
  } finally {
    importing.value = false
  }
}

// 表单处理函数
const handleTypeChange = () => {
  ruleForm.condition = ''
}

const handleActionChange = () => {
  if (ruleForm.action !== 'redirect') {
    ruleForm.redirect_url = ''
  }
  if (ruleForm.action !== 'modify') {
    ruleForm.modify_config = ''
  }
}

// 工具函数
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    url: '',
    domain: 'success',
    ip: 'warning',
    user_agent: 'info'
  }
  return typeMap[type] || ''
}

const getTypeText = (type: string) => {
  const textMap: Record<string, string> = {
    url: 'URL',
    domain: '域名',
    ip: 'IP',
    user_agent: 'UA'
  }
  return textMap[type] || type
}

const getActionTagType = (action: string) => {
  const actionMap: Record<string, string> = {
    allow: 'success',
    block: 'danger',
    redirect: 'warning',
    modify: 'info'
  }
  return actionMap[action] || ''
}

const getActionText = (action: string) => {
  const textMap: Record<string, string> = {
    allow: '允许',
    block: '阻止',
    redirect: '重定向',
    modify: '修改'
  }
  return textMap[action] || action
}

const getConditionPlaceholder = (type: string) => {
  const placeholderMap: Record<string, string> = {
    url: '例如: https://example.com/* 或 *.example.com',
    domain: '例如: example.com 或 *.example.com',
    ip: '例如: 192.168.1.1 或 192.168.1.0/24',
    user_agent: '例如: *Chrome* 或 Mozilla/5.0*'
  }
  return placeholderMap[type] || '请输入匹配条件'
}

const getConditionTip = (type: string) => {
  const tipMap: Record<string, string> = {
    url: '支持通配符 * 和正则表达式',
    domain: '支持通配符 * 匹配子域名',
    ip: '支持单个IP或CIDR格式的IP段',
    user_agent: '支持通配符 * 匹配用户代理字符串'
  }
  return tipMap[type] || ''
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载
onMounted(() => {
  fetchRules()
})
</script>

<style scoped>
.rules-management {
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

/* 筛选栏 */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 16px 20px;
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
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

/* 规则表格 */
.rules-table {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  overflow: hidden;
}

.rule-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name-text {
  font-weight: 500;
  color: #fff;
}

.rule-condition {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  color: #00ff88;
  font-size: 13px;
}

.match-count {
  color: #ccc;
  font-weight: 500;
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

/* 优先级调整 */
.priority-form p {
  color: #ccc;
  margin-bottom: 20px;
}

/* 导入预览 */
.import-preview {
  margin-top: 20px;
  padding: 16px;
  background-color: #2a2a2a;
  border-radius: 8px;
}

.import-preview h4 {
  margin: 0 0 12px 0;
  color: #fff;
  font-size: 14px;
}

.preview-list {
  max-height: 200px;
  overflow-y: auto;
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid #333;
}

.preview-item:last-child {
  border-bottom: none;
}

.preview-item .rule-name {
  flex: 1;
  color: #ccc;
  font-size: 13px;
}

.more-items {
  padding: 8px 0;
  color: #888;
  font-size: 12px;
  text-align: center;
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
}

/* Element Plus 组件样式覆盖 */
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
  background-color: #1a1a1a;
}

:deep(.el-table__row:hover) {
  background-color: rgba(0, 255, 136, 0.05);
}

:deep(.el-table__row--striped) {
  background-color: #1f1f1f;
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

:deep(.el-upload) {
  width: 100%;
}

:deep(.el-upload__tip) {
  color: #888;
  margin-top: 8px;
}
</style>