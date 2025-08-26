package handlers

// Response 通用API响应结构
// @Description API响应的通用结构
type Response struct {
	Success bool        `json:"success" example:"true"`         // 操作是否成功
	Message string      `json:"message" example:"操作成功"`         // 响应消息
	Data    interface{} `json:"data,omitempty"`                 // 响应数据
	Error   string      `json:"error,omitempty" example:"错误信息"` // 错误信息（仅在失败时返回）
}

// PaginatedResponse 分页响应结构
// @Description 分页数据的响应结构
type PaginatedResponse struct {
	Success bool        `json:"success" example:"true"`   // 操作是否成功
	Message string      `json:"message" example:"获取数据成功"` // 响应消息
	Data    interface{} `json:"data"`                     // 响应数据
	Meta    PageMeta    `json:"meta"`                     // 分页元信息
}

// PageMeta 分页元信息
// @Description 分页的元信息
type PageMeta struct {
	CurrentPage int `json:"current_page" example:"1"` // 当前页码
	PerPage     int `json:"per_page" example:"10"`    // 每页条数
	Total       int `json:"total" example:"100"`      // 总记录数
	TotalPages  int `json:"total_pages" example:"10"` // 总页数
}
