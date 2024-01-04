package common

import (
	// 	"encoding/json"
	"encoding/hex"
	"strconv"
	"time"
	"errors"
	"fmt"
	"reflect"
	"math/rand"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

var (
	timeTemplate1 = "2006-01-02 15:04:05" //常规类型
	timeTemplate2 = "2006/01/02 15:04:05" //其他类型
	timeTemplate3 = "2006-01-02"          //其他类型
	timeTemplate4 = "15:04:05"            //其他类型
	timeTemplate5 = "2006-01-02 15:04"    //其他类型
)

/**
 *@desc 默认返回方式
 *@author Carver
 */
type ReturnJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

/**
 *@desc 用户信息
 *@author Carver
 */
type returnUserData struct {
	UserId    int    `json:"userId,omitempty"`
	UserName  string `json:"userName,omitempty"`
	UserEmail string `json:"userEmail,omitempty"`
	UserToken string `json:"userToken,omitempty"`
	UserTel   string `json:"userTel,omitempty"`
}

/**
 *@desc 文章信息
 *@author Carver
 */
type returnArticleData struct {
	Count int64       `json:"count,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Hot   interface{} `json:"hot,omitempty"`
}

/**
 *@desc 标签或分类信息
 *@author Carver
 */
type returnTagOrCategoryData struct {
	Id   interface{} `json:"id,omitempty"`
	Name interface{} `json:"name,omitempty"`
}

/**
 *@desc 网站信息
 *@author Carver
 */
type returnSiteData struct {
	Info interface{} `json:"info,omitempty"`
}

// 长度为62
var bytes []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

// 设置jwt生成需要的参数
type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

/**
 *@desc 获取随机字符串
 *@author Carver
 */
func RandStr(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}

/**
 *@desc 进行hash加密
 *@author Carver
 */
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

/**
 *@desc 进行hash解密
 *@author Carver
 */
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/**
 *@desc 格式化时间戳
 *@params time_str 需要转化的时间戳字符串
 *@return dataTime (string) 转换后的时间字符串
 *@author Carver
 */
func UnixToTime(time_stamp_str string) (dataTime string) {
	t, _ := strconv.ParseInt(time_stamp_str, 10, 64) //外部传入的时间戳（秒为单位），必须为int64类型
	dataTime = time.Unix(t, 0).Format(timeTemplate5)
	return
}

/**
 *@desc 格式化时间
 *@params time_str 需要转化的时间字符串
 *@return stamp (int64) 转换后的int64时间戳
 *@author Carver
 */
func TimeToUnix(time_str string) time.Time {
	t1 := "2019-01-08 13:50:30"                                     //外部传入的时间字符串
	stamp, _ := time.ParseInLocation(timeTemplate1, t1, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	return stamp
}

func TimeInt64ToInt(int64p int64) int {
	str := strconv.FormatInt(int64p, 10)
	ints, _ := strconv.Atoi(str)
	return ints
}

/**
 *@desc 提取 []struct 中 column 列。desk中存储为 slice
 		提取 []struct 中的 index 列作为 key，column列作为值。 desk 中存储为map
		用到的方法有： 1. StructColumn 2. structColumn 3. findStructValByColumnKey 4. structIndexColumn 5. findStructValByIndexKey
 *@param desk [slice|map] 指针类型，方法最终的存储位置
 *@param input []struct，待转换的结构体切片
 *@param columnKey string
 *@param indexKey string
 *@author Carver
*/
func StructColumn(desk, input interface{}, columnKey, indexKey string) (err error) {
	deskValue := reflect.ValueOf(desk)
	if deskValue.Kind() != reflect.Ptr {
		return errors.New("desk must be ptr")
	}

	rv := reflect.ValueOf(input)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return errors.New("input must be map slice or array")
	}

	rt := reflect.TypeOf(input)
	if rt.Elem().Kind() != reflect.Struct {
		return errors.New("input's elem must be struct")
	}

	if len(indexKey) > 0 {
		return structIndexColumn(desk, input, columnKey, indexKey)
	}
	return structColumn(desk, input, columnKey)
}

func structColumn(desk, input interface{}, columnKey string) (err error) {
	if len(columnKey) == 0 {
		return errors.New("columnKey cannot not be empty")
	}

	deskElemType := reflect.TypeOf(desk).Elem()
	if deskElemType.Kind() != reflect.Slice {
		return errors.New("desk must be slice")
	}

	rv := reflect.ValueOf(input)
	rt := reflect.TypeOf(input)

	var columnVal reflect.Value
	deskValue := reflect.ValueOf(desk)
	direct := reflect.Indirect(deskValue)

	for i := 0; i < rv.Len(); i++ {
		columnVal, err = findStructValByColumnKey(rv.Index(i), rt.Elem(), columnKey)
		if err != nil {
			return
		}
		if deskElemType.Elem().Kind() != columnVal.Kind() {
			return errors.New(fmt.Sprintf("your slice must be []%s", columnVal.Kind()))
		}

		direct.Set(reflect.Append(direct, columnVal))
	}
	return
}

func findStructValByColumnKey(curVal reflect.Value, elemType reflect.Type, columnKey string) (columnVal reflect.Value, err error) {
	columnExist := false
	for i := 0; i < elemType.NumField(); i++ {
		curField := curVal.Field(i)
		if elemType.Field(i).Name == columnKey {
			columnExist = true
			columnVal = curField
			continue
		}
	}
	if !columnExist {
		return columnVal, errors.New(fmt.Sprintf("columnKey %s not found in %s's field", columnKey, elemType))
	}
	return
}

func structIndexColumn(desk, input interface{}, columnKey, indexKey string) (err error) {
	deskValue := reflect.ValueOf(desk)
	if deskValue.Elem().Kind() != reflect.Map {
		return errors.New("desk must be map")
	}
	deskElem := deskValue.Type().Elem()
	if len(columnKey) == 0 && deskElem.Elem().Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("desk's elem expect struct, got %s", deskElem.Elem().Kind()))
	}

	rv := reflect.ValueOf(input)
	rt := reflect.TypeOf(input)
	elemType := rt.Elem()

	var indexVal, columnVal reflect.Value
	direct := reflect.Indirect(deskValue)
	mapReflect := reflect.MakeMap(deskElem)
	deskKey := deskValue.Type().Elem().Key()

	for i := 0; i < rv.Len(); i++ {
		curVal := rv.Index(i)

		indexVal, columnVal, err = findStructValByIndexKey(curVal, elemType, indexKey, columnKey)
		if err != nil {
			return
		}
		if deskKey.Kind() != indexVal.Kind() {
			return errors.New(fmt.Sprintf("cant't convert %s to %s, your map'key must be %s", indexVal.Kind(), deskKey.Kind(), indexVal.Kind()))
		}
		if len(columnKey) == 0 {
			mapReflect.SetMapIndex(indexVal, curVal)
			direct.Set(mapReflect)
		} else {
			if deskElem.Elem().Kind() != columnVal.Kind() {
				return errors.New(fmt.Sprintf("your map must be map[%s]%s", indexVal.Kind(), columnVal.Kind()))
			}
			mapReflect.SetMapIndex(indexVal, columnVal)
			direct.Set(mapReflect)
		}
	}
	return
}

func findStructValByIndexKey(curVal reflect.Value, elemType reflect.Type, indexKey, columnKey string) (indexVal, columnVal reflect.Value, err error) {
	indexExist := false
	columnExist := false

	for i := 0; i < elemType.NumField(); i++ {
		curField := curVal.Field(i)

		if elemType.Field(i).Name == indexKey {

			switch curField.Kind() {
			case reflect.String, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int, reflect.Float64, reflect.Float32:
				indexExist = true
				indexVal = curField
			default:
				return indexVal, columnVal, errors.New("indexKey must be int float or string")
			}
		}
		if elemType.Field(i).Name == columnKey {
			columnExist = true
			columnVal = curField
			continue
		}
	}
	if !indexExist {
		return indexVal, columnVal, errors.New(fmt.Sprintf("indexKey %s not found in %s's field", indexKey, elemType))
	}
	if len(columnKey) > 0 && !columnExist {
		return indexVal, columnVal, errors.New(fmt.Sprintf("columnKey %s not found in %s's field", columnKey, elemType))
	}
	return
}

// 生成token【jwt】
func GenerateToken(username, role, secretKey string, expiration time.Duration) (string, error) {
	// 设置过期时间
	expireAt := time.Now().Add(expiration).Unix()
	// 创建生成的参数
	claims := CustomClaims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt,
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 根据秘钥去签名token
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("生成token失败: %v", err)
	}

	return signedToken, nil
}

// 验证token【jwt】
func VerifyToken(tokenString, secretKey string) (*CustomClaims, error) {
	// 验证token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名方式不正确: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(secretKey), nil
	})

	// token验证失败
	if err != nil {
		return nil, fmt.Errorf("token验证失败: %v", err)
	}

	// 验证通过的token
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的token!")
}
