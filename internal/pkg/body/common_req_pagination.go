package common_body

type CommonRequestPagination struct {
	Limit  int `json:"limit"  example:"20"`
	Offset int `json:"offset" example:"0"`
}
