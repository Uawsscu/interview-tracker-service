package config

import (
	"crypto/rsa"
	"errors"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

var EnvConfig *Config

type Config struct {
	DatabaseUrl      string
	HttpPort         string
	RedisAddr        string
	JWTSecret        string
	AccessTTLMinutes int
	RefreshTTLDays   int
	JWTKeys          JWTKeys
}
type JWTKeys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func LoadConfig() {
	loadJWTKeys, err := LoadJWTKeys()
	if err != nil {
		panic(err)
	}
	EnvConfig = &Config{
		DatabaseUrl:      os.Getenv("DATABASE_URL"),
		HttpPort:         os.Getenv("HTTP_PORT"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		AccessTTLMinutes: envAsInt("ACCESS_TTL_MINUTES", 15),
		RefreshTTLDays:   envAsInt("REFRESH_TTL_DAYS", 7),
		JWTKeys:          *loadJWTKeys,
	}
}

func envAsInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func LoadJWTKeys() (*JWTKeys, error) {
	privPath := os.Getenv("JWT_PRIVATE_KEY_PEM")
	pubPath := os.Getenv("JWT_PUBLIC_KEY_PEM")

	var (
		priv *rsa.PrivateKey
		pub  *rsa.PublicKey
	)

	if privPath != "" {
		b, err := os.ReadFile(privPath)
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
		if err != nil {
			return nil, err
		}
		priv = key
		pub = &key.PublicKey
	}

	if pubPath != "" {
		b, err := os.ReadFile(pubPath)
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(b)
		if err != nil {
			return nil, err
		}
		pub = key
	}

	if priv == nil && pub == nil {
		return nil, errors.New("no JWT keys found")
	}

	return &JWTKeys{Private: priv, Public: pub}, nil
}
