import type { ThemeConfig, ThemeMode } from './theme-types'
import { FONT_SIZE_MAP, BORDER_RADIUS_MAP, LAYOUT_CLASSES, THEME_CLASSES } from './theme-config'

export class ThemeDOMManager {
  private root: HTMLElement

  constructor() {
    this.root = document.documentElement
  }

  applyTheme(config: ThemeConfig): void {
    this.applyThemeMode(config.themeMode, config.isDark)
    this.applyPrimaryColor(config.primaryColor)
    this.applyFontSize(config.fontSize)
    this.applyBorderRadius(config.borderRadius)
    this.applyAnimations(config.animationEnabled)
    this.applyCompactMode(config.compactMode)
    this.applyLayoutMode(config.layoutMode)
  }

  private applyThemeMode(themeMode: ThemeMode, isDark: boolean): void {
    let actualIsDark = isDark

    if (themeMode === 'auto') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      actualIsDark = prefersDark
    } else {
      actualIsDark = themeMode === 'dark'
    }

    if (actualIsDark) {
      this.root.classList.add(THEME_CLASSES.dark)
      this.root.classList.remove(THEME_CLASSES.light)
    } else {
      this.root.classList.add(THEME_CLASSES.light)
      this.root.classList.remove(THEME_CLASSES.dark)
    }
  }

  private applyPrimaryColor(color: string): void {
    this.root.style.setProperty('--primary-color', color)
  }

  private applyFontSize(fontSize: string): void {
    const size = FONT_SIZE_MAP[fontSize as keyof typeof FONT_SIZE_MAP] || '16px'
    this.root.style.setProperty('--base-font-size', size)
  }

  private applyBorderRadius(borderRadius: string): void {
    const radius = BORDER_RADIUS_MAP[borderRadius as keyof typeof BORDER_RADIUS_MAP] || '8px'
    this.root.style.setProperty('--border-radius', radius)
  }

  private applyAnimations(enabled: boolean): void {
    if (!enabled) {
      this.root.classList.add(THEME_CLASSES.noAnimations)
    } else {
      this.root.classList.remove(THEME_CLASSES.noAnimations)
    }
  }

  private applyCompactMode(enabled: boolean): void {
    if (enabled) {
      this.root.classList.add(THEME_CLASSES.compactMode)
    } else {
      this.root.classList.remove(THEME_CLASSES.compactMode)
    }
  }

  private applyLayoutMode(layoutMode: string): void {
    // 移除所有布局类
    Object.values(LAYOUT_CLASSES).forEach(className => {
      this.root.classList.remove(className)
    })
    
    // 添加当前布局类
    const layoutClass = LAYOUT_CLASSES[layoutMode as keyof typeof LAYOUT_CLASSES]
    if (layoutClass) {
      this.root.classList.add(layoutClass)
    }
  }

  watchSystemTheme(callback: () => void): () => void {
    if (typeof window === 'undefined') {
      return () => {}
    }

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    
    const handleChange = () => {
      callback()
    }
    
    mediaQuery.addEventListener('change', handleChange)
    
    // 返回清理函数
    return () => {
      mediaQuery.removeEventListener('change', handleChange)
    }
  }

  getCurrentSystemTheme(): 'light' | 'dark' {
    if (typeof window === 'undefined') {
      return 'light'
    }
    
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }
}
