# 公共组件目录

## 目录结构
```
components/
├── common/           # 通用组件
│   ├── DataTable/    # 数据表格组件
│   ├── FilterBar/    # 筛选栏组件  
│   ├── StatsCard/    # 统计卡片组件
│   ├── FormDialog/   # 表单对话框组件
│   └── ActionBar/    # 操作栏组件
├── business/         # 业务组件
│   ├── UserManagement/
│   ├── RoleManagement/
│   └── SystemConfig/
└── ui/              # UI基础组件
    ├── LoadingSpinner/
    ├── EmptyState/
    └── ConfirmDialog/
```

## 设计原则

1. **可复用性**: 组件应该是通用的，可以在多个地方使用
2. **可配置性**: 通过props和slots提供灵活的配置选项
3. **一致性**: 保持统一的API设计和视觉风格
4. **文档化**: 每个组件都应该有清晰的使用说明

## 重构指南

### 1. 抽取重复组件
- UserManagement → components/business/UserManagement
- RoleManagement → components/business/RoleManagement  
- SystemConfiguration → components/business/SystemConfig

### 2. 创建通用组件
- DataTable: 统一的数据表格组件
- FilterBar: 统一的筛选条件组件
- StatsCard: 统一的统计卡片组件
- FormDialog: 统一的表单对话框组件

### 3. 组件API设计
每个组件应该有清晰的props、events、slots定义，支持:
- 数据绑定
- 事件回调
- 插槽定制
- 主题适配
