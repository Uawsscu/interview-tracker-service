package common_body

type DirectionEnum string

const (
	DESC DirectionEnum = "desc"
	ASC  DirectionEnum = "asc"
)

type OrderBy struct {
	Column string        `json:"column" validate:"required" example:"created_at"`
	Sort   DirectionEnum `json:"sort" validate:"required" example:"ASC"`
}
