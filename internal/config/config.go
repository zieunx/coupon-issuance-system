package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config는 애플리케이션의 모든 설정을 담는 구조체
type Config struct {
	MySQL       MySQLConfig
	Redis       RedisConfig
	AdminServer AdminServerConfig
	IssueServer AdminServerConfig
}

// MySQLConfig는 MySQL 데이터베이스 연결 정보를 담는 구조체
type MySQLConfig struct {
	DSN string // Data Source Name
}

// RedisConfig는 Redis 연결 정보를 담는 구조체
type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

// ServerConfig는 서버 관련 설정을 담는 구조체
type AdminServerConfig struct {
	Host string
	Port string
}

var appConfig *Config

// LoadConfig는 .env 파일과 환경변수에서 설정을 로드합니다.
func LoadConfig() {
	godotenv.Load()

	appConfig = &Config{
		MySQL: MySQLConfig{
			DSN: getEnv("MYSQL_DSN", "coupon:coupon123@tcp(127.0.0.1:3306)/coupondb?parseTime=true"),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		AdminServer: AdminServerConfig{
			Host: getEnv("ADMIN_SERVER_HOST", "localhost"),
			Port: getEnv("ADMIN_SERVER_PORT", "8081"),
		},
		IssueServer: AdminServerConfig{
			Port: getEnv("ISSUE_SERVER_PORT", "8082"),
		},
	}
}

// GetConfig는 로드된 설정을 반환합니다.
func GetConfig() *Config {
	if appConfig == nil {
		log.Fatal("Config not loaded. Call LoadConfig first.")
	}
	return appConfig
}

// getEnv는 환경변수를 읽고, 값이 없으면 기본값을 반환합니다.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt는 환경변수를 정수로 읽고, 실패 시 기본값을 반환합니다.
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
