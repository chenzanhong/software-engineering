// config/config.go
package config

import (
	"io/ioutil"
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

// LoadConfig 加载配置：优先环境变量，降级使用 YAML
func LoadConfig() (*Config, error) {
	var config Config

	// 1️⃣ 先尝试从环境变量加载
	config = Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			Name:     getEnv("DB_NAME", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
		},
		OSS: OSSConfig{
			OSS_REGION:            getEnv("OSS_REGION", ""),
			OSS_ACCESS_KEY_ID:     getEnv("OSS_ACCESS_KEY_ID", ""),
			OSS_ACCESS_KEY_SECRET: getEnv("OSS_ACCESS_KEY_SECRET", ""),
			OSS_BUCKET:            getEnv("OSS_BUCKET", ""),
		},
		Model: ModelConfig{
			GEN_API_KEY: getEnv("GEN_API_KEY", ""),
			Time:        getEnv("MODEL_TIMEOUT", ""),
		},
	}

	// 2️⃣ 检查是否有字段为空，如果有，则从 YAML 文件加载
	// 如果环境变量中缺少关键字段（如 DB_HOST），则加载 YAML
	if config.DB.Host == "" || config.DB.Name == "" {
		log.Println("环境变量中缺少关键配置，尝试加载 config.yaml...")
		yamlPath := GetDBConfigPath()

		yamlFile, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			log.Printf("⚠️ 无法读取 config.yaml: %v", err)
			// 即使文件不存在，也不 panic，继续使用环境变量（可能部分为空）
			return &config, nil // 返回部分配置
		}

		var yamlConfig Config
		err = yaml.Unmarshal(yamlFile, &yamlConfig)
		if err != nil {
			log.Printf("❌ 解析 config.yaml 失败: %v", err)
			return &config, nil // 降级返回环境变量配置
		}

		// 合并：仅填充环境变量中为空的字段
		if config.DB.Host == "" {
			config.DB.Host = yamlConfig.DB.Host
		}
		if config.DB.Port == "" {
			config.DB.Port = yamlConfig.DB.Port
		}
		if config.DB.Name == "" {
			config.DB.Name = yamlConfig.DB.Name
		}
		if config.DB.User == "" {
			config.DB.User = yamlConfig.DB.User
		}
		if config.DB.Password == "" {
			config.DB.Password = yamlConfig.DB.Password
		}

		// OSS 和 Model 同理
		if config.OSS.OSS_REGION == "" {
			config.OSS.OSS_REGION = yamlConfig.OSS.OSS_REGION
		}
		if config.OSS.OSS_ACCESS_KEY_ID == "" {
			config.OSS.OSS_ACCESS_KEY_ID = yamlConfig.OSS.OSS_ACCESS_KEY_ID
		}
		if config.OSS.OSS_ACCESS_KEY_SECRET == "" {
			config.OSS.OSS_ACCESS_KEY_SECRET = yamlConfig.OSS.OSS_ACCESS_KEY_SECRET
		}
		if config.OSS.OSS_BUCKET == "" {
			config.OSS.OSS_BUCKET = yamlConfig.OSS.OSS_BUCKET
		}

		if config.Model.GEN_API_KEY == "" {
			config.Model.GEN_API_KEY = yamlConfig.Model.GEN_API_KEY
		}
		if config.Model.Time == "" {
			config.Model.Time = yamlConfig.Model.Time
		}

		log.Printf("✅ 已从 config.yaml 补充配置: DB=%s@%s:%s", config.DB.User, config.DB.Host, config.DB.Port)
	}

	// 3️⃣ 最终检查关键字段
	if config.DB.Host == "" || config.DB.Name == "" || config.DB.User == "" {
		return nil, &ConfigError{"缺少数据库配置（DB_HOST、DB_NAME、DB_USER）"}
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
