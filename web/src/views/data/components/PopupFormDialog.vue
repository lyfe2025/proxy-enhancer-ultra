<template>
  <el-dialog
    v-model="dialogVisible"
    :title="editingPopup ? '编辑弹窗' : '创建弹窗'"
    width="800px"
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
          <el-form-item label="弹窗名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入弹窗名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="弹窗类型" prop="type">
            <el-select v-model="formData.type" placeholder="选择弹窗类型">
              <el-option label="信息收集" value="form" />
              <el-option label="问卷调查" value="survey" />
              <el-option label="反馈收集" value="feedback" />
              <el-option label="订阅邮件" value="newsletter" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="状态" prop="status">
            <el-select v-model="formData.status" placeholder="选择状态">
              <el-option label="启用" value="active" />
              <el-option label="禁用" value="inactive" />
              <el-option label="草稿" value="draft" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="显示位置" prop="position">
            <el-select v-model="formData.position" placeholder="选择显示位置">
              <el-option label="页面中央" value="center" />
              <el-option label="页面顶部" value="top" />
              <el-option label="页面底部" value="bottom" />
              <el-option label="右下角" value="bottom-right" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="弹窗标题" prop="title">
        <el-input v-model="formData.title" placeholder="请输入弹窗标题" />
      </el-form-item>

      <el-form-item label="弹窗描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入弹窗描述"
        />
      </el-form-item>

      <el-form-item label="表单字段" prop="fields">
        <div class="fields-editor">
          <div
            v-for="(field, index) in formData.fields"
            :key="index"
            class="field-item"
          >
            <el-select v-model="field.type" placeholder="字段类型" style="width: 120px;">
              <el-option label="文本" value="text" />
              <el-option label="邮箱" value="email" />
              <el-option label="电话" value="tel" />
              <el-option label="数字" value="number" />
              <el-option label="多行文本" value="textarea" />
              <el-option label="选择框" value="select" />
              <el-option label="复选框" value="checkbox" />
              <el-option label="单选框" value="radio" />
            </el-select>
            
            <el-input
              v-model="field.name"
              placeholder="字段名称"
              style="width: 150px;"
            />
            
            <el-input
              v-model="field.label"
              placeholder="字段标签"
              style="width: 150px;"
            />
            
            <el-input
              v-model="field.placeholder"
              placeholder="占位符"
              style="width: 150px;"
            />
            
            <el-checkbox v-model="field.required">必填</el-checkbox>
            
            <el-button
              size="small"
              type="danger"
              @click="removeField(index)"
              :disabled="formData.fields.length <= 1"
            >
              删除
            </el-button>
          </div>
          
          <el-button
            type="primary"
            size="small"
            @click="addField"
            style="margin-top: 10px;"
          >
            <el-icon><Plus /></el-icon>
            添加字段
          </el-button>
        </div>
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="显示时机" prop="trigger">
            <el-select v-model="formData.trigger" placeholder="选择显示时机">
              <el-option label="页面加载时" value="load" />
              <el-option label="停留时间" value="time" />
              <el-option label="滚动位置" value="scroll" />
              <el-option label="退出意图" value="exit" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="延迟时间" prop="delay" v-if="formData.trigger === 'time'">
            <el-input-number
              v-model="formData.delay"
              :min="0"
              :max="300"
              controls-position="right"
              style="width: 100%;"
            />
            <span style="color: var(--el-text-color-regular); font-size: 12px;">秒</span>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          {{ editingPopup ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

interface PopupField {
  type: string
  name: string
  label: string
  placeholder: string
  required: boolean
  options?: string[]
}

interface PopupFormData {
  name: string
  type: string
  status: string
  position: string
  title: string
  description: string
  fields: PopupField[]
  trigger: string
  delay: number
}

const props = defineProps<{
  editingPopup?: any
}>()

const emit = defineEmits<{
  submit: [data: PopupFormData]
  close: []
}>()

const dialogVisible = defineModel<boolean>({ required: true })
const formRef = ref()
const submitLoading = ref(false)

const formData = reactive<PopupFormData>({
  name: '',
  type: '',
  status: 'active',
  position: 'center',
  title: '',
  description: '',
  fields: [
    {
      type: 'text',
      name: 'name',
      label: '姓名',
      placeholder: '请输入您的姓名',
      required: true
    }
  ],
  trigger: 'load',
  delay: 3
})

const formRules = {
  name: [
    { required: true, message: '请输入弹窗名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择弹窗类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入弹窗标题', trigger: 'blur' }
  ]
}

// 监听编辑数据变化
watch(
  () => props.editingPopup,
  (newValue) => {
    if (newValue) {
      Object.assign(formData, {
        name: newValue.name || '',
        type: newValue.type || '',
        status: newValue.status || 'active',
        position: newValue.position || 'center',
        title: newValue.title || '',
        description: newValue.description || '',
        fields: newValue.fields || formData.fields,
        trigger: newValue.trigger || 'load',
        delay: newValue.delay || 3
      })
    }
  },
  { immediate: true }
)

const addField = () => {
  formData.fields.push({
    type: 'text',
    name: '',
    label: '',
    placeholder: '',
    required: false
  })
}

const removeField = (index: number) => {
  if (formData.fields.length > 1) {
    formData.fields.splice(index, 1)
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitLoading.value = true
    
    // 验证字段配置
    const invalidFields = formData.fields.filter(field => !field.name || !field.label)
    if (invalidFields.length > 0) {
      ElMessage.error('请完善所有字段的名称和标签')
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
    status: 'active',
    position: 'center',
    title: '',
    description: '',
    fields: [
      {
        type: 'text',
        name: 'name',
        label: '姓名',
        placeholder: '请输入您的姓名',
        required: true
      }
    ],
    trigger: 'load',
    delay: 3
  })
  
  emit('close')
}
</script>

<style scoped>
.fields-editor {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  padding: 12px;
  background: var(--el-bg-color-page);
}

.field-item {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.field-item:last-child {
  margin-bottom: 0;
}

.dialog-footer {
  text-align: right;
}

/* 深色主题适配 */
.dark .fields-editor {
  background: var(--el-bg-color);
  border-color: var(--el-border-color);
}
</style>
