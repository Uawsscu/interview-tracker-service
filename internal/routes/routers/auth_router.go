package routers

import (
	"interview-tracker/internal/adapters/handlers"
	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/config"
	"interview-tracker/internal/usecases"

	"github.com/gin-gonic/gin"
)

func Auth(r *gin.RouterGroup) {
	db := config.DB
	rdb := config.Rdb
	userRepo := repositories.NewUserRepo(db)
	authUseCases := usecases.NewAuthUsecase(userRepo, rdb)
	authHandler := handlers.NewAuthHandler(authUseCases)
	r.POST("/internal/v1/auth/login", authHandler.Login)
	r.POST("/internal/v1/auth/refresh", authHandler.Refresh)
	r.POST("/internal/v1/auth/logout", authHandler.Logout)
}
