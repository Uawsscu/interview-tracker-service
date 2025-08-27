package routers

import (
	"interview-tracker/internal/adapters/handlers"
	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/config"
	"interview-tracker/internal/usecases"

	"github.com/gin-gonic/gin"
)

func User(r *gin.RouterGroup) {
	db := config.DB
	userRepo := repositories.NewUserRepo(db)
	roleRepo := repositories.NewRoleRepo(db)
	userUseCases := usecases.NewUserUsecase(userRepo, roleRepo)
	userHandler := handlers.NewUserHandler(userUseCases)
	r.POST("/internal/v1/users/create", userHandler.Create)
	r.GET("/internal/v1/users/details/:userId", userHandler.GetById)
	r.GET("/internal/v1/users/role-list", userHandler.GetListRole)
}
