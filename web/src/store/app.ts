import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 状态
  const sidebarCollapsed = ref(false)
  const theme = ref<'dark' | 'light'>('dark')
  const loading = ref(false)
  const breadcrumbs = ref<{ title: string; path?: string }[]>([])
  
  // 设备类型检测
  const isMobile = ref(false)
  const isTablet = ref(false)
  
  // 计算属性
  const deviceType = computed(() => {
    if (isMobile.value) return 'mobile'
    if (isTablet.value) return 'tablet'
    return 'desktop'
  })
  
  // 切换侧边栏
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem('sidebarCollapsed', String(sidebarCollapsed.value))
  }
  
  // 设置侧边栏状态
  const setSidebarCollapsed = (collapsed: boolean) => {
    sidebarCollapsed.value = collapsed
    localStorage.setItem('sidebarCollapsed', String(collapsed))
  }
  
  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    localStorage.setItem('theme', theme.value)
    updateThemeClass()
  }
  
  // 设置主题
  const setTheme = (newTheme: 'dark' | 'light') => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    updateThemeClass()
  }
  
  // 更新主题类名
  const updateThemeClass = () => {
    const html = document.documentElement
    if (theme.value === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }
  
  // 设置全局加载状态
  const setLoading = (isLoading: boolean) => {
    loading.value = isLoading
  }
  
  // 设置面包屑
  const setBreadcrumbs = (crumbs: { title: string; path?: string }[]) => {
    breadcrumbs.value = crumbs
  }
  
  // 添加面包屑
  const addBreadcrumb = (crumb: { title: string; path?: string }) => {
    breadcrumbs.value.push(crumb)
  }
  
  // 清空面包屑
  const clearBreadcrumbs = () => {
    breadcrumbs.value = []
  }
  
  // 检测设备类型
  const detectDevice = () => {
    const width = window.innerWidth
    isMobile.value = width < 768
    isTablet.value = width >= 768 && width < 1024
    
    // 移动端自动收起侧边栏
    if (isMobile.value && !sidebarCollapsed.value) {
      setSidebarCollapsed(true)
    }
  }
  
  // 初始化应用状态
  const initApp = () => {
    // 从localStorage恢复状态
    const savedSidebarState = localStorage.getItem('sidebarCollapsed')
    if (savedSidebarState !== null) {
      sidebarCollapsed.value = savedSidebarState === 'true'
    }
    
    const savedTheme = localStorage.getItem('theme') as 'dark' | 'light'
    if (savedTheme) {
      theme.value = savedTheme
    }
    
    // 应用主题
    updateThemeClass()
    
    // 检测设备类型
    detectDevice()
    
    // 监听窗口大小变化
    window.addEventListener('resize', detectDevice)
  }
  
  // 清理资源
  const cleanup = () => {
    window.removeEventListener('resize', detectDevice)
  }
  
  return {
    // 状态
    sidebarCollapsed: readonly(sidebarCollapsed),
    theme: readonly(theme),
    loading: readonly(loading),
    breadcrumbs: readonly(breadcrumbs),
    isMobile: readonly(isMobile),
    isTablet: readonly(isTablet),
    
    // 计算属性
    deviceType,
    
    // 方法
    toggleSidebar,
    setSidebarCollapsed,
    toggleTheme,
    setTheme,
    setLoading,
    setBreadcrumbs,
    addBreadcrumb,
    clearBreadcrumbs,
    detectDevice,
    initApp,
    cleanup
  }
})