package entities

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // hashed
	RoleID    string    `gorm:"type:uuid;not null" json:"role_id"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedBy string    `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy string    `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Role *Role `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;foreignKey:RoleID;references:ID" json:"role,omitempty"`
}
