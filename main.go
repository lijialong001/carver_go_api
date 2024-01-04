package main

import (
	"config"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"api/controller"
	//"github.com/garyburd/redigo/redis"
	"api/common"
	"os"
)

var AuthInfo struct {
	UserId string `json:"id"`
	//UserName string `json:"username"`
	//Role     string `json:"role"`
}

func main() {

	r := mux.NewRouter()

	//日志中间件【全局】
	r.Use(loggingMiddleware)

	// 首页
	r.HandleFunc("/", homeHandler)

	/****************************************************【不需要登录模块】****************************************************/
	//TODO 用户模块
	//用户登录
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", controller.UserLogin)

	/****************************************************【需要登录模块】****************************************************/
	//TODO 测试模块
	//测试1
	testRouter := r.PathPrefix("/test").Subrouter()
	testRouter.HandleFunc("/test1", controller.GetCustomerInfo)
	testRouter.HandleFunc("/test2", controller.GetRedisInfo)

	//注册登录中间件
	testRouter.Use(authMiddleware)

	// 启动服务器
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// 路由处理函数
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the home page!")
}

// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: action:%s url:%s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// 登录中间件
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 验证token
		token := r.Header.Get("_token")
		payload, _ := ioutil.ReadAll(r.Body)
		userInfo := string(payload)
		err := json.Unmarshal([]byte(userInfo), &AuthInfo)
		if err != nil {
			fmt.Fprintln(w, "授权失败:", err)
			return
		}

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

		_, verifyError := common.VerifyToken(token, configInfo.Other.SecretKey)
		if verifyError != nil {
			fmt.Fprintln(w, verifyError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
