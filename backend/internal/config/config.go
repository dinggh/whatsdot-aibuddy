package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
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
	ForceDevWeChat bool
	LogDir         string
	UploadDir      string

	OpenAIBaseURL string
	OpenAIAPIKey  string
	OpenAIModel   string
	AnalyzeMock   bool

	RateLimitCapacity int
	RateLimitRefill   int
}

func Load() Config {
	loadEnvFile(".env")
	expireHours := getEnvInt("JWT_EXPIRE_HOURS", 168)
	cfg := Config{
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/aibuddy?sslmode=disable"),
		WeChatAppID:    os.Getenv("WECHAT_APP_ID"),
		WeChatSecret:   os.Getenv("WECHAT_APP_SECRET"),
		WeChatAPIBase:  getEnv("WECHAT_API_BASE", "https://api.weixin.qq.com"),
		JWTSecret:      getEnv("JWT_SECRET", "change-this-jwt-secret"),
		JWTExpireAfter: time.Duration(expireHours) * time.Hour,
		ForceDevWeChat: getEnvBool("FORCE_DEV_WECHAT", false),
		LogDir:         getEnv("LOG_DIR", "logs"),
		UploadDir:      getEnv("UPLOAD_DIR", "uploads"),

		OpenAIBaseURL: getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:   getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		AnalyzeMock:   getEnvBool("ANALYZE_MOCK", false),

		RateLimitCapacity: getEnvInt("RATE_LIMIT_CAPACITY", 6),
		RateLimitRefill:   getEnvInt("RATE_LIMIT_REFILL_PER_MIN", 6),
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return strings.EqualFold(v, "true") || v == "1"
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

// loadEnvFile reads a .env file and sets variables into os.Environ.
// It skips empty lines and lines starting with #. File not found is ignored.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		if key == "" {
			continue
		}
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		os.Setenv(key, val)
	}
}
