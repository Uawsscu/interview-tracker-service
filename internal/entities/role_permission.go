package entities

import "time"

type RolePermission struct {
	RoleID       string    `gorm:"type:uuid;primaryKey" json:"role_id"`
	PermissionID string    `gorm:"type:uuid;primaryKey" json:"permission_id"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedBy    *string   `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy    *string   `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Role       *Role       `gorm:"constraint:OnDelete:CASCADE;foreignKey:RoleID;references:ID" json:"role,omitempty"`
	Permission *Permission `gorm:"constraint:OnDelete:CASCADE;foreignKey:PermissionID;references:ID" json:"permission,omitempty"`
}
