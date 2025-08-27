package entities

import "time"

type Role struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null;default:''" json:"description"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Users []User `gorm:"foreignKey:RoleID" json:"-"`
}
