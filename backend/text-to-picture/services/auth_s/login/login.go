package auth_s

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	middlewire "text-to-picture/middlewire/jwt" // JWT 认证中间件
	models "text-to-picture/models/init"        // 数据库模型初始化
	"text-to-picture/models/repository/user_r"  // 用户数据访问层
	userLogin "text-to-picture/models/user"     // 用户模型
	"time"

	"github.com/dgrijalva/jwt-go" // JWT 库
	"github.com/gin-gonic/gin"    // Gin 框架
	_ "github.com/lib/pq"         // PostgreSQL 驱动
	"gorm.io/gorm"                // GORM ORM
)

// 注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags auth
// @Accept json
// @Produce json
// @Param user body userLogin.UserInformation true "用户信息"
// @Success 200 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]interface{} "请求数据格式错误"
// @Failure 500 {object} map[string]interface{} "用户创建失败"
// @Router /register [post]
func Register(c *gin.Context) {
	// 定义用于接收 JSON 数据的结构体
	var input userLogin.UserInformation

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		// 如果解析失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求数据格式错误"})
		return
	}
	input.Avatar_url = "https://chuhsing-blog-bucket.oss-cn-shenzhen.aliyuncs.com.aliyuncs.com/chuhsing/202412092143859.png"
	//插入数据
	if err := user_r.InsertUserInformation(models.DB, &input); err != nil {
		// 如果插入失败，返回 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "用户创建失败",
			"error":   err.Error()})
		return
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "注册成功",
	})
}

// 登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags auth
// @Accept json
// @Produce json
// @Param user body userLogin.UserInformation true "用户名和密码"
// @Success 200  {object} map[string]interface{} "登录成功"
// @Failure 400  {object} map[string]interface{} "请求数据格式错误"
// @Failure 401  {object} map[string]interface{} "用户不存在或密码错误"
// @Failure 500 {object} map[string]interface{} "生成 token 错误"
// @Router /login [post]
func Login(c *gin.Context) {
	// 定义用于接收 JSON 数据的结构体
	var input struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		// 如果解析失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请求数据格式错误"})
		return
	}
	fmt.Println(input.Password)

	// 查找用户
	user, err := user_r.GetUserByName(models.DB, input.Name)
	if err != nil {
		log.Printf("name:%v,%v", input.Name, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，返回 401 错误
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "用户不存在"})
			return
		}
		// 返回查询失败信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "数据库查询错误"})
		return
	}

	// 验证密码将用户的明文密码与数据库中的哈希密码对比
	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "密码错误"})
		return
	}

	//生成 JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middlewire.Claims{
		Username: input.Name, // 基于输入的用户名
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middlewire.JwtKey)
	if err != nil {
		// 如果生成 Token 失败，返回 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "生成 token 错误"})
		return
	}

	// 更新用户的 token
	updates := map[string]interface{}{
		"token": tokenString,
	}
	if err := user_r.UpdateUserInfo(models.DB, user.UserName, updates); err != nil {
		// 如果更新失败，返回 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "登录时更新用户 token 失败", "error": err.Error()})
		return
	}

	// 登录成功，返回用户信息和 Token
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登录成功",
		"token":   tokenString,
		"avatar":  user.Avatar_url, // 前端登录后及时刷新头像需要用到
	})
}
