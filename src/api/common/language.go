package common

var AllZhTips = map[int]string{
	//系统操作模块
	SystemHandleSuccessCode: "操作成功",
	SystemHandleErrorCode:   "操作失败",
	//用户操作模块
	UserHandleSuccessCode:   "登录成功!",
	UserNameCode:            "用户名不存在!",
	UserPwdCode:             "密码不正确!",
	UserNameDoneCode:        "用户名已存在!",
	UserRegisterSuccessCode: "注册成功!",
	UserFormParamsErrorCode: "用户名 | 密码 | 确认密码缺一不可哦!",
	UserPwdNoSameErrorCode:  "密码和确认密码不正确!",
	//文章操作模块
	ArticleErrorCode:          "请求异常!",
	ArticleSuccessCode:        "获取成功!",
	ArticlePublishSuccessCode: "发布成功!",
	ArticlePublishErrorCode:   "发布失败!",
	//标签模块
	TagSaveCode:  "该标签已经关联过文章，不能进行其它操作!",
	TagExistCode: "该标签已经存在!",
	//分类模块
	CategoryExistCode: "该分类已经存在!",
	CategorySaveCode:  "该分类已经关联过文章，不能进行其它操作!",
}
