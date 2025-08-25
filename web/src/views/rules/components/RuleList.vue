<template>
  <el-card class="table-card">
    <template #header>
      <div class="card-header">
        <span>规则列表</span>
        <div class="header-actions">
          <el-button size="small" @click="batchEnable" :disabled="!selectedRules.length">
            <el-icon><CircleCheck /></el-icon>
            批量启用
          </el-button>
          <el-button size="small" @click="batchDisable" :disabled="!selectedRules.length">
            <el-icon><CircleClose /></el-icon>
            批量禁用
          </el-button>
          <el-button size="small" type="danger" @click="batchDelete" :disabled="!selectedRules.length">
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
        </div>
      </div>
    </template>

    <el-table
      v-loading="loading"
      :data="rulesList"
      @selection-change="handleSelectionChange"
      stripe
      style="width: 100%"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="规则名称" min-width="150">
        <template #default="{ row }">
          <div class="rule-name">
            <el-icon class="rule-icon"><SetUp /></el-icon>
            <span>{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)">{{ getTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="priority" label="优先级" width="100">
        <template #default="{ row }">
          <el-tag :type="getPriorityTagType(row.priority)" size="small">
            {{ getPriorityText(row.priority) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="conditions" label="条件" min-width="200">
        <template #default="{ row }">
          <div class="conditions">
            <el-tag
              v-for="(condition, index) in row.conditions.slice(0, 2)"
              :key="index"
              size="small"
              class="condition-tag"
            >
              {{ formatCondition(condition) }}
            </el-tag>
            <span v-if="row.conditions.length > 2" class="more-conditions">
              +{{ row.conditions.length - 2 }}
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="action" label="动作" width="120">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ getActionText(row.action) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            active-value="enabled"
            inactive-value="disabled"
            @change="toggleStatus(row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="executedCount" label="执行次数" width="120">
        <template #default="{ row }">
          <span>{{ row.executedCount || 0 }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="lastExecuted" label="最后执行" width="160">
        <template #default="{ row }">
          <span class="text-muted">{{ formatDate(row.lastExecuted) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="moveUp(row)" :disabled="isFirst(row)">
            <el-icon><ArrowUp /></el-icon>
            上移
          </el-button>
          <el-button size="small" @click="moveDown(row)" :disabled="isLast(row)">
            <el-icon><ArrowDown /></el-icon>
            下移
          </el-button>
          <el-button size="small" @click="edit(row)">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-button size="small" type="danger" @click="deleteRule(row)">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  SetUp,
  CircleCheck,
  CircleClose,
  Delete,
  ArrowUp,
  ArrowDown,
  Edit
} from '@element-plus/icons-vue'

interface RuleItem {
  id: string
  name: string
  type: string
  priority: string
  conditions: any[]
  action: string
  status: string
  executedCount: number
  lastExecuted: string
}

interface PaginationData {
  page: number
  size: number
  total: number
}

const props = defineProps<{
  rulesList: RuleItem[]
  loading: boolean
  pagination: PaginationData
}>()

const emit = defineEmits<{
  batchEnable: [ids: string[]]
  batchDisable: [ids: string[]]
  batchDelete: [ids: string[]]
  toggleStatus: [rule: RuleItem]
  moveUp: [rule: RuleItem]
  moveDown: [rule: RuleItem]
  edit: [rule: RuleItem]
  delete: [rule: RuleItem]
  sizeChange: [size: number]
  currentChange: [page: number]
}>()

const selectedRules = ref<RuleItem[]>([])

const handleSelectionChange = (selection: RuleItem[]) => {
  selectedRules.value = selection
}

const batchEnable = async () => {
  try {
    await ElMessageBox.confirm('确认批量启用选中的规则吗？', '提示', {
      type: 'warning'
    })
    emit('batchEnable', selectedRules.value.map(rule => rule.id))
  } catch {}
}

const batchDisable = async () => {
  try {
    await ElMessageBox.confirm('确认批量禁用选中的规则吗？', '提示', {
      type: 'warning'
    })
    emit('batchDisable', selectedRules.value.map(rule => rule.id))
  } catch {}
}

const batchDelete = async () => {
  try {
    await ElMessageBox.confirm('确认批量删除选中的规则吗？此操作不可恢复！', '危险操作', {
      type: 'error'
    })
    emit('batchDelete', selectedRules.value.map(rule => rule.id))
  } catch {}
}

const toggleStatus = (rule: RuleItem) => {
  emit('toggleStatus', rule)
}

const moveUp = (rule: RuleItem) => {
  emit('moveUp', rule)
}

const moveDown = (rule: RuleItem) => {
  emit('moveDown', rule)
}

const edit = (rule: RuleItem) => {
  emit('edit', rule)
}

const deleteRule = async (rule: RuleItem) => {
  try {
    await ElMessageBox.confirm(`确认删除规则"${rule.name}"吗？此操作不可恢复！`, '危险操作', {
      type: 'error'
    })
    emit('delete', rule)
  } catch {}
}

const handleSizeChange = (size: number) => {
  emit('sizeChange', size)
}

const handleCurrentChange = (page: number) => {
  emit('currentChange', page)
}

const isFirst = (rule: RuleItem) => {
  return props.rulesList.indexOf(rule) === 0
}

const isLast = (rule: RuleItem) => {
  return props.rulesList.indexOf(rule) === props.rulesList.length - 1
}

// 工具函数
const getTypeTagType = (type: string): "success" | "info" | "warning" | "danger" | "primary" => {
  const typeMap: Record<string, "success" | "info" | "warning" | "danger" | "primary"> = {
    route: 'primary',
    filter: 'success',
    rewrite: 'warning',
    rate_limit: 'danger'
  }
  return typeMap[type] || 'info'
}

const getTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    route: '路由规则',
    filter: '过滤规则',
    rewrite: '重写规则',
    rate_limit: '限流规则'
  }
  return typeMap[type] || type
}

const getPriorityTagType = (priority: string): "success" | "info" | "warning" | "danger" | "primary" => {
  const priorityMap: Record<string, "success" | "info" | "warning" | "danger" | "primary"> = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return priorityMap[priority] || 'info'
}

const getPriorityText = (priority: string) => {
  const priorityMap: Record<string, string> = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return priorityMap[priority] || priority
}

const getActionText = (action: string) => {
  const actionMap: Record<string, string> = {
    proxy: '代理',
    redirect: '重定向',
    block: '阻止',
    modify: '修改'
  }
  return actionMap[action] || action
}

const formatCondition = (condition: any) => {
  if (typeof condition === 'string') return condition
  if (condition.field && condition.operator && condition.value) {
    return `${condition.field} ${condition.operator} ${condition.value}`
  }
  return JSON.stringify(condition)
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.table-card {
  border: 1px solid var(--el-border-color-light);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.rule-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rule-icon {
  color: var(--el-color-primary);
}

.conditions {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.condition-tag {
  margin: 0;
}

.more-conditions {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}

.text-muted {
  color: var(--el-text-color-placeholder);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 深色主题适配 */
.dark .table-card {
  background: var(--el-bg-color-overlay);
  border-color: var(--el-border-color);
}
</style>
