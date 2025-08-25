<template>
  <div class="proxy-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>代理管理</h2>
        <p>管理和配置代理服务器设置</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新增代理
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filter-bar">
      <div class="filter-left">
        <el-input
          v-model="searchQuery"
          placeholder="搜索代理名称或地址"
          style="width: 300px"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select
          v-model="statusFilter"
          placeholder="状态筛选"
          style="width: 120px"
          clearable
          @change="handleFilter"
        >
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="错误" value="error" />
        </el-select>
        
        <el-select
          v-model="typeFilter"
          placeholder="类型筛选"
          style="width: 120px"
          clearable
          @change="handleFilter"
        >
          <el-option label="HTTP" value="http" />
          <el-option label="HTTPS" value="https" />
          <el-option label="SOCKS5" value="socks5" />
        </el-select>
      </div>
      
      <div class="filter-right">
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="batchDelete" :disabled="selectedProxies.length === 0">
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
      </div>
    </div>

    <!-- 代理列表 -->
    <div class="proxy-table">
      <el-table
        v-loading="loading"
        :data="filteredProxies"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="代理名称" min-width="150">
          <template #default="{ row }">
            <div class="proxy-name">
              <span class="name-text">{{ row.name }}</span>
              <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ row.type.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="host" label="地址" min-width="200">
          <template #default="{ row }">
            <span class="proxy-address">{{ row.host }}:{{ row.port }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small">
              <el-icon><component :is="getStatusIcon(row.status)" /></el-icon>
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="auth_required" label="认证" width="80">
          <template #default="{ row }">
            <el-icon v-if="row.auth_required" color="#00ff88"><Check /></el-icon>
            <el-icon v-else color="#666"><Close /></el-icon>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                type="text"
                size="small"
                @click="testProxy(row)"
                :loading="row.testing"
              >
                <el-icon><Connection /></el-icon>
                测试
              </el-button>
              
              <el-button
                type="text"
                size="small"
                @click="toggleProxy(row)"
                :loading="row.toggling"
              >
                <el-icon><component :is="row.enabled ? 'VideoPause' : 'VideoPlay'" /></el-icon>
                {{ row.enabled ? '停用' : '启用' }}
              </el-button>
              
              <el-button type="text" size="small" @click="editProxy(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              
              <el-button
                type="text"
                size="small"
                @click="deleteProxy(row)"
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

    <!-- 创建/编辑代理对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑代理' : '新增代理'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="proxyForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="代理名称" prop="name">
          <el-input v-model="proxyForm.name" placeholder="请输入代理名称" />
        </el-form-item>
        
        <el-form-item label="代理类型" prop="type">
          <el-select v-model="proxyForm.type" style="width: 100%">
            <el-option label="HTTP" value="http" />
            <el-option label="HTTPS" value="https" />
            <el-option label="SOCKS5" value="socks5" />
          </el-select>
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="16">
            <el-form-item label="主机地址" prop="host">
              <el-input v-model="proxyForm.host" placeholder="请输入主机地址" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="端口" prop="port">
              <el-input-number
                v-model="proxyForm.port"
                :min="1"
                :max="65535"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item>
          <el-checkbox v-model="proxyForm.auth_required">需要认证</el-checkbox>
        </el-form-item>
        
        <template v-if="proxyForm.auth_required">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="proxyForm.username" placeholder="请输入用户名" />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="proxyForm.password"
              type="password"
              placeholder="请输入密码"
              show-password
            />
          </el-form-item>
        </template>
        
        <el-form-item>
          <el-checkbox v-model="proxyForm.is_default">设为默认代理</el-checkbox>
        </el-form-item>
        
        <el-form-item label="描述">
          <el-input
            v-model="proxyForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入代理描述（可选）"
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
  Check,
  Close,
  VideoPlay,
  VideoPause,
  SuccessFilled,
  WarningFilled,
  CircleCloseFilled
} from '@element-plus/icons-vue'
import { proxyApi } from '@/api'
import type { ProxyConfig } from '@/types'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const typeFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedProxies = ref<ProxyConfig[]>([])
const proxies = ref<ProxyConfig[]>([])

// 表单引用和数据
const formRef = ref()
const proxyForm = reactive({
  id: '',
  name: '',
  type: 'http',
  host: '',
  port: 8080,
  auth_required: false,
  username: '',
  password: '',
  is_default: false,
  description: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入代理名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择代理类型', trigger: 'change' }
  ],
  host: [
    { required: true, message: '请输入主机地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: '端口号范围 1-65535', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 计算属性
const filteredProxies = computed(() => {
  let result = proxies.value
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(proxy => 
      proxy.name.toLowerCase().includes(query) ||
      proxy.host.toLowerCase().includes(query)
    )
  }
  
  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(proxy => proxy.status === statusFilter.value)
  }
  
  // 类型过滤
  if (typeFilter.value) {
    result = result.filter(proxy => proxy.type === typeFilter.value)
  }
  
  return result
})

// 获取代理列表
const fetchProxies = async () => {
  try {
    loading.value = true
    const response = await proxyApi.getProxies({
      page: currentPage.value,
      limit: pageSize.value
    })
    
    proxies.value = response.data.items.map(proxy => ({
      ...proxy,
      testing: false,
      toggling: false
    }))
    total.value = response.data.total
    
  } catch (error: any) {
    ElMessage.error(error.message || '获取代理列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  // 这里可以添加防抖逻辑
}

// 筛选处理
const handleFilter = () => {
  currentPage.value = 1
}

// 刷新数据
const refreshData = () => {
  fetchProxies()
}

// 选择变化处理
const handleSelectionChange = (selection: ProxyConfig[]) => {
  selectedProxies.value = selection
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchProxies()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchProxies()
}

// 显示创建对话框
const showCreateDialog = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 编辑代理
const editProxy = (proxy: ProxyConfig) => {
  isEdit.value = true
  Object.assign(proxyForm, {
    id: proxy.id,
    name: proxy.name,
    type: proxy.type,
    host: proxy.host,
    port: proxy.port,
    auth_required: proxy.auth_required,
    username: proxy.username || '',
    password: '', // 不显示原密码
    is_default: proxy.is_default,
    description: proxy.description || ''
  })
  dialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  Object.assign(proxyForm, {
    id: '',
    name: '',
    type: 'http',
    host: '',
    port: 8080,
    auth_required: false,
    username: '',
    password: '',
    is_default: false,
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
    
    const data = { ...proxyForm }
    delete data.id
    
    if (isEdit.value) {
      await proxyApi.updateProxy(proxyForm.id, data)
      ElMessage.success('代理更新成功')
    } else {
      await proxyApi.createProxy(data)
      ElMessage.success('代理创建成功')
    }
    
    dialogVisible.value = false
    fetchProxies()
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 测试代理
const testProxy = async (proxy: ProxyConfig) => {
  try {
    proxy.testing = true
    const response = await proxyApi.testProxy(proxy.id)
    
    if (response.data.success) {
      ElMessage.success('代理测试成功')
    } else {
      ElMessage.error(`代理测试失败: ${response.data.error}`)
    }
    
  } catch (error: any) {
    ElMessage.error(error.message || '代理测试失败')
  } finally {
    proxy.testing = false
  }
}

// 切换代理状态
const toggleProxy = async (proxy: ProxyConfig) => {
  try {
    proxy.toggling = true
    
    if (proxy.enabled) {
      await proxyApi.disableProxy(proxy.id)
      proxy.enabled = false
      proxy.status = 'stopped'
      ElMessage.success('代理已停用')
    } else {
      await proxyApi.enableProxy(proxy.id)
      proxy.enabled = true
      proxy.status = 'running'
      ElMessage.success('代理已启用')
    }
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    proxy.toggling = false
  }
}

// 删除代理
const deleteProxy = async (proxy: ProxyConfig) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除代理 "${proxy.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await proxyApi.deleteProxy(proxy.id)
    ElMessage.success('代理删除成功')
    fetchProxies()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 批量删除
const batchDelete = async () => {
  if (selectedProxies.value.length === 0) {
    ElMessage.warning('请选择要删除的代理')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedProxies.value.length} 个代理吗？此操作不可恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ids = selectedProxies.value.map(proxy => proxy.id)
    await Promise.all(ids.map(id => proxyApi.deleteProxy(id)))
    
    ElMessage.success('批量删除成功')
    selectedProxies.value = []
    fetchProxies()
    
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '批量删除失败')
    }
  }
}

// 工具函数
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    http: '',
    https: 'success',
    socks5: 'warning'
  }
  return typeMap[type] || ''
}

const getStatusTagType = (status: string) => {
  const statusMap: Record<string, string> = {
    running: 'success',
    stopped: 'info',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusIcon = (status: string) => {
  const iconMap: Record<string, string> = {
    running: 'SuccessFilled',
    stopped: 'VideoPause',
    error: 'CircleCloseFilled'
  }
  return iconMap[status] || 'WarningFilled'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '错误'
  }
  return textMap[status] || '未知'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载
onMounted(() => {
  fetchProxies()
})
</script>

<style scoped>
.proxy-management {
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

/* 代理表格 */
.proxy-table {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  overflow: hidden;
}

.proxy-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name-text {
  font-weight: 500;
  color: #fff;
}

.proxy-address {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  color: #00ff88;
  font-size: 13px;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .filter-left,
  .filter-right {
    justify-content: flex-start;
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
</style>