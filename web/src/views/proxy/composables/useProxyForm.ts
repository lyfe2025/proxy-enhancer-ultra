import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { proxyApi, type ProxyConfig } from '@/api/proxy'
import type { ProxyCreateRequest, ProxyUpdateRequest } from '@/types/proxy'

export function useProxyForm() {
  // 响应式数据
  const saving = ref(false)
  const showFormDialog = ref(false)
  const editingProxy = ref<ProxyConfig | null>(null)

  // 编辑代理
  const editProxy = (proxy: ProxyConfig) => {
    editingProxy.value = proxy
    showFormDialog.value = true
  }

  // 保存代理（新增或编辑）
  const saveProxy = async (form: Partial<ProxyConfig> & {
    timeout?: number
    retry_count?: number
    enabled?: boolean
    tags?: string | string[]
  }) => {
    try {
      saving.value = true
      
      if (editingProxy.value) {
        // 更新代理
        const updateData: ProxyUpdateRequest = {
          name: form.name,
          type: form.type as 'http' | 'https' | 'socks5',
          host: form.host,
          port: form.port,
          username: form.username || undefined,
          password: form.password || undefined,
          description: form.description || undefined
        }
        await proxyApi.updateProxy(editingProxy.value.id!, updateData)
        ElMessage.success('更新成功')
      } else {
        // 创建代理
        const createData: Omit<ProxyConfig, 'id' | 'createdAt' | 'updatedAt'> = {
          name: form.name!,
          type: form.type!,
          host: form.host!,
          port: form.port!,
          username: form.username || undefined,
          password: form.password || undefined,
          status: 'inactive' as const,
          description: form.description || undefined
        }
        await proxyApi.createProxy(createData)
        ElMessage.success('添加成功')
      }
      
      showFormDialog.value = false
      return true // 表示保存成功
    } catch (error) {
      ElMessage.error(editingProxy.value ? '更新失败' : '添加失败')
      return false // 表示保存失败
    } finally {
      saving.value = false
    }
  }

  // 重置表单
  const resetForm = () => {
    editingProxy.value = null
  }

  // 显示添加对话框
  const showAddDialog = () => {
    editingProxy.value = null
    showFormDialog.value = true
  }

  // 关闭对话框
  const closeDialog = () => {
    showFormDialog.value = false
    resetForm()
  }

  return {
    // 响应式数据
    saving,
    showFormDialog,
    editingProxy,
    
    // 方法
    editProxy,
    saveProxy,
    resetForm,
    showAddDialog,
    closeDialog
  }
}