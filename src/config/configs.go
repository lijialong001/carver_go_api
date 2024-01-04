package config

import (
	"database/sql"
	. "fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

type TomlConfig struct {
	Title   string
	Admin   adminInfo
	Db      database  `toml:"database"`
	Redis   redisInfo `toml:"redis"`
	Servers map[string]server
	Other   other
}

type adminInfo struct {
	Name     string
	Port     int `toml:"ListenPort"`
	Debug    bool
	LogLevel string
}

type database struct {
	UserName string `toml:"username"`
	Password string
	Database string
	Addr     string
	Ports    []int
	ConnMax  int `toml:"connection_max"`
	Enabled  bool
	Port     int `toml:"system_port"`
}

type redisInfo struct {
	Addr  string
	Auth  string
	Dial  string `toml:"DialTimeout"`
	Read  string `toml:"ReadTimeout"`
	Write string `toml:"WriteTimeout"`
}

type server struct {
	Db  string
	Pwd string
}

type other struct {
	Data      [][]interface{}
	Type      []string
	SecretKey string
}

var DB *sql.DB

var Pool *redis.Pool //创建redis连接池

/**
 *@desc 获取mysql数据库连接
 *@author Carver
 */
func GetInitMysqlConnect() {

	pwd, pwdError := os.Getwd()
	if pwdError != nil {
		os.Exit(1)
		panic(pwdError)
	}

	envUrl := pwd + "/src/config/local.config.toml"

	var config TomlConfig
	_, configError := toml.DecodeFile(envUrl, &config)

	if configError != nil {
		panic(configError)
	}

	connStr := Sprintf("%s:%s@(%s)/%s", config.Db.UserName, config.Db.Password, config.Db.Addr, config.Db.Database)
	//设置连接数据库的参数
	db, _ := sql.Open("mysql", connStr)
	// 设置连接池参数
	db.SetMaxOpenConns(10) // 设置最大连接数
	db.SetMaxIdleConns(5)  // 设置最大空闲连接数
	//连接数据库
	mysqlConnectErr := db.Ping()
	if mysqlConnectErr != nil {
		Println("数据库连接失败❌ ")
		return
	} else {
		Println("数据库连接成功✅ ")
	}
	DB = db
}

/**
 *@desc 获取redis数据库连接
 *@author Carver
 */
func GetInitRedisConnect() {

	pwd, pwdError := os.Getwd()
	if pwdError != nil {
		os.Exit(1)
		panic(pwdError)
	}

	envUrl := pwd + "/src/config/local.config.toml"

	var config TomlConfig
	_, configError := toml.DecodeFile(envUrl, &config)

	if configError != nil {
		panic(configError)
	}

	Pool = &redis.Pool{ //实例化一个连接池
		MaxIdle: 16, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			RedisConnect, _ := redis.Dial("tcp", config.Redis.Addr)
			// 设置密码认证
			_, redisConnectErr := RedisConnect.Do("AUTH", config.Redis.Auth)
			if redisConnectErr != nil {
				log.Fatal(redisConnectErr)
			}
			return RedisConnect, redisConnectErr
		},
	}
}
