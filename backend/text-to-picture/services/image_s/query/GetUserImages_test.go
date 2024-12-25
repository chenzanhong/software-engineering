package query

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	getDB "text-to-picture/config"
	middlewire "text-to-picture/middlewire/jwt"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
	"text-to-picture/models/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/auth/user/images",middlewire.JWTAuthMiddleware(), GetUserImages)
	return r
}

//模拟JWT中间件，主要是不要设置Set，使得Get无法获取username上下文，模拟未找到用户信息
func MockJWTAuthMiddlewareNoUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
    }
}

// DBConfig 结构体定义
type DBConfig struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

// TestGetUserImages 是测试GetUserImages的入口函数
func TestGetUserImages(t *testing.T) {

	//设置gin运行模式为测试模式
	gin.SetMode(gin.TestMode)

	// 读取测试数据库配置
	yamlFile, err := os.ReadFile(getDB.GetDBConfigPath())
	if err != nil {
		fmt.Printf("Error reading DBconfig.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Printf("Error parsing DBconfig.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 设置数据库连接的环境变量
	os.Setenv("DB_USER", dbconfig.DB.User)
	os.Setenv("DB_PASSWORD", dbconfig.DB.Password)
	os.Setenv("DB_HOST", dbconfig.DB.Host)
	os.Setenv("DB_PORT", dbconfig.DB.Port)
	os.Setenv("DB_NAME", dbconfig.DB.Name)

	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	//成功获取用户收藏图片信息的测试
	t.Run("GetUserImages_Success",func(t *testing.T) {

		// 创建gin运行模型为测试模式
		gin.SetMode(gin.TestMode)

		//创建路由
		router := SetupRouter()

		//创建一个数据库尚未存在的用户作为测试用户
		testUserName := "test"

		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: testUserName, 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)

		//为了避免违反外键约束，先插入用户信息
		db.DB.Create(&user.UserInformation{ID:100,UserName: testUserName, Email:testUserName+"@qq.com",
		Password: "123456",Avatar_url:testUserName+".jpg",Score: 100,Token: tokenString,Create_time: time.Now()})

		//插入生成图片信息
		create_time := time.Now().UTC()
		db.DB.Create(&image.ImageInformation{ID:100,UserName: testUserName,Params: "param",LikeCount: 1,
		Picture: testUserName+".jpg",Create_time: create_time,})

		//创建一个GET请求
		request , _ := http.NewRequest("GET","/auth/user/images",nil)
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response,request)

		//检查响应码
		assert.Equal(t,http.StatusOK,response.Code)

		//检查响应体
		expectedResponse := gin.H{
			"message":"获取用户的图像成功",
			"images":[]interface{}{
				map[string]interface {}{
					"id":float64(100),
					"username": testUserName,
					"params": "param",
					"likecount": float64(1),
					"picture": testUserName+".jpg",
					"create_time": create_time.Format("2006-01-02T15:04:05.000000Z"),
					//测试时，有时实际结果的时间是5位小数，导致出错，可以再次运行一次，或者去掉一位小数
				},
			},
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(),&actualResponse)
		fmt.Println("实际响应为：",actualResponse)
		assert.Equal(t,expectedResponse,actualResponse)

		//最后清理测试数据
		db.DB.Where("username = ?", testUserName).Delete(&image.ImageInformation{})
		db.DB.Where("username = ?", testUserName).Delete(&user.UserInformation{})
	})

	//未找到用户信息的测试
	t.Run("GetUserImages_UserNotExist",func(t *testing.T) {

		// 创建gin运行模型为测试模式
		gin.SetMode(gin.TestMode)

		// 创建路由，其中MockJWTAuthMiddlewareNoUser()使GetUserImags的Get无法从上下文获取username
		router := gin.Default()
		router.GET("/auth/user/images", MockJWTAuthMiddlewareNoUser(), GetUserImages)
		
		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: "test", 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)

		//这个测试不涉及数据库，不需要向数据库插入数据

		//创建一个GET请求
		request , _ := http.NewRequest("GET","/auth/user/images",nil)
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response,request)

		//检查响应码
		assert.Equal(t,http.StatusUnauthorized,response.Code)

		//检查响应体
		expectedResponse := gin.H{
			"success": false,
			"message": "未找到用户信息",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(),&actualResponse)
		fmt.Println("实际响应为：",actualResponse)
		assert.Equal(t,expectedResponse,actualResponse)
	})

	//查询用户图片失败的测试(数据库查询失败)
	t.Run("GetUserImages_DatabaseQueryFailed",func(t *testing.T) {

		// 创建gin运行模型为测试模式
		gin.SetMode(gin.TestMode)

		//创建路由
		router := SetupRouter()

		//创建一个数据库尚未存在的用户作为测试用户
		testUserName := "test"

		// 创建有效的Token
		claims := &middlewire.Claims{
			Username: testUserName, 
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(middlewire.JwtKey)

		//为了避免违反外键约束，先插入用户信息
		db.DB.Create(&user.UserInformation{ID:100,UserName: testUserName, Email:testUserName+"@qq.com",
		Password: "123456",Avatar_url:testUserName+".jpg",Score: 100,Token: tokenString,Create_time: time.Now()})

		//插入生成图片信息
		create_time := time.Now().UTC()
		db.DB.Create(&image.ImageInformation{ID:100,UserName: testUserName,Params: "param",LikeCount: 1,
		Picture: testUserName+".jpg",Create_time: create_time,})

		
		//断开数据库连接
		// 获取底层的 sql.DB 对象并关闭连接
		sqlDB, err := db.DB.DB()
		if err != nil {
			t.Fatalf("Failed to get sql.DB object: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("Failed to close database connection: %v", err)
		}

		//创建一个GET请求
		request , _ := http.NewRequest("GET","/auth/user/images",nil)
		request.Header.Set("Authorization", tokenString)

		//创建一个响应器
		response := httptest.NewRecorder()

		//执行请求
		router.ServeHTTP(response,request)

		//检查响应码
		assert.Equal(t,http.StatusInternalServerError,response.Code)

		//检查响应体
		expectedResponse := gin.H{
			"message": "查询用户图片失败", 
			"error": map[string]interface{}{},
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(),&actualResponse)
		fmt.Println("实际响应为：",actualResponse)
		assert.Equal(t,expectedResponse,actualResponse)

		// 重新建立数据库连接
		if err := db.ConnectDatabase(); err != nil {
			t.Fatalf("Failed to reconnect to test database: %v", err)
		}

		//最后清理测试数据
		db.DB.Where("username = ?", testUserName).Delete(&image.ImageInformation{})
		db.DB.Where("username = ?", testUserName).Delete(&user.UserInformation{})
	})

}