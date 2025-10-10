package user_r

import (
	"errors"
	"fmt"
	"regexp"
	"text-to-picture/models/image"
	userLogin "text-to-picture/models/user"
	"time"

	"gorm.io/gorm"
)

// 正则表达式验证邮箱格式
func isValidEmail(email string) bool {

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email) // 返回是否匹配
}

// 向用户信息表插入数据
func InsertUserInformation(db *gorm.DB, user *userLogin.UserInformation) error {
	if user.UserName == "" {
		return fmt.Errorf("名字为空") // 检查用户名是否为空
	}
	if user.Email == "" {
		return fmt.Errorf("邮箱为空") // 检查邮箱是否为空
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("密码少于6位")// 检查密码长度
	}
	if isValidEmail(user.Email) == false {
		return fmt.Errorf("邮箱格式不正确") // 验证邮箱格式
	}
	var existingUserLogin userLogin.UserInformation

	// 检查用户名是否已存在
	result := db.Where("UserName = ?", user.UserName).First(&existingUserLogin)
	if result.Error == nil {
		return fmt.Errorf("用户名已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return fmt.Errorf("查询用户名时发生错误: %v", result.Error)
	}

	result = db.Where("Email = ?", user.Email).First(&existingUserLogin)
	if result.Error == nil {
		return fmt.Errorf("邮箱已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询邮箱时发生错误: %v", result.Error)
	}
	user.Create_time = time.Now() // 设置创建时间
	user.Score = 100 // 注册时给一点积分

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("插入用户表失败: %v", err)
	}

	return nil // 插入成功
}

// 向图片信息表插入数据
func InsertImageInformation(db *gorm.DB, user *image.ImageInformation) error {
	if user.Params == "" {
		return fmt.Errorf("参数为空") // 检查参数是否为空
	}
	if user.Picture == "" {
		return fmt.Errorf("结果为空") // 检查图像结果是否为空
	}
	if user.Create_time.IsZero() {
		return fmt.Errorf("时间参数为空") // 检查创建时间是否为空
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("插入图片信息表失败: %v", err) // 插入失败
	}

	return nil

}

// 向用户收藏表插入数据
func InsertFavoritedImage(db *gorm.DB, user *image.ImageInformation) error {
	if user.Picture == "" { // 检查图像结果是否为空
		return fmt.Errorf("结果为空") 
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("插入用户收藏表失败: %v", err)
	}
	return nil
}
