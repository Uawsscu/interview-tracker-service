package entities

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CandidateName string    `json:"candidate_name"`
	ScheduledAt   time.Time `json:"scheduled_at"`
	StatusCode    string    `json:"status_code" gorm:"not null"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedBy     uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy     uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
