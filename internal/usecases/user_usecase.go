package usecases

import (
	"errors"
	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/entities"
	"interview-tracker/internal/models/user_models"
	"strings"
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
	// normalize
	email := strings.TrimSpace(strings.ToLower(item.Email))

	// ---- check duplicate email ----
	if existing, _ := u.userRepo.GetByEmail(email); existing != nil {
		return nil, errors.New("email already in use")
	}

	// ---- check exists roleId ----
	ok, err := u.roleRepo.ExistsByID(item.RoleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("role not found")
	}
	// ---- hash password ----
	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ---- build payload ----
	var currentDtm = time.Now()
	payload := entities.User{
		Name:      item.Name,
		Email:     item.Email,
		Password:  string(hashed),
		RoleID:    item.RoleID,
		IsActive:  true,
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
