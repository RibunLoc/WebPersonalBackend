package application

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddress  string // địa chỉ redis server
	RedisUsername string // tên user login redis
	RedisPassword string // mật khẩu login
	MongoURI      string
	ServerPort    uint16 // cổng lắng nghe của backend
	JwtSecret     string // Secret JWT
}

func LoadConfig() Config {
	_ = godotenv.Load()
	cfg := Config{
		ServerPort: 3000, // default server port
	}

	if redisAddr, exist := os.LookupEnv("REDIS_ADDR"); exist {
		cfg.RedisAddress = redisAddr
	}

	if redisUser, exits := os.LookupEnv("REDIS_USERNAME"); exits {
		cfg.RedisUsername = redisUser
	}

	if redisPass, exist := os.LookupEnv("REDIS_PASSWORD"); exist {
		cfg.RedisPassword = redisPass
	}

	if mongoURI, exist := os.LookupEnv("MONGODB_URI"); exist {
		cfg.MongoURI = mongoURI
	}

	if serverPort, exist := os.LookupEnv("SERVER_PORT"); exist {
		// Chuyển kiểu dạng số về chuỗi (string)
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	if jwtSecret, exist := os.LookupEnv("JWT_SECRET_KEY"); exist {
		cfg.JwtSecret = jwtSecret
	}

	// Kiểm tra các trường bắt buộc
	if cfg.JwtSecret == "" || cfg.MongoURI == "" {
		log.Fatal("Missing required env: JWT_SECRET_KEY or MONGODB_URI")
	}

	return cfg
}
