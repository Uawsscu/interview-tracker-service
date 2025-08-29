package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"interview-tracker/internal/config"
	"interview-tracker/internal/pkg/logs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	jwt.RegisteredClaims
}

type Session struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Perms  []string  `json:"perms"`
}

// core authn: validate JWT + redis session, return refID + session struct
func authenticate(c *gin.Context) (*Session, string, bool) {
	logs.Logger.Println("[Authn] start...")

	// 1) header
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return nil, "", false
	}
	tokenStr := parts[1]

	// 2) public key
	pub := config.EnvConfig.JWTKeys.Public
	if pub == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "public key not loaded"})
		return nil, "", false
	}

	// 3) parse + verify
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return pub, nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return nil, "", false
	}

	// 4) issuer
	if claims.Issuer != "interview-tracker" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid issuer"})
		return nil, "", false
	}

	// 5) exp
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return nil, "", false
	}

	// 6) subject
	refID := claims.Subject
	if refID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing sub"})
		return nil, "", false
	}

	// 7) session in redis
	val, err := config.Rdb.Get(c, "session:"+refID).Result()
	if err != nil || val == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired or revoked"})
		return nil, "", false
	}

	var s Session
	_ = json.Unmarshal([]byte(val), &s)

	// inject to context
	c.Set("refID", refID)
	logs.Logger.Println(">>>>>>")
	logs.LogJson(s)
	c.Set("session", s)

	return &s, refID, true
}

// ใช้แค่ auth (ไม่เช็ค permission)
func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, _, ok := authenticate(c); !ok {
			return
		}
		c.Next()
	}
}

// ใช้ auth + เช็ค permission
func Authorize(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, _, ok := authenticate(c)
		if !ok {
			return
		}
		// check permission
		has := false
		for _, p := range s.Perms {
			if p == required {
				has = true
				break
			}
		}
		if !has {
			logs.Logger.Printf("[Authorize] deny: need %s", required)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func GetSession(c *gin.Context) *Session {
	if v, ok := c.Get("session"); ok {
		if s, ok := v.(Session); ok {
			return &s
		}
	}
	return nil
}
