package application

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort uint16
	GRPCPort   uint16
	MongoURI   string

	TurnstileSecret  string
	TurnstileDisable bool

	SMTPHost     string
	SMTPPort     uint16
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	NotifyEmail  string
}

func LoadConfig() Config {
	// Dev-only: giúp chạy `go run` đọc .env; trong Docker không cần
	_ = godotenv.Load()
	cfg := Config{
		ServerPort: 8082,
		GRPCPort:   50052,
		SMTPPort:   587,
	}
	// Load server port from env
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 16); err == nil {
			cfg.ServerPort = uint16(p)
		}
	}
	// Load gRPC port from env
	if v := os.Getenv("GRPC_PORT"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 16); err == nil {
			cfg.GRPCPort = uint16(p)
		}
	}
	// Load MongoDB URI from env
	if v := os.Getenv("MONGODB_URI"); v != "" {
		cfg.MongoURI = v
	}
	// Load Turnstile secret from env
	cfg.TurnstileSecret = os.Getenv("TURNSTILE_SECRET")
	if cfg.MongoURI == "" || cfg.TurnstileSecret == "" {
		log.Fatal("Missing required env: MONGODB_URI or TURNSTILE_SECRET")
	}
	// Load Turnstile disable from env
	if v := os.Getenv("TURNSTILE_DISABLE"); v != "" {
		b, _ := strconv.ParseBool(v)
		cfg.TurnstileDisable = b
	}
	// Load SMTP config from env
	cfg.SMTPHost = os.Getenv("SMTP_HOST")
	// Load SMTP port from env
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 16); err == nil {
			cfg.SMTPPort = uint16(p)
		}
	}
	// Load SMTP user from env
	if v := os.Getenv("SMTP_USER"); v != "" {
		cfg.SMTPUser = v
	}
	// Load SMTP password from env
	cfg.SMTPUser = os.Getenv("SMTP_USER")
	cfg.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	cfg.FromEmail = os.Getenv("FROM_EMAIL")
	cfg.NotifyEmail = os.Getenv("NOTIFY_EMAIL")

	return cfg
}
