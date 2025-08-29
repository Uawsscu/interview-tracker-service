package repositories

import (
	"interview-tracker/internal/entities"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetList() ([]*entities.Role, error)
	ExistsByID(id string) (bool, error)
}

type roleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) RoleRepository { return &roleRepo{db} }

func (r *roleRepo) GetList() ([]*entities.Role, error) {
	var roles []*entities.Role
	if err := r.db.
		Model(&entities.Role{}).
		Order("created_at ASC").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleRepo) ExistsByID(id string) (bool, error) {
	var cnt int64
	if err := r.db.Model(&entities.Role{}).Where("id = ? AND is_active = TRUE", id).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}
