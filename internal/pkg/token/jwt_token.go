package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"interview-tracker/internal/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func MustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func SignAccess(refID string, ttl time.Duration) (string, error) {
	if config.EnvConfig == nil || config.EnvConfig.JWTKeys.Private == nil {
		return "", errors.New("jwt private key not loaded")
	}
	privateKey := config.EnvConfig.JWTKeys.Private

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": refID,
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
		"iss": "interview-tracker",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}
