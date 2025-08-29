package entities

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CandidateName string     `json:"candidate_name"`
	ScheduledAt   time.Time  `json:"scheduled_at"`
	StatusCode    string     `json:"status_code" gorm:"not null"`
	CreatedBy     uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"`
	AssigneeID    *uuid.UUID `json:"assignee_id,omitempty" gorm:"type:uuid"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
