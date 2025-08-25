import type { ThemeConfig } from './theme-types'
import { DEFAULT_THEME_CONFIG, THEME_STORAGE_KEY } from './theme-config'

export class ThemeStorage {
  static load(): Partial<ThemeConfig> {
    try {
      const savedSettings = localStorage.getItem(THEME_STORAGE_KEY)
      if (savedSettings) {
        return JSON.parse(savedSettings)
      }
      return {}
    } catch (error) {
      console.error('加载主题设置失败:', error)
      return {}
    }
  }

  static save(config: ThemeConfig): void {
    try {
      localStorage.setItem(THEME_STORAGE_KEY, JSON.stringify(config))
    } catch (error) {
      console.error('保存主题设置失败:', error)
    }
  }

  static loadWithDefaults(): ThemeConfig {
    const saved = this.load()
    return {
      isDark: saved.isDark ?? DEFAULT_THEME_CONFIG.isDark,
      themeMode: saved.themeMode ?? DEFAULT_THEME_CONFIG.themeMode,
      sidebarCollapsed: saved.sidebarCollapsed ?? DEFAULT_THEME_CONFIG.sidebarCollapsed,
      layoutMode: saved.layoutMode ?? DEFAULT_THEME_CONFIG.layoutMode,
      primaryColor: saved.primaryColor ?? DEFAULT_THEME_CONFIG.primaryColor,
      fontSize: saved.fontSize ?? DEFAULT_THEME_CONFIG.fontSize,
      borderRadius: saved.borderRadius ?? DEFAULT_THEME_CONFIG.borderRadius,
      animationEnabled: saved.animationEnabled ?? DEFAULT_THEME_CONFIG.animationEnabled,
      compactMode: saved.compactMode ?? DEFAULT_THEME_CONFIG.compactMode,
      showBreadcrumb: saved.showBreadcrumb ?? DEFAULT_THEME_CONFIG.showBreadcrumb,
      showTabsView: saved.showTabsView ?? DEFAULT_THEME_CONFIG.showTabsView,
      fixedHeader: saved.fixedHeader ?? DEFAULT_THEME_CONFIG.fixedHeader,
      fixedSidebar: saved.fixedSidebar ?? DEFAULT_THEME_CONFIG.fixedSidebar
    }
  }

  static clear(): void {
    try {
      localStorage.removeItem(THEME_STORAGE_KEY)
    } catch (error) {
      console.error('清除主题设置失败:', error)
    }
  }

  static export(): string {
    try {
      const config = this.load()
      return JSON.stringify(config, null, 2)
    } catch (error) {
      console.error('导出主题设置失败:', error)
      return '{}'
    }
  }

  static import(configJson: string): boolean {
    try {
      const config = JSON.parse(configJson)
      this.save(config)
      return true
    } catch (error) {
      console.error('导入主题设置失败:', error)
      return false
    }
  }
}
