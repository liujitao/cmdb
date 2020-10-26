package common

type Code int64

const (
	RequestSuccess             Code = 0
	RequestOtherError          Code = 4000
	RequestKeyNotFound         Code = 4001
	RequestParameterTypeError  Code = 4002
	RequestAuthorizedFailed    Code = 4003
	RequestNotFound            Code = 4004
	RequestParameterMiss       Code = 4005
	RequestMethodNotAllowed    Code = 4006
	RequestExpired             Code = 4007
	RequestAccessDeny          Code = 4009
	RequestParameterRangeError Code = 4002
)

type Response struct {
	Code    int64       `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应描述
	Data    interface{} `json:"data"`    // 返回数据
}

type ResponseList struct {
	Code    int64  `json:"code"`    // 响应码
	Message string `json:"message"` // 响应描述
	Data    List   `json:"data"`    // 返回数据
}

type List struct {
	Index int64         `json:"index"` // 页码
	Size  int64         `json:"size"`  // 大小
	Page  int64         `json:"page"`  // 页数
	Total int64         `json:"total"` // 总记录数
	List  []interface{} `json:"list"`
}
