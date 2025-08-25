export type ThemeMode = 'light' | 'dark' | 'auto'
export type SidebarMode = 'expanded' | 'collapsed'
export type LayoutMode = 'default' | 'compact' | 'comfortable'
export type FontSize = 'small' | 'medium' | 'large'
export type BorderRadius = 'none' | 'small' | 'medium' | 'large'

export interface ThemeConfig {
  isDark: boolean
  themeMode: ThemeMode
  sidebarCollapsed: boolean
  layoutMode: LayoutMode
  primaryColor: string
  fontSize: FontSize
  borderRadius: BorderRadius
  animationEnabled: boolean
  compactMode: boolean
  showBreadcrumb: boolean
  showTabsView: boolean
  fixedHeader: boolean
  fixedSidebar: boolean
}

export const DEFAULT_THEME_CONFIG: ThemeConfig = {
  isDark: true,
  themeMode: 'dark',
  sidebarCollapsed: false,
  layoutMode: 'default',
  primaryColor: '#00ff88',
  fontSize: 'medium',
  borderRadius: 'medium',
  animationEnabled: true,
  compactMode: false,
  showBreadcrumb: true,
  showTabsView: true,
  fixedHeader: true,
  fixedSidebar: true
}
