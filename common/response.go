package common

/*
错误码长度4位，0表示成功
*/

type Code int64

const (
	RequestSuccess Code = 0
)

type Response struct {
	Code    int64       `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应描述
	Error   string      `json:"error"`   // 错误信息
	Data    interface{} `json:"data"`    // 返回数据
}

type ResponseList struct {
	Code    int64  `json:"code"`    // 响应码
	Message string `json:"message"` // 响应描述
	Error   string `json:"error"`   // 错误信息
	Data    List   `json:"data"`    // 返回数据
}

type List struct {
	Index int64         `json:"index"` // 页码
	Limit int64         `json:"limit"` // 每页记录数
	Page  int64         `json:"page"`  // 页数
	Total int64         `json:"total"` // 全部记录数
	List  []interface{} `json:"list"`  // 记录列表
}
