import type { FontSize, BorderRadius } from './theme-types'
export { DEFAULT_THEME_CONFIG } from './theme-types'

export const FONT_SIZE_MAP: Record<FontSize, string> = {
  small: '14px',
  medium: '16px',
  large: '18px'
}

export const BORDER_RADIUS_MAP: Record<BorderRadius, string> = {
  none: '0px',
  small: '4px',
  medium: '8px',
  large: '12px'
}

export const THEME_STORAGE_KEY = 'theme-settings'

export const PRESET_COLORS = [
  '#00ff88', // 默认亮绿色
  '#409eff', // Element Plus 蓝色
  '#67c23a', // 绿色
  '#e6a23c', // 橙色
  '#f56c6c', // 红色
  '#909399', // 灰色
  '#9c88ff', // 紫色
  '#ff88c2', // 粉色
  '#88d4ff', // 天蓝色
  '#88ffaa'  // 薄荷绿
]

export const LAYOUT_CLASSES = {
  default: 'layout-default',
  compact: 'layout-compact',
  comfortable: 'layout-comfortable'
}

export const THEME_CLASSES = {
  dark: 'dark',
  light: 'light',
  noAnimations: 'no-animations',
  compactMode: 'compact-mode'
}
