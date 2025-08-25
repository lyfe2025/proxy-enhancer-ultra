// 系统管理API - 重新导出拆分后的模块以保持向后兼容性

// 重新导出类型定义
export type {
  User,
  Role,
  Permission,
  SystemConfig,
  OperationLog,
  SystemStats,
  SystemHealth,
  FileUploadResponse,
  BackupInfo,
  UserFormData,
  RoleFormData,
  PermissionFormData
} from './system-types'

// 重新导出API函数
export * from './user-api'
export * from './role-api'
export * from './permission-api'
export * from './config-api'
export * from './log-api'
export * from './stats-api'
export * from './file-api'

// 重新导出API对象（保持向后兼容）
export { userApi } from './user-api'
export { roleApi } from './role-api'
export { permissionApi } from './permission-api'
export { configApi as systemApi } from './config-api'
export { logApi } from './log-api'
export { statsApi } from './stats-api'
export { fileApi } from './file-api'

// 兼容性导出（旧的API对象名称）
export { userApi as systemUserApi } from './user-api'
export { roleApi as systemRoleApi } from './role-api'
export { permissionApi as systemPermissionApi } from './permission-api'