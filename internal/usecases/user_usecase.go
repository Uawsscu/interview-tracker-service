package usecases

import (
	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/config"
	"interview-tracker/internal/entities"
	"interview-tracker/internal/models/user_models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewUserUsecase(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (u *UserUsecase) CreateUser(item user_models.CreateUserRequest) (*entities.User, error) {
	// TODO check duplicate email
	// TODO check exists roleId

	// ---- hash password ----
	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ---- build payload ----
	var currentDtm = time.Now()
	var userAdminDefault = config.EnvConfig.UserAdminDefault
	payload := entities.User{
		Name:      item.Name,
		Email:     item.Email,
		Password:  string(hashed),
		RoleID:    item.RoleID,
		IsActive:  true,
		CreatedBy: userAdminDefault,
		UpdatedBy: userAdminDefault,
		CreatedAt: currentDtm,
		UpdatedAt: currentDtm,
	}

	return u.userRepo.Create(payload)
}

func (u *UserUsecase) GetUserById(id string) (*entities.User, error) {
	return u.userRepo.GetById(id)
}

func (u *UserUsecase) GetRoleList() ([]*entities.Role, error) {
	return u.roleRepo.GetList()
}
