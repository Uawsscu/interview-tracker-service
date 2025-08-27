package config

import "os"

type Config struct {
	UserAdminDefault string
	DatabaseUrl      string
	HttpPort         string
}

func LoadConfig() *Config {
	return &Config{
		UserAdminDefault: os.Getenv("USER_ADMIN_DEFAULT"),
		DatabaseUrl:      os.Getenv("DATABASE_URL"),
		HttpPort:         os.Getenv("HTTP_PORT"),
	}
}
