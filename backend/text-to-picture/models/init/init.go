// backend\text-to-picture\models\init\init.go
package models

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	image2 "text-to-picture/models/image"
	user2 "text-to-picture/models/user"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SQL语句用于创建相关的数据库表，使用postgresSQL
const createTableSQL = `
CREATE TABLE IF NOT EXISTS UserInformation (
    id SERIAL PRIMARY KEY,                       -- 用户ID，自动递增
	email VARCHAR(50) UNIQUE NOT NULL,           -- 用户邮箱，唯一且不能为空
	username VARCHAR(30) UNIQUE NOT NULL,        -- 用户名，唯一且不能为空
    password VARCHAR(256) NOT NULL,              -- 用户密码，不能为空
	avatar_url VARCHAR(255) NOT NULL,            -- 头像图片的URL，不能为空
    score INT DEFAULT 0,                         -- 用户积分，默认值为0
	create_time TIMESTAMP DEFAULT NOW(),         -- 创建时间，默认为当前时间
    token VARCHAR(255)                           -- 用户的认证token
);
CREATE TABLE IF NOT EXISTS UserScore (
    id SERIAL PRIMARY KEY,                       -- 积分记录ID，自动递增
	username VARCHAR(30) NOT NULL,               -- 用户名
    record TEXT,                                  -- 积分记录内容
	create_time TIMESTAMP DEFAULT NOW(),         -- 创建时间，默认为当前时间
    FOREIGN KEY (username) REFERENCES UserInformation(username) ON DELETE CASCADE -- 外键约束
);
CREATE TABLE IF NOT EXISTS ImageInformation (
    id SERIAL PRIMARY KEY,                       -- 图像信息的ID，自动递增
	userName VARCHAR(30) NOT NULL,               -- 上传图像的用户名
    params TEXT,                                 -- 图像生成参数
    picture TEXT UNIQUE,                          -- 图像URL，唯一
    likecount INT DEFAULT 0,                     -- 点赞数，默认值为0
    create_time TIMESTAMP DEFAULT NOW(),         -- 创建时间，默认为当前时间
    FOREIGN KEY (userName) REFERENCES UserInformation(username) ON DELETE CASCADE -- 外键约束
);
CREATE TABLE IF NOT EXISTS ImageLike (
    id SERIAL PRIMARY KEY,                       -- 点赞记录的ID，自动递增
    picture TEXT,                                -- 点赞的图像URL
    username TEXT,                               -- 点赞用户的用户名
    num INT DEFAULT 0,                           -- 点赞数，默认值为0
    create_time TIMESTAMP DEFAULT NOW(),         -- 创建时间，默认为当前时间
	FOREIGN KEY (username) REFERENCES UserInformation(username) ON DELETE CASCADE, -- 外键约束
    FOREIGN KEY (picture) REFERENCES ImageInformation(picture) ON DELETE CASCADE -- 外键约束
);

CREATE TABLE IF NOT EXISTS FavoritedImage (
	id SERIAL PRIMARY KEY,                       -- 收藏记录的ID，自动递增
	userName VARCHAR(30) NOT NULL,               -- 收藏图像的用户名
	picture TEXT,                                -- 收藏的图像URL
	create_time TIMESTAMP DEFAULT NOW(),         -- 创建时间，默认为当前时间
	FOREIGN KEY (userName) REFERENCES UserInformation(username) ON DELETE CASCADE, -- 外键约束
    FOREIGN KEY (picture) REFERENCES ImageInformation(picture) ON DELETE CASCADE -- 外键约束
);
`

// DB 是全局数据库连接实例
var DB *gorm.DB

// ConnectDatabase 连接到数据库
func ConnectDatabase() error {
	var err error

	// 获取数据库连接信息
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	// 使用gorm打开数据库连接
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err // 返回连接错误
	}
	return nil
}

// InitDB 初始化数据库，创建所需的表
func InitDB() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized") // 检查数据库连接是否已初始化
	}

	tx := DB.Begin() // 开始事务
	if tx.Error != nil {
		return tx.Error // 返回事务错误
	}

	if err := tx.Exec(createTableSQL).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err    // 返回创建表时的错误
	}

	if err := tx.Commit().Error; err != nil {
		return err // 返回提交事务时的错误
	}

	return nil
}

// InitTestUser 初始化测试用户和相关数据
func InitTestUser() error {
	tx := DB.Begin() // 开始事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 如果发生panic，回滚事务
		}
	}()

	var user user2.UserInformation
	currentTime := time.Now()
	result := tx.Where("username=?", "root").First(&user) // 查找用户名为root的用户

	if result.Error == nil {
		log.Printf("User already exists") // 用户已存在

		// 更新用户分数
		user.Score = 10000
		if result := tx.Save(&user); result.Error != nil {
			log.Printf("Failed to update user score: %v", result.Error)
			tx.Rollback()       // 回滚事务
			return result.Error // 返回更新错误
		}

		tx.Commit() // 提交事务
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to find user: %v", result.Error)
		tx.Rollback()       // 回滚事务
		return result.Error // 返回查找用户错误
	}

	// 创建用户信息
	user = user2.UserInformation{
		Email:       "root@example.com",                                                                                    // 用户邮箱
		UserName:    "root",                                                                                                // 用户名
		Password:    "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a",                                    // 密码（加密后的，111111）
		Avatar_url:  "https://whxh-czh.oss-cn-shenzhen.aliyuncs.com/generate/main.jpg", // 头像URL
		Score:       10000,                                                                                                 // 用户分数
		Create_time: currentTime.AddDate(-1, 0, 0),                                                                         // 创建时间为一年前
	}
	if result := tx.Create(&user); result.Error != nil {
		log.Printf("Failed to create user: %v", result.Error)
		tx.Rollback()       // 回滚事务
		return result.Error // 返回创建用户错误
	}

	// 创建用户积分记录
	userscore := user2.UserScore{
		Username:    "root",                        // 用户名
		Record:      "积分+100",                      // 积分记录内容
		Create_time: currentTime.AddDate(0, -1, 0), // 创建时间为一个月前
	}
	if result := tx.Create(&userscore); result.Error != nil {
		log.Printf("Failed to create record: %v", result.Error)
		tx.Rollback()       // 回滚事务
		return result.Error // 返回创建记录错误
	}

	// 处理图片文件路径
	ImgfilePath := filepath.Join("assets", "examples", "images", "image_urls.txt")
	file, err := os.Open(ImgfilePath) // 打开图像URL文件
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		tx.Rollback() // 回滚事务
		return err    // 返回打开文件错误
	}
	defer file.Close() // 确保文件在函数结束时关闭

	var imageUrls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		imageUrls = append(imageUrls, scanner.Text()) // 读取文件中的每一行作为图像URL
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read file : %v", err)
		tx.Rollback()                // 回滚事务
		return fmt.Errorf("&v", err) // 返回读取文件错误
	}
	i := 0

	// 创建图像记录
	for _, url := range imageUrls {
		i++
		test := fmt.Sprintf("test%d", i)                   // 生成测试参数
		createTime := currentTime.AddDate(0, 0, -(11 + i)) // 创建时间递减
		imageInfo := image2.ImageInformation{
			UserName:    "root",                                                                                                                      // 假设用户名为root
			Params:      "\"Prompt\": \"" + test + "\", \"Width\": \"512\", \"Height\": \"512\", \"Steps\": \"20\", \"SamplingMethod\": \"Euler a\"", // 图像生成参数
			Picture:     url,                                                                                                                         // 图像URL
			Create_time: createTime,                                                                                                                  // 创建时间
		}

		result := tx.Create(&imageInfo) // 创建图像信息记录
		if result.Error != nil {
			log.Printf("Failed to create image information: %v", result.Error)
			return result.Error // 返回创建图像信息错误
		}
	}

	// 如果超过10张图像，随机选取10个用于点赞和收藏
	if len(imageUrls) > 10 {
		rand.Seed(time.Now().UnixNano())                                                                         // 设置随机种子
		rand.Shuffle(len(imageUrls), func(i, j int) { imageUrls[i], imageUrls[j] = imageUrls[j], imageUrls[i] }) // 随机打乱图像URL顺序
		imageUrls = imageUrls[:10]                                                                               // 只保留前10个图像URL
	}

	currentTime = time.Now() // 更新当前时间

	// 创建点赞和收藏记录
	for i, url := range imageUrls {
		createTime := currentTime.AddDate(0, 0, -(i + 1)) // 创建时间递减
		// 创建图像收藏记录
		imagefavor := image2.FavoritedImages{
			UserName:    "root",     // 假设用户名为root
			Picture:     url,        // 图像URL
			Create_time: createTime, // 创建时间
		}
		if result := tx.Create(&imagefavor); result.Error != nil {
			log.Printf("Failed to create image favor for URL %s: %v", url, result.Error)
			tx.Rollback()       // 回滚事务
			return result.Error // 返回创建图像收藏错误
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err) // 返回提交事务时的错误
	}

	return nil // 成功完成初始化
}
