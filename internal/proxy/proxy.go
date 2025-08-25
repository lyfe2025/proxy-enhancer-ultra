package proxy

// 代理服务器包 - 重新导出拆分后的组件以保持向后兼容

// 注意：这个文件现在只作为重新导出，实际实现已经拆分到各个专门的文件中

// 你可以直接从各个专门的文件中使用组件：
// proxy_server.go - 代理服务器核心功能
// request_processor.go - HTTP请求/响应处理
// html_injector.go - HTML内容注入
// url_rewriter.go - URL重写功能
// html_parser.go - HTML节点操作
// assets_provider.go - 静态资源提供

// 为了向后兼容，重新导出主要的类型
// 但建议直接使用拆分后的专门组件
