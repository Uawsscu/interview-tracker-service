package entities

import (
	"time"

	"github.com/google/uuid"
)

type CardComment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CardID    string    `json:"card_id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Content   string    `json:"content"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
