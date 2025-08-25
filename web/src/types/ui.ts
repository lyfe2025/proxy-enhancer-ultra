// UI组件相关类型定义

// 表格列配置类型
export interface TableColumn {
  prop: string
  label: string
  width?: string | number
  minWidth?: string | number
  fixed?: boolean | 'left' | 'right'
  sortable?: boolean
  filterable?: boolean
  searchable?: boolean
  formatter?: (row: any, column: any, cellValue: any, index: number) => string
  render?: (h: any, params: any) => any
}

// 菜单项类型
export interface MenuItem {
  id: string
  name: string
  path: string
  icon?: string
  component?: string
  redirect?: string
  meta?: {
    title: string
    requireAuth?: boolean
    permissions?: string[]
    roles?: string[]
    hidden?: boolean
    cache?: boolean
    affix?: boolean
    breadcrumb?: boolean
  }
  children?: MenuItem[]
}

// 路由配置类型
export interface RouteConfig {
  path: string
  name: string
  component: any
  redirect?: string
  meta?: {
    title: string
    requireAuth?: boolean
    permissions?: string[]
    roles?: string[]
    hidden?: boolean
    cache?: boolean
    affix?: boolean
    breadcrumb?: boolean
    icon?: string
  }
  children?: RouteConfig[]
}
