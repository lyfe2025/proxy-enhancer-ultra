<template>
  <el-dialog
    v-model="dialogVisible"
    :title="editingRule ? '编辑规则' : '添加规则'"
    width="900px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="规则名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入规则名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="规则类型" prop="type">
            <el-select v-model="formData.type" placeholder="选择规则类型">
              <el-option label="路由规则" value="route" />
              <el-option label="过滤规则" value="filter" />
              <el-option label="重写规则" value="rewrite" />
              <el-option label="限流规则" value="rate_limit" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="优先级" prop="priority">
            <el-select v-model="formData.priority" placeholder="选择优先级">
              <el-option label="高" value="high" />
              <el-option label="中" value="medium" />
              <el-option label="低" value="low" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="状态" prop="status">
            <el-select v-model="formData.status" placeholder="选择状态">
              <el-option label="启用" value="enabled" />
              <el-option label="禁用" value="disabled" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="规则描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入规则描述"
        />
      </el-form-item>

      <el-form-item label="匹配条件" prop="conditions">
        <RuleConditionsEditor v-model="formData.conditions" />
      </el-form-item>

      <el-form-item label="执行动作" prop="action">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-select v-model="formData.action.type" placeholder="选择动作类型">
              <el-option label="代理转发" value="proxy" />
              <el-option label="重定向" value="redirect" />
              <el-option label="阻止请求" value="block" />
              <el-option label="修改请求" value="modify" />
            </el-select>
          </el-col>
          <el-col :span="16">
            <el-input
              v-model="formData.action.target"
              placeholder="目标地址或配置"
              v-if="formData.action.type !== 'block'"
            />
            <el-input
              v-model="formData.action.reason"
              placeholder="阻止原因"
              v-else
            />
          </el-col>
        </el-row>
      </el-form-item>

      <el-form-item label="高级设置">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="超时时间" prop="timeout">
              <el-input-number
                v-model="formData.timeout"
                :min="1"
                :max="300"
                controls-position="right"
                style="width: 100%;"
              />
              <span style="color: var(--el-text-color-regular); font-size: 12px; margin-left: 8px;">秒</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="重试次数" prop="retries">
              <el-input-number
                v-model="formData.retries"
                :min="0"
                :max="5"
                controls-position="right"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-checkbox v-model="formData.logEnabled">启用日志</el-checkbox>
          </el-col>
        </el-row>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          {{ editingRule ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import RuleConditionsEditor from './RuleConditionsEditor.vue'

interface RuleCondition {
  field: string
  operator: string
  value: string
  logic: 'AND' | 'OR'
}

interface RuleAction {
  type: string
  target?: string
  reason?: string
}

interface RuleFormData {
  name: string
  type: string
  priority: string
  status: string
  description: string
  conditions: RuleCondition[]
  action: RuleAction
  timeout: number
  retries: number
  logEnabled: boolean
}

const props = defineProps<{
  editingRule?: any
}>()

const emit = defineEmits<{
  submit: [data: RuleFormData]
  close: []
}>()

const dialogVisible = defineModel<boolean>({ required: true })
const formRef = ref()
const submitLoading = ref(false)

const formData = reactive<RuleFormData>({
  name: '',
  type: '',
  priority: 'medium',
  status: 'enabled',
  description: '',
  conditions: [
    {
      field: 'url',
      operator: 'contains',
      value: '',
      logic: 'AND'
    }
  ],
  action: {
    type: 'proxy',
    target: ''
  },
  timeout: 30,
  retries: 3,
  logEnabled: true
})

const formRules = {
  name: [
    { required: true, message: '请输入规则名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择规则类型', trigger: 'change' }
  ],
  priority: [
    { required: true, message: '请选择优先级', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  'action.type': [
    { required: true, message: '请选择动作类型', trigger: 'change' }
  ]
}

// 监听编辑数据变化
watch(
  () => props.editingRule,
  (newValue) => {
    if (newValue) {
      Object.assign(formData, {
        name: newValue.name || '',
        type: newValue.type || '',
        priority: newValue.priority || 'medium',
        status: newValue.status || 'enabled',
        description: newValue.description || '',
        conditions: newValue.conditions || formData.conditions,
        action: newValue.action || formData.action,
        timeout: newValue.timeout || 30,
        retries: newValue.retries || 3,
        logEnabled: newValue.logEnabled !== undefined ? newValue.logEnabled : true
      })
    }
  },
  { immediate: true }
)

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitLoading.value = true
    
    // 验证条件配置
    const invalidConditions = formData.conditions.filter(condition => 
      !condition.field || !condition.operator || !condition.value
    )
    if (invalidConditions.length > 0) {
      ElMessage.error('请完善所有匹配条件')
      return
    }
    
    // 验证动作配置
    if (formData.action.type !== 'block' && !formData.action.target) {
      ElMessage.error('请输入目标地址或配置')
      return
    }
    
    emit('submit', { ...formData })
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    submitLoading.value = false
  }
}

const handleClose = () => {
  // 重置表单
  Object.assign(formData, {
    name: '',
    type: '',
    priority: 'medium',
    status: 'enabled',
    description: '',
    conditions: [
      {
        field: 'url',
        operator: 'contains',
        value: '',
        logic: 'AND'
      }
    ],
    action: {
      type: 'proxy',
      target: ''
    },
    timeout: 30,
    retries: 3,
    logEnabled: true
  })
  
  emit('close')
}
</script>

<style scoped>
.dialog-footer {
  text-align: right;
}
</style>
