package card_models

import (
	"time"
)

type ListCardsQuery struct {
	Page     int    `form:"page,default=1" example:"1"`
	PageSize int    `form:"page_size,default=10" example:"10"`
	Status   string `form:"status" example:"todo"` // todo|in_progress|done
}

type CreateCardReq struct {
	Title       string    `json:"title" binding:"required" example:"นัดสัมภาษณ์งาน 1"`
	Description string    `json:"description" example:"สัมภาษณ์ตำแหน่ง Backend Developer"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required" example:"2023-01-01T10:00:00Z"`
	// AssigneeID  *uuid.UUID `json:"assignee_id" example:"888f2c6b-cc1a-4e94-bd6e-d8ba0ac36fc3"`
}

type UpdateCardReq struct {
	Title       *string    `json:"title" example:"นัดสัมภาษณ์งาน 2"`
	Description *string    `json:"description" example:"สัมภาษณ์ตำแหน่ง Fullstack Developer"`
	ScheduledAt *time.Time `json:"scheduled_at" example:"2023-01-02T15:00:00Z"`
}

type UpdateCardStatusReq struct {
	Status string `json:"status" binding:"required,oneof=todo in_progress done" example:"in_progress"`
}

type AddCommentReq struct {
	Content string `json:"content" binding:"required" example:"ควรปรับปรุง portfolio ให้ละเอียดขึ้น"`
}

type UpdateCommentReq struct {
	Content string `json:"content" binding:"required" example:"ใช้ได้"`
}
