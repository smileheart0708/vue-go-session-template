package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	Port          int    // 服务监听端口
	DataDir       string // 数据持久化目录
	LogLevel      string // 日志等级
	AuthKey       string // 管理员身份验证密钥，同时用于 Session 签名
	CookieSecure  bool   // Session Cookie 是否启用 Secure
	IsAutoAuthKey bool   // AuthKey 是否自动生成
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	// 若当前目录存在 .env，先加载到进程环境变量（不会覆盖已存在变量）。
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	cfg := &Config{
		Port:         getEnvAsInt("PORT", 8080),
		DataDir:      getEnv("DATA_DIR", ".data"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		AuthKey:      getEnv("AUTH_KEY", ""),
		CookieSecure: getEnvAsBool("COOKIE_SECURE", false),
	}

	// 如果 AUTH_KEY 未设置，生成随机 12 位字符串
	if cfg.AuthKey == "" {
		cfg.AuthKey = generateRandomKey(12)
		cfg.IsAutoAuthKey = true
	}

	return cfg, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取整数类型的环境变量，如果不存在或解析失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsBool 获取布尔类型环境变量，支持 true/false（不区分大小写）
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// generateRandomKey 生成指定长度的随机十六进制字符串
func generateRandomKey(length int) string {
	// 每个字节生成 2 个十六进制字符，所以需要 length/2 个字节
	bytes := make([]byte, (length+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用固定的默认值
		return "default_auth_key"
	}

	key := hex.EncodeToString(bytes)
	// 截取到指定长度
	if len(key) > length {
		key = key[:length]
	}

	return key
}
