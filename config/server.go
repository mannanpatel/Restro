package config

import "os"

type ServerConfig struct {
	Port      string
	Mode      string // "debug", "release", "test"
	Env       string // "development", "production", "staging"
	JwtSecret string

	// // HTTP Timeouts — production ke liye important
	// ReadTimeoutSecs  int // request read karne ka max time
	// WriteTimeoutSecs int // response bhejne ka max time
	// IdleTimeoutSecs  int // keep-alive connection ka max idle time

	// // CORS — frontend se requests allow karne ke liye
	// AllowedOrigins []string // ["http://localhost:3000", "https://yourdomain.com"]
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:      os.Getenv("SERVER_PORT"),
		Mode:      os.Getenv("SERVER_MODE"),
		Env:       os.Getenv("SERVER_ENV"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
