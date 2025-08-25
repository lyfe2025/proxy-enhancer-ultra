<template>
  <div class="conditions-editor">
    <div
      v-for="(condition, index) in conditions"
      :key="index"
      class="condition-item"
    >
      <el-select v-model="condition.field" placeholder="字段" style="width: 120px;">
        <el-option label="URL" value="url" />
        <el-option label="域名" value="host" />
        <el-option label="路径" value="path" />
        <el-option label="方法" value="method" />
        <el-option label="请求头" value="header" />
        <el-option label="查询参数" value="query" />
        <el-option label="IP地址" value="ip" />
        <el-option label="用户代理" value="user_agent" />
      </el-select>
      
      <el-select v-model="condition.operator" placeholder="操作符" style="width: 120px;">
        <el-option label="等于" value="equals" />
        <el-option label="包含" value="contains" />
        <el-option label="开始于" value="starts_with" />
        <el-option label="结束于" value="ends_with" />
        <el-option label="正则匹配" value="regex" />
        <el-option label="不等于" value="not_equals" />
        <el-option label="不包含" value="not_contains" />
      </el-select>
      
      <el-input
        v-model="condition.value"
        placeholder="值"
        style="width: 200px;"
      />
      
      <el-select 
        v-model="condition.logic" 
        placeholder="逻辑" 
        style="width: 80px;"
        v-if="index < conditions.length - 1"
      >
        <el-option label="并且" value="AND" />
        <el-option label="或者" value="OR" />
      </el-select>
      
      <el-button
        size="small"
        type="danger"
        @click="removeCondition(index)"
        :disabled="conditions.length <= 1"
      >
        删除
      </el-button>
    </div>
    
    <el-button
      type="primary"
      size="small"
      @click="addCondition"
      style="margin-top: 10px;"
    >
      <el-icon><Plus /></el-icon>
      添加条件
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'

interface RuleCondition {
  field: string
  operator: string
  value: string
  logic: 'AND' | 'OR'
}

const conditions = defineModel<RuleCondition[]>({ required: true })

const addCondition = () => {
  conditions.value.push({
    field: 'url',
    operator: 'contains',
    value: '',
    logic: 'AND'
  })
}

const removeCondition = (index: number) => {
  if (conditions.value.length > 1) {
    conditions.value.splice(index, 1)
  }
}
</script>

<style scoped>
.conditions-editor {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  padding: 12px;
  background: var(--el-bg-color-page);
}

.condition-item {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.condition-item:last-child {
  margin-bottom: 0;
}

/* 深色主题适配 */
.dark .conditions-editor {
  background: var(--el-bg-color);
  border-color: var(--el-border-color);
}
</style>
