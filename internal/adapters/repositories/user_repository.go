package repositories

import (
	"interview-tracker/internal/entities"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user entities.User) (*entities.User, error)
	GetById(id string) (*entities.User, error)
}

type userRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) UserRepository { return &userRepo{db} }

func (r *userRepo) Create(user entities.User) (*entities.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetById(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
