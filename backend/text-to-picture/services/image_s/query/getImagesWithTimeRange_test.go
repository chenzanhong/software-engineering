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

// SetupRouter2 设置 Gin 路由
func SetupRouter2() *gin.Engine {
	r := gin.Default()
	r.GET("/auth/user/images/timeRange", middlewire.JWTAuthMiddleware(), GetImagesWithinTimeRange)
	return r
}

// Helper function to create a JWT token
func createToken(username string) string {
	claims := &middlewire.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middlewire.JwtKey)
	return tokenString
}

// Helper function to perform a GET request
func performGetRequest(router *gin.Engine, tokenString, startTime, endTime string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("GET", fmt.Sprintf("/auth/user/images/timeRange?start_time=%s&end_time=%s", startTime, endTime), nil)
	request.Header.Set("Authorization", tokenString)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

// TestGetImagesWithinTimeRange 是测试 GetImagesWithinTimeRange 函数的入口函数
func TestGetImagesWithinTimeRange(t *testing.T) {
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

	// 清理数据库的函数
	defer func() {
		db.DB.Where("username = ?", "test").Delete(&user.UserInformation{})
		db.DB.Where("username = ?", "test").Delete(&image.ImageInformation{})
	}()

	// 创建路由
	router := SetupRouter2()

	// 测试缺少用户名
	t.Run("Missing_Username", func(t *testing.T) {
		router1 := gin.Default()
		router1.GET("/auth/user/images/timeRange", MockJWTAuthMiddlewareNoUser(), GetImagesWithinTimeRange)

		tokenString := createToken("test")

		response := performGetRequest(router1, tokenString, "2023-10-01", "2024-10-31")

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		expectedResponse := gin.H{
			"code":    float64(401),
			"message": "未找到用户信息",
		}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	// 测试未找到图像
	t.Run("No_Images_Found", func(t *testing.T) {
		testUsername := "test"
		tokenString := createToken(testUsername)

		response := performGetRequest(router, tokenString, "2023-10-01T00:00:00Z", "2024-10-31T23:59:59Z")

		assert.Equal(t, http.StatusOK, response.Code)

		expectedResponse := gin.H{"code": float64(200), "message": "查询图像列表成功", "images": []interface{}{}}
		var actualResponse gin.H
		json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	// 测试无效的开始时间
	t.Run("Inalid_StartTime_Range", func(t *testing.T) {
		testUsername := "test"
		tokenString := createToken(testUsername)

		response := performGetRequest(router, tokenString, "2023.10.1", "2050-12-31T23:59:59.000000Z")

		assert.Equal(t, http.StatusBadRequest, response.Code)

		var actualResponse gin.H
		err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		fmt.Println("实际的响应是：",actualResponse)
		assert.Equal(t, float64(400), actualResponse["code"])
		assert.Equal(t, "无效的开始时间格式", actualResponse["message"])
	})

	// 测试无效的结束时间
	t.Run("Inalid_EndTime_Range", func(t *testing.T) {
		testUsername := "test"
		tokenString := createToken(testUsername)

		response := performGetRequest(router, tokenString, "2023-10-01T00:00:00.000000Z", "2050.10.1")

		assert.Equal(t, http.StatusBadRequest, response.Code)

		var actualResponse gin.H
		err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		fmt.Println("实际的响应是：",actualResponse)
		assert.Equal(t, float64(400), actualResponse["code"])
		assert.Equal(t, "无效的结束时间格式", actualResponse["message"])
	})

	// 测试有效的时间范围，成功查询
	t.Run("Valid_Time_Range", func(t *testing.T) {
		testUsername := "test_Valid_Time_Range"
		create_time := time.Now()
		db.DB.Create(&user.UserInformation{ID: 100, UserName: testUsername, Email: testUsername + "@qq.com", Password: "123456", Create_time: create_time})

		tokenString := createToken(testUsername)

		imageInfo := &image.ImageInformation{
			UserName:   testUsername,
			Picture:    "http://example.com/test_Valid_Time_Range.jpg",
			Create_time: time.Now(),
		}
		db.DB.Create(&imageInfo)

		response := performGetRequest(router, tokenString, "2023-10-01", "2050-12-31")

		assert.Equal(t, http.StatusOK, response.Code)

		var actualResponse gin.H
		err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		fmt.Println("实际的响应是：",actualResponse)
		assert.Equal(t, float64(200), actualResponse["code"])
		assert.Equal(t, "查询图像列表成功", actualResponse["message"])
		assert.NotEmpty(t, actualResponse["images"].([]interface{}))

		// 清理测试数据
		db.DB.Where("username = ?", testUsername).Delete(&image.ImageInformation{})
		db.DB.Where("username = ?", testUsername).Delete(&user.UserInformation{})
	})
}
