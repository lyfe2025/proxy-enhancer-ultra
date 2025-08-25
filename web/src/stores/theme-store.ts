import { ref, watch } from 'vue'
import { defineStore } from 'pinia'
import type { ThemeMode, LayoutMode, ThemeConfig, FontSize, BorderRadius } from './theme-types'
import { DEFAULT_THEME_CONFIG } from './theme-config'
import { ThemeStorage } from './theme-storage'
import { ThemeDOMManager } from './theme-dom'

export const useThemeStore = defineStore('theme', () => {
  // DOM管理器
  const domManager = new ThemeDOMManager()

  // 状态
  const isDark = ref(DEFAULT_THEME_CONFIG.isDark)
  const themeMode = ref<ThemeMode>(DEFAULT_THEME_CONFIG.themeMode)
  const sidebarCollapsed = ref(DEFAULT_THEME_CONFIG.sidebarCollapsed)
  const layoutMode = ref<LayoutMode>(DEFAULT_THEME_CONFIG.layoutMode)
  const primaryColor = ref(DEFAULT_THEME_CONFIG.primaryColor)
  const fontSize = ref(DEFAULT_THEME_CONFIG.fontSize)
  const borderRadius = ref(DEFAULT_THEME_CONFIG.borderRadius)
  const animationEnabled = ref(DEFAULT_THEME_CONFIG.animationEnabled)
  const compactMode = ref(DEFAULT_THEME_CONFIG.compactMode)
  const showBreadcrumb = ref(DEFAULT_THEME_CONFIG.showBreadcrumb)
  const showTabsView = ref(DEFAULT_THEME_CONFIG.showTabsView)
  const fixedHeader = ref(DEFAULT_THEME_CONFIG.fixedHeader)
  const fixedSidebar = ref(DEFAULT_THEME_CONFIG.fixedSidebar)

  // 从本地存储加载设置
  const loadSettings = () => {
    const config = ThemeStorage.loadWithDefaults()
    isDark.value = config.isDark
    themeMode.value = config.themeMode
    sidebarCollapsed.value = config.sidebarCollapsed
    layoutMode.value = config.layoutMode
    primaryColor.value = config.primaryColor
    fontSize.value = config.fontSize
    borderRadius.value = config.borderRadius
    animationEnabled.value = config.animationEnabled
    compactMode.value = config.compactMode
    showBreadcrumb.value = config.showBreadcrumb
    showTabsView.value = config.showTabsView
    fixedHeader.value = config.fixedHeader
    fixedSidebar.value = config.fixedSidebar
  }

  // 保存设置到本地存储
  const saveSettings = () => {
    const config = getThemeConfig()
    ThemeStorage.save(config)
  }

  // 应用主题到DOM
  const applyTheme = () => {
    const config = getThemeConfig()
    domManager.applyTheme(config)
    
    // 更新isDark状态
    if (themeMode.value === 'auto') {
      isDark.value = domManager.getCurrentSystemTheme() === 'dark'
    } else {
      isDark.value = themeMode.value === 'dark'
    }
  }

  // 切换主题模式
  const toggleTheme = () => {
    if (themeMode.value === 'light') {
      themeMode.value = 'dark'
    } else if (themeMode.value === 'dark') {
      themeMode.value = 'auto'
    } else {
      themeMode.value = 'light'
    }
    applyTheme()
    saveSettings()
  }

  // 设置主题模式
  const setThemeMode = (mode: ThemeMode) => {
    themeMode.value = mode
    applyTheme()
    saveSettings()
  }

  // 侧边栏操作
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
    saveSettings()
  }

  const setSidebarCollapsed = (collapsed: boolean) => {
    sidebarCollapsed.value = collapsed
    saveSettings()
  }

  // 布局设置
  const setLayoutMode = (mode: LayoutMode) => {
    layoutMode.value = mode
    applyTheme()
    saveSettings()
  }

  // 主色调设置
  const setPrimaryColor = (color: string) => {
    primaryColor.value = color
    applyTheme()
    saveSettings()
  }

  // 字体大小设置
  const setFontSize = (size: FontSize) => {
    fontSize.value = size
    applyTheme()
    saveSettings()
  }

  // 圆角设置
  const setBorderRadius = (radius: BorderRadius) => {
    borderRadius.value = radius
    applyTheme()
    saveSettings()
  }

  // 切换功能
  const toggleAnimation = () => {
    animationEnabled.value = !animationEnabled.value
    applyTheme()
    saveSettings()
  }

  const toggleCompactMode = () => {
    compactMode.value = !compactMode.value
    applyTheme()
    saveSettings()
  }

  const toggleBreadcrumb = () => {
    showBreadcrumb.value = !showBreadcrumb.value
    saveSettings()
  }

  const toggleTabsView = () => {
    showTabsView.value = !showTabsView.value
    saveSettings()
  }

  const toggleFixedHeader = () => {
    fixedHeader.value = !fixedHeader.value
    saveSettings()
  }

  const toggleFixedSidebar = () => {
    fixedSidebar.value = !fixedSidebar.value
    saveSettings()
  }

  // 重置所有设置
  const resetSettings = () => {
    const config = DEFAULT_THEME_CONFIG
    isDark.value = config.isDark
    themeMode.value = config.themeMode
    sidebarCollapsed.value = config.sidebarCollapsed
    layoutMode.value = config.layoutMode
    primaryColor.value = config.primaryColor
    fontSize.value = config.fontSize
    borderRadius.value = config.borderRadius
    animationEnabled.value = config.animationEnabled
    compactMode.value = config.compactMode
    showBreadcrumb.value = config.showBreadcrumb
    showTabsView.value = config.showTabsView
    fixedHeader.value = config.fixedHeader
    fixedSidebar.value = config.fixedSidebar
    
    applyTheme()
    saveSettings()
  }

  // 获取当前主题配置
  const getThemeConfig = (): ThemeConfig => {
    return {
      isDark: isDark.value,
      themeMode: themeMode.value,
      sidebarCollapsed: sidebarCollapsed.value,
      layoutMode: layoutMode.value,
      primaryColor: primaryColor.value,
      fontSize: fontSize.value,
      borderRadius: borderRadius.value,
      animationEnabled: animationEnabled.value,
      compactMode: compactMode.value,
      showBreadcrumb: showBreadcrumb.value,
      showTabsView: showTabsView.value,
      fixedHeader: fixedHeader.value,
      fixedSidebar: fixedSidebar.value
    }
  }

  // 初始化主题
  const initTheme = () => {
    loadSettings()
    applyTheme()
    
    // 监听系统主题变化
    domManager.watchSystemTheme(() => {
      if (themeMode.value === 'auto') {
        applyTheme()
      }
    })
  }

  // 监听设置变化并自动保存
  watch(
    [isDark, themeMode, sidebarCollapsed, layoutMode, primaryColor, fontSize, borderRadius, animationEnabled, compactMode, showBreadcrumb, showTabsView, fixedHeader, fixedSidebar],
    () => {
      saveSettings()
    },
    { deep: true }
  )

  return {
    // 状态
    isDark,
    themeMode,
    sidebarCollapsed,
    layoutMode,
    primaryColor,
    fontSize,
    borderRadius,
    animationEnabled,
    compactMode,
    showBreadcrumb,
    showTabsView,
    fixedHeader,
    fixedSidebar,
    
    // 方法
    toggleTheme,
    setThemeMode,
    toggleSidebar,
    setSidebarCollapsed,
    setLayoutMode,
    setPrimaryColor,
    setFontSize,
    setBorderRadius,
    toggleAnimation,
    toggleCompactMode,
    toggleBreadcrumb,
    toggleTabsView,
    toggleFixedHeader,
    toggleFixedSidebar,
    resetSettings,
    getThemeConfig,
    initTheme,
    applyTheme,
    loadSettings,
    saveSettings
  }
})
