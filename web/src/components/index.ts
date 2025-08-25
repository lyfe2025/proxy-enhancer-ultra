// 公共组件导出

// 通用组件
export { default as DataTable } from './common/DataTable/DataTable.vue'

// 业务组件
export { default as UserManagement } from './business/UserManagement/UserManagement.vue'
export { default as UserFormDialog } from './business/UserManagement/UserFormDialog.vue'
export { default as RoleManagement } from './business/RoleManagement/RoleManagement.vue'
export { default as RoleFormDialog } from './business/RoleManagement/RoleFormDialog.vue'
export { default as SystemConfiguration } from './business/SystemConfiguration/SystemConfiguration.vue'

// 组件类型导出
export type { TableColumn, PaginationData } from './common/DataTable/DataTable.vue'
