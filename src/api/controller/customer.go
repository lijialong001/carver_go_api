package controller

import (
	"api/service"
	"config"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "log"
	"net/http"
)

func GetCustomerInfo(w http.ResponseWriter, r *http.Request) {
	//获取service逻辑
	customer := service.GetCustomerInfo()
	// 自定义返回的JSON格式
	jsonData, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func GetRedisInfo(w http.ResponseWriter, r *http.Request) {
	//初始化数据库
	config.GetInitRedisConnect()

	c := config.Pool.Get() //从连接池，取一个链接

	defer c.Close() //函数运行结束 ，把连接放回连接池

	//选择数据库
	_, selectErr := c.Do("SELECT", 1)
	//操作数据库
	_, err := c.Do("Set", "test", 200)
	if err != nil {
		fmt.Println(err)
		return
	}
	if selectErr != nil {
		fmt.Println(selectErr)
		return
	}

	req, err := redis.Int(c.Do("Get", "test"))
	if err != nil {
		fmt.Println("get test faild :", err)
		return
	}
	fmt.Println(req)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	//获取service逻辑
	customer := service.UserLoginServices(username, password)
	// 自定义返回的JSON格式
	jsonData, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
