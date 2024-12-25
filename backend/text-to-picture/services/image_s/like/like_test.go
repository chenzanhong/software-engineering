package like

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	middlewire "text-to-picture/middlewire/jwt"
	"text-to-picture/models/image"
	db "text-to-picture/models/init"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLikeImage(t *testing.T) {
	// 初始化 Gin 引擎
	r := gin.Default()
	r.Use(middlewire.JWTAuthMiddleware()) // 假设这是你的 JWT 中间件
	r.POST("/auth/like", LikeImage)

	// 模拟数据库
	db.DB = MockDB()

	defer func() {
		// 清理表中的数据，确保每次测试都是独立的
		if db.DB.Migrator().HasTable(&image.ImageLike{}) {
			db.DB.Where("username = ? OR picture = ?", "test_user", "http://example.com/image.jpg").Delete(&image.ImageLike{})
		}
		if db.DB.Migrator().HasTable(&image.ImageInformation{}) {
			db.DB.Where("picture = ?", "http://example.com/image.jpg").Delete(&image.ImageInformation{})
		}
	}()

	// 测试用例
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		username       string
		mockDBSetup    func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "无效的token",
			requestBody: map[string]interface{}{"url": "http://example.com/image.jpg"},
			username:    "test_user",

			mockDBSetup:    func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"code": float64(401), "message": "无效的Token"},
		},
		{
			name:        "成功点赞",
			requestBody: map[string]interface{}{"url": "http://example.com/image.jpg"},
			username:    "test_user",

			mockDBSetup: func() {
				// 插入初始数据，如图片信息
				imgInfo := image.ImageInformation{Picture: "http://example.com/image.jpg", UserName: "test_user", LikeCount: 0}
				db.DB.Create(&imgInfo)
			},
			expectedStatus: http.StatusOK, //在JSON解码的过程中，Go 的 encoding/json 默认会将数字解析为 float64
			expectedBody:   map[string]interface{}{"code":float64(200), "current_likes": float64(1), "message": "Image liked successfully"},
		},
		{
			name:           "缺少图片 URL",
			requestBody:    map[string]interface{}{},
			username:       "test_user",
			mockDBSetup:    func() {}, // 不需要数据库操作
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"code": float64(400), "error": "Missing image URL"},
		},
		{
			name:        "用户已点赞",
			requestBody: map[string]interface{}{"url": "http://example.com/image.jpg"},
			username:    "test_user",
			mockDBSetup: func() {
				// 初始化用户点赞记录
				db.DB.Create(&image.ImageLike{Picture: "http://example.com/image.jpg", UserName: "test_user", Num: 1})
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   map[string]interface{}{"code": float64(409), "error": "用户已经点赞过该图片"},
		},
		{
			name:           "获取用户点赞数失败",
			requestBody:    map[string]interface{}{"url": "http://example.com/invalidImage.jpg"},
			username:       "test_user",
			mockDBSetup:    func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"code": float64(500), "error": "sql: no rows in result set"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 初始化数据库状态
			if test.mockDBSetup != nil {
				test.mockDBSetup()
			}
			// 模拟请求
			body, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest("POST", "/auth/like", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			var tokenString string
			var err error
			if test.name != "无效的token" {
				// 创建一个有效的Token
				claims := &middlewire.Claims{
					Username: "test_user",
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // 设置有效的过期时间
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err = token.SignedString(middlewire.JwtKey)
				if err != nil {
					t.Fatalf("Failed to create token: %v", err)
				}
			} else {
				tokenString = "invaildToken"
			}
			req.Header.Set("Authorization", tokenString) // 模拟 Token

			// // 模拟上下文
			// w := httptest.NewRecorder()
			// c, _ := gin.CreateTestContext(w)
			// c.Request = req
			// c.Set("username", test.username) // 模拟上下文中解析到的用户名
			// 执行路由，不再手动创建 Gin 上下文
			w := httptest.NewRecorder()

			// 执行路由
			r.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, test.expectedStatus, w.Code)
			fmt.Println("Response Body:", string(w.Body.Bytes()))
			var responseMap map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &responseMap); err != nil {
				t.Errorf("failed to unmarshal response body: %v", err)
			}
			assert.Equal(t, test.expectedBody, responseMap)
		})
	}
}

// MockDB 模拟数据库操作
func MockDB() *gorm.DB {
	// 使用内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 执行自动迁移并确认表已创建
	if err := db.AutoMigrate(&image.ImageLike{}, &image.ImageInformation{}); err != nil {
		panic("failed to migrate database tables")
	}

	// 检查表是否存在
	if !db.Migrator().HasTable(&image.ImageLike{}) {
		panic("image_likes table does not exist after migration")
	}
	if !db.Migrator().HasTable(&image.ImageInformation{}) {
		panic("image_information table does not exist after migration")
	}

	return db
}
