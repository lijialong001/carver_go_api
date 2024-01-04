package common

import (
	"encoding/json"
)

/**
 *@desc 格式化返回的JSON数据
 *@author Carver
 */
type errInfo struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 返回的数据
}

/**
 *@desc 设置系统默认的状态码和提示信息
 *@author Carver
 */
func NewError(code int, msg string) Error {
	return &errInfo{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

/**
 *@desc 添加返回的数据
 *@author Carver
 */
func (e *errInfo) JsonWithData(data interface{}) Error {
	e.Data = data
	return e
}

/**
 *@desc 设置返回JSON 格式
 *@author Carver
 */
func (e *errInfo) ToJson() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
		ID   string      `json:"id"`
	}{
		Code: e.Code,
		Msg:  e.Msg,
		Data: e.Data,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}

var _ Error = (*errInfo)(nil)

/**
 *@desc 格式化接口的统一方法
 *@author Carver
 */
type Error interface {
	//添加返回的数据
	JsonWithData(data interface{}) Error
	//设置返回JSON 格式
	ToJson() string
}
