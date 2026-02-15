package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr     string
	DatabaseURL    string
	WeChatAppID    string
	WeChatSecret   string
	WeChatAPIBase  string
	JWTSecret      string
	JWTExpireAfter time.Duration
}

func Load() Config {
	expireHours := getEnvInt("JWT_EXPIRE_HOURS", 168)
	cfg := Config{
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/aibuddy?sslmode=disable"),
		WeChatAppID:    os.Getenv("WECHAT_APP_ID"),
		WeChatSecret:   os.Getenv("WECHAT_APP_SECRET"),
		WeChatAPIBase:  getEnv("WECHAT_API_BASE", "https://api.weixin.qq.com"),
		JWTSecret:      getEnv("JWT_SECRET", "change-this-jwt-secret"),
		JWTExpireAfter: time.Duration(expireHours) * time.Hour,
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
