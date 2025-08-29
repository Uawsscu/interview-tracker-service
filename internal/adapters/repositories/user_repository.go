package repositories

import (
	"interview-tracker/internal/entities"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user entities.User) (*entities.User, error)
	GetById(id string) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	GetPermissionsByUserID(userID string) ([]string, error)
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

func (r *userRepo) GetByEmail(email string) (*entities.User, error) {
	var u entities.User
	if err := r.db.
		Preload("Role").
		First(&u, "email = ? AND is_active = true", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetPermissionsByUserID(userID string) ([]string, error) {
	type row struct{ Code string }
	var rows []row
	q := `
			SELECT p.code
			FROM users u
			JOIN roles r ON r.id = u.role_id AND r.is_active = true
			JOIN role_permissions rp ON rp.role_id = r.id AND rp.is_active = true
			JOIN permissions p ON p.id = rp.permission_id AND p.is_active = true
			WHERE u.id = ? AND u.is_active = true
		`
	if err := r.db.Raw(q, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	perms := make([]string, 0, len(rows))
	for _, x := range rows {
		perms = append(perms, x.Code)
	}
	return perms, nil
}
