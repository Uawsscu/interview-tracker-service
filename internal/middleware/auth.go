package middleware

import (
	"net/http"
	"strings"
	"time"

	"interview-tracker/internal/config"
	"interview-tracker/internal/pkg/logs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
}

func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		logs.Logger.Println("[Authn] start...")

		// 1) อ่าน header
		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		logs.Logger.Println("[Authn] raw Authorization:", auth)

		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			logs.Logger.Println("[Authn] missing bearer prefix")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := parts[1]
		logs.Logger.Println("[Authn] extracted token length:", len(tokenStr))

		// 2) ต้องมี public key
		pub := config.EnvConfig.JWTKeys.Public
		if pub == nil {
			logs.Logger.Println("[Authn] public key not loaded")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "public key not loaded"})
			return
		}
		logs.Logger.Println("[Authn] public key loaded")

		// 3) parse + verify RS256
		claims := &AccessClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				logs.Logger.Println("[Authn] unexpected signing method:", t.Method.Alg())
				return nil, jwt.ErrSignatureInvalid
			}
			return pub, nil
		})
		if err != nil {
			logs.Logger.Println("[Authn] parse error:", err)
		}
		if err != nil || !token.Valid {
			logs.Logger.Println("[Authn] invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		logs.Logger.Println("[Authn] token signature valid")

		// 4) ตรวจ issuer
		if claims.Issuer != "interview-tracker" {
			logs.Logger.Println("[Authn] invalid issuer:", claims.Issuer)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid issuer"})
			return
		}
		logs.Logger.Println("[Authn] issuer ok:", claims.Issuer)

		// 5) ตรวจ exp
		if claims.ExpiresAt == nil {
			logs.Logger.Println("[Authn] exp is nil")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing exp"})
			return
		}
		now := time.Now()
		if now.After(claims.ExpiresAt.Time) {
			logs.Logger.Println("[Authn] token expired at", claims.ExpiresAt.Time, "now:", now)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}
		logs.Logger.Println("[Authn] token not expired, exp:", claims.ExpiresAt.Time)

		// 6) sub (refID)
		if claims.Subject == "" {
			logs.Logger.Println("[Authn] missing sub")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing sub"})
			return
		}
		logs.Logger.Println("[Authn] refID (sub):", claims.Subject)

		// 7) เช็ค session ใน Redis
		val, err := config.Rdb.Get(c, "session:"+claims.Subject).Result()
		if err != nil || val == "" {
			logs.Logger.Println("[Authn] Redis session not found for refID:", claims.Subject)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired or revoked"})
			return
		}
		logs.Logger.Println("[Authn] Redis session OK for refID:", claims.Subject)

		c.Set("refID", claims.Subject)

		logs.Logger.Println("[Authn] success, passing to next handler")
		c.Next()
	}
}
