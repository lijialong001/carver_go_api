package common

/**
 * desc 自定义返回错误码
 * author Carver
 */
const (
	SystemHandleSuccessCode = 200
	SystemHandleErrorCode   = 500

	UserNameCode            = 1001
	UserPwdCode             = 1002
	UserNameDoneCode        = 1003
	UserFormParamsErrorCode = 1004
	UserPwdNoSameErrorCode  = 1005
	ArticlePublishErrorCode = 1006
	TagSaveCode             = 1007
	TagExistCode            = 1008
	CategoryExistCode       = 1009
	CategorySaveCode        = 1010

	UserHandleSuccessCode     = 2001
	UserRegisterSuccessCode   = 2003
	ArticleSuccessCode        = 2002
	ArticlePublishSuccessCode = 2004

	ArticleErrorCode = 5001
)

/**
 * desc 默认的返回信息
 * author Carver
 */
var (
	//操作成功
	HandleSuccess = NewError(0, "操作成功")

	//服务端操作
	ErrServer    = NewError(10001, "服务异常，请联系管理员")
	ErrParam     = NewError(10002, "参数有误")
	ErrSignParam = NewError(10003, "签名参数有误")
)
