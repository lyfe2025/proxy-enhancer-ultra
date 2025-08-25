// 系统监控相关类型定义

// 系统信息类型
export interface SystemInfo {
  version: string
  buildTime: string
  gitCommit: string
  environment: string
  uptime: number
  server: {
    os: string
    arch: string
    platform: string
    hostname: string
    cpuCount: number
    totalMemory: number
    nodeVersion: string
  }
  database: {
    type: string
    version: string
    host: string
    port: number
    name: string
    status: 'connected' | 'disconnected'
  }
  cache: {
    type: string
    status: 'connected' | 'disconnected'
    memory?: {
      used: number
      total: number
    }
  }
}

// 系统状态类型
export interface SystemStatus {
  status: 'healthy' | 'warning' | 'error'
  services: {
    name: string
    status: 'running' | 'stopped' | 'error'
    uptime: number
    memory: number
    cpu: number
  }[]
  metrics: {
    cpu: {
      usage: number
      cores: number
    }
    memory: {
      used: number
      total: number
      usage: number
    }
    disk: {
      used: number
      total: number
      usage: number
    }
    network: {
      bytesIn: number
      bytesOut: number
      packetsIn: number
      packetsOut: number
    }
  }
  database: {
    connections: {
      active: number
      idle: number
      total: number
    }
    queries: {
      total: number
      slow: number
      failed: number
    }
    size: number
  }
}

// 系统服务类型
export interface SystemService {
  name: string
  displayName: string
  description: string
  status: 'running' | 'stopped' | 'error'
  pid?: number
  uptime: number
  memory: number
  cpu: number
  autoStart: boolean
  restartCount: number
  lastRestart?: string
}

// 系统指标类型
export interface SystemMetrics {
  timestamp: string
  cpu: {
    usage: number
    load: number[]
  }
  memory: {
    used: number
    total: number
    usage: number
    swap: {
      used: number
      total: number
    }
  }
  disk: {
    used: number
    total: number
    usage: number
    io: {
      read: number
      write: number
    }
  }
  network: {
    bytesIn: number
    bytesOut: number
    packetsIn: number
    packetsOut: number
    errors: number
  }
  database: {
    connections: number
    queries: number
    slowQueries: number
    size: number
  }
  application: {
    requests: number
    errors: number
    responseTime: number
    activeUsers: number
  }
}
