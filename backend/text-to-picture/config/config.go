// config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type OSSConfig struct {
	OSS_REGION            string `yaml:"OSS_REGION"`
	OSS_ACCESS_KEY_ID     string `yaml:"OSS_ACCESS_KEY_ID"`
	OSS_ACCESS_KEY_SECRET string `yaml:"OSS_ACCESS_KEY_SECRET"`
	OSS_BUCKET            string `yaml:"OSS_BUCKET"`
}

type ModelConfig struct {
	GEN_API_KEY string `yaml:"GEN_API_KEY"`
	Time        string `yaml:"timeout"`
}

type Config struct {
	DB    DBConfig    `yaml:"db"`
	OSS   OSSConfig   `yaml:"oss"`
	Model ModelConfig `yaml:"model"`
}

// 工作目录相对路径（+ flag/env）	⭐⭐⭐⭐☆（非常普遍）	✅ 推荐
// runtime.Caller() 构建路径	⭐（极少）	❌ 不推荐
// 绝对路径硬编码	⭐	❌ 不推荐

func getDBConfigPath() string { 
	// 获取调用者的文件名（即 login_test.go 或 findByFeature.go）
	_, filename, _, ok := runtime.Caller(2) // 注意这里使用 Caller(2)
	if !ok {
		log.Fatal("无法获取运行时调用者信息")
	}

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(filename)

	// 构建到项目根目录的相对路径
	dbConfigPath := filepath.Join(currentDir, "..", "config", "configs", "config.yaml")

	// 将路径转换为绝对路径并简化路径（移除冗余的 '..'）
	absPath, err := filepath.Abs(dbConfigPath)
	if err != nil {
		log.Fatalf("无法获取绝对路径: %v", err)
	}

	simplifiedPath := filepath.Clean(absPath)

	return simplifiedPath
}

func GetDBConfigPath() string {
	return getDBConfigPath()
}

// getEnv 读取环境变量，如果为空则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getConfigPath 获取 config.yaml 的路径（支持相对路径）
func getConfigPath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("无法获取调用者信息")
	}
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "configs", "config.yaml")
}

// LoadConfig 加载配置：YAML 为默认值，环境变量优先覆盖
func LoadConfig() (*Config, error) {
	// 1️⃣ 先加载 config.yaml 作为基础配置（默认值）
	var config Config
	yamlPath := GetDBConfigPath()

	log.Printf("正在读取配置文件: %s", yamlPath)
	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, &ConfigError{fmt.Sprintf("config.yaml 读取失败: %v", err)}
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, &ConfigError{fmt.Sprintf("config.yaml 解析失败: %v", err)}
	}

	log.Println("✅ 已从 config.yaml 加载默认配置")

	// 2️⃣ 用环境变量覆盖 YAML 中的值（环境变量优先）
	if v := os.Getenv("DB_HOST"); v != "" {
		config.DB.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		config.DB.Port = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		config.DB.Name = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		config.DB.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		config.DB.Password = v
	}

	if v := os.Getenv("OSS_REGION"); v != "" {
		config.OSS.OSS_REGION = v
	}
	if v := os.Getenv("OSS_ACCESS_KEY_ID"); v != "" {
		config.OSS.OSS_ACCESS_KEY_ID = v
	}
	if v := os.Getenv("OSS_ACCESS_KEY_SECRET"); v != "" {
		config.OSS.OSS_ACCESS_KEY_SECRET = v
	}
	if v := os.Getenv("OSS_BUCKET"); v != "" {
		config.OSS.OSS_BUCKET = v
	}

	if v := os.Getenv("GEN_API_KEY"); v != "" {
		config.Model.GEN_API_KEY = v
	}
	if v := os.Getenv("MODEL_TIMEOUT"); v != "" {
		config.Model.Time = v
	}

	// 3️⃣ 最终校验关键字段
	if config.DB.Host == "" {
		return nil, &ConfigError{"DB_HOST 不能为空"}
	}
	if config.DB.Name == "" {
		return nil, &ConfigError{"DB_NAME 不能为空"}
	}
	if config.DB.User == "" {
		return nil, &ConfigError{"DB_USER 不能为空"}
	}
	if config.Model.GEN_API_KEY == "" {
		return nil, &ConfigError{"GEN_API_KEY 不能为空（请在环境变量或 config.yaml 中设置）"}
	}

	return &config, nil
}

// ConfigError 自定义配置错误
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "配置错误: " + e.Message
}
