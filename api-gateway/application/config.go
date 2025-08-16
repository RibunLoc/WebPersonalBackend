package application

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      uint16 // configures the server port
	AuthGRPCAddr    string
	AuthHTTPBase    string
	ContactGRPCAddr string // ví dụ: http://localhost:50052
	JWTSecret       string // dùng cho middleware verify token
}

func LoadConfig() Config {
	_ = godotenv.Load()

	cfg := Config{
		ServerPort:   8080,
		AuthGRPCAddr: "localhost:50051",
		AuthHTTPBase: "http://localhost:3000",

		ContactGRPCAddr: "localhost:50052",
	}

	if v := os.Getenv("GATEWAY_PORT"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 16); err == nil {
			cfg.ServerPort = uint16(p)
		}
	}
	if v := os.Getenv("AUTH_GRPC_ADDR"); v != "" {
		cfg.AuthGRPCAddr = v
	}
	if v := os.Getenv("AUTH_HTTP_BASE"); v != "" {
		cfg.AuthHTTPBase = v
	}
	if v := os.Getenv("JWT_SECRET_KEY"); v != "" { // nên đồng bộ với auth-service
		cfg.JWTSecret = v
	}
	if v := os.Getenv("CONTACT_GRPC_ADDR"); v != "" {
		cfg.ContactGRPCAddr = v
	}
	return cfg
}
