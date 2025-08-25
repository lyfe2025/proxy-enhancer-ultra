/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface ImportMetaEnv {
  // API 配置
  readonly VITE_API_BASE_URL: string
  
  // 应用配置
  readonly VITE_APP_TITLE: string
  readonly VITE_APP_VERSION: string
  
  // 开发服务器配置
  readonly VITE_DEV_PORT: string
  
  // 功能开关
  readonly VITE_DEBUG: string
  readonly VITE_DEV_MODE: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}