package user_models

type CreateUserResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleID   string `json:"role_id"`
	IsActive bool   `json:"is_active"`
}
