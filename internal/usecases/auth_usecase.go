package usecases

import (
	"context"
	"errors"
	"time"

	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/config"
	"interview-tracker/internal/pkg/token"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// ===== unchanged interfaces/types =====
type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken, refID string, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccess, newRefresh, refID string, err error)
	Logout(ctx context.Context, refID string) error
}

type authUsecase struct {
	userRepo repositories.UserRepository
	rdb      *redis.Client
}

func NewAuthUsecase(userRepo repositories.UserRepository, rdb *redis.Client) AuthUsecase {
	return &authUsecase{userRepo: userRepo, rdb: rdb}
}

type sessionData struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Role   string   `json:"role"`
	Perms  []string `json:"perms"`
}

// ===== main flows =====

func (s *authUsecase) Login(ctx context.Context, email, password string) (string, string, string, error) {
	u, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	perms, err := s.userRepo.GetPermissionsByUserID(u.ID)
	if err != nil {
		return "", "", "", err
	}

	refID := uuid.NewString()
	accessTTL := readMinutes(config.EnvConfig.AccessTTLMinutes)
	refreshTTLDays := config.EnvConfig.RefreshTTLDays
	sessionData := sessionData{UserID: u.ID, Email: u.Email, Role: u.Role.Code, Perms: perms}

	if err := s.rdb.Set(ctx, "session:"+refID, token.MustJSON(sessionData), accessTTL).Err(); err != nil {
		return "", "", "", err
	}

	access, err := token.SignAccess(refID, accessTTL)
	if err != nil {
		_ = s.rdb.Del(ctx, "session:"+refID).Err()
		return "", "", "", err
	}

	rt, err := token.NewRefreshToken()
	if err != nil {
		_ = s.rdb.Del(ctx, "session:"+refID).Err()
		return "", "", "", err
	}
	if err := s.saveRefresh(ctx, refID, rt, refreshTTLDays); err != nil {
		_ = s.rdb.Del(ctx, "session:"+refID).Err()
		return "", "", "", err
	}

	return access, rt, refID, nil
}

func (s *authUsecase) Refresh(ctx context.Context, refreshToken string) (string, string, string, error) {
	hash := token.Sha256Hex(refreshToken)
	refID, err := s.rdb.Get(ctx, "rt:"+hash).Result()
	if err != nil {
		return "", "", "", errors.New("invalid or expired refresh token")
	}

	_ = s.rdb.Del(ctx, "rt:"+hash).Err()

	accessTTL := readMinutes(config.EnvConfig.AccessTTLMinutes)
	val, err := s.rdb.Get(ctx, "session:"+refID).Result()
	if err != nil {
		return "", "", "", errors.New("session not found; please login again")
	}
	_ = s.rdb.Expire(ctx, "session:"+refID, accessTTL).Err()

	newAccess, err := token.SignAccess(refID, accessTTL)
	if err != nil {
		return "", "", "", err
	}

	refreshTTLDays := config.EnvConfig.RefreshTTLDays
	newRT, err := token.NewRefreshToken()
	if err != nil {
		return "", "", "", err
	}
	if err := s.saveRefresh(ctx, refID, newRT, refreshTTLDays); err != nil {
		return "", "", "", err
	}

	_ = val
	return newAccess, newRT, refID, nil
}

func (s *authUsecase) Logout(ctx context.Context, refID string) error {
	_ = s.rdb.Del(ctx, "session:"+refID).Err()
	iter := s.rdb.Scan(ctx, 0, "rt:*", 1000).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		rid, _ := s.rdb.Get(ctx, key).Result()
		if rid == refID {
			_ = s.rdb.Del(ctx, key).Err()
		}
	}
	return iter.Err()
}

func (s *authUsecase) saveRefresh(ctx context.Context, refID, raw string, ttlDays int) error {
	hash := token.Sha256Hex(raw)
	return s.rdb.Set(ctx, "rt:"+hash, refID, days(ttlDays)).Err()
}

// ===== helpers =====
func readMinutes(numMin int) time.Duration { return time.Duration(numMin) * time.Minute }
func days(n int) time.Duration             { return time.Hour * 24 * time.Duration(n) }
