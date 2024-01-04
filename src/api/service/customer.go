package service

import (
	"api/common"
	"api/model"
	"config"
	"github.com/BurntSushi/toml"
	"os"
	"time"

	//_ "encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetCustomerInfo() common.ReturnJson {
	//初始化数据库
	config.GetInitMysqlConnect()
	var returnData common.ReturnJson
	var CustomerField model.CustomerField
	db := config.DB
	defer db.Close() //函数运行结束 ，把连接放回连接池

	returnData.Code = 200
	returnData.Msg = "请求成功！"

	customerQueryInfo, errs := db.Query("select id,module,secret from nwp_dingding_robot_config where robot_name = ? order by id desc limit ?", 1, 5)
	if errs != nil {
		fmt.Println("读取响应失败:", errs)
		return returnData
	}

	var customerList []model.CustomerField
	for customerQueryInfo.Next() {
		customerQueryInfo.Scan(&CustomerField.Id, &CustomerField.UserName, &CustomerField.Role)
		customerList = append(customerList, CustomerField)
	}

	returnData.Data = customerList
	return returnData
}

// 用户登录
func UserLoginServices(username string, password string) common.ReturnJson {
	//初始化数据库
	config.GetInitMysqlConnect()
	var returnData common.ReturnJson
	var CustomerField model.CustomerField

	db := config.DB
	defer db.Close() //函数运行结束 ，把连接放回连接池

	returnData.Code = 200
	returnData.Msg = "请求成功！"

	row, _ := db.Query("select id,username,password,role from users where username = ?", username)
	for row.Next() {
		userNameErr := row.Scan(&CustomerField.Id, &CustomerField.UserName, &CustomerField.Password, &CustomerField.Role)

		if userNameErr != nil {
			returnData.Code = common.UserNameCode
			returnData.Msg = common.AllZhTips[common.UserNameCode]
			return returnData
		} else {
			//验证用户密码
			match := common.PasswordVerify(password, CustomerField.Password)
			if match == false {
				returnData.Code = common.UserPwdCode
				returnData.Msg = common.AllZhTips[common.UserPwdCode]
			} else {
				returnData.Code = common.UserHandleSuccessCode
				returnData.Msg = common.AllZhTips[common.UserHandleSuccessCode]
				//获取秘钥
				pwd, pwdError := os.Getwd()
				if pwdError != nil {
					os.Exit(1)
					fmt.Println(pwdError)
					panic(pwdError)
				}

				envUrl := pwd + "/src/config/local.config.toml"
				var configInfo config.TomlConfig
				_, configError := toml.DecodeFile(envUrl, &configInfo)
				if configError != nil {
					fmt.Println(pwdError)
					panic(configError)
				}

				//生成并保存token
				userToken, _ := common.GenerateToken(CustomerField.UserName, CustomerField.Role, configInfo.Other.SecretKey, 2*time.Hour)

				//初始化数据库
				config.GetInitRedisConnect()

				c := config.Pool.Get() //从连接池，取一个链接

				defer c.Close() //函数运行结束 ，把连接放回连接池

				//选择数据库
				_, selectErr := c.Do("SELECT", 1)
				if selectErr != nil {
					fmt.Println(selectErr)
					panic(selectErr)
				}
				// 定义要设置的字段和值的映射关系
				fieldsAndValues := map[string]interface{}{
					"id":       CustomerField.Id,
					"username": CustomerField.UserName,
					"role":     CustomerField.Role,
					"token":    userToken,
				}

				hashName := CustomerField.Id
				// 构建命令参数列表
				args := []interface{}{hashName}
				for field, value := range fieldsAndValues {
					args = append(args, field, value)
				}
				//操作数据库
				_, actionErr := c.Do("HMSET", args...)
				if actionErr != nil {
					fmt.Println(actionErr)
					panic(actionErr)
				}

				//返回用户信息
				return_user_info := model.ReturnUserData{
					Id:       CustomerField.Id,
					UserName: CustomerField.UserName,
					Token:    userToken,
				}
				returnData.Data = return_user_info
			}
		}
	}
	return returnData
}
