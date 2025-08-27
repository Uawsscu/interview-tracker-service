package common_body

type CommonResponsePagination struct {
	TotalCount int `json:"totalCount"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`

	// CurrentPage  int `json:"currentPage"`
	// PageSize     int `json:"pageSize"`
	// TotalRecords int `json:"totalRecords"`
	// TotalPages   int `json:"totalPages"`
}
