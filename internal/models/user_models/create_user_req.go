package user_models

type CreateUserRequest struct {
	Name     string `json:"name"  validate:"required,max=120" example:"Soda Pop"`
	Email    string `json:"email" validate:"required" example:"example@example.com"`
	Password string `json:"password" validate:"required" example:"P@ssw0rd"`
	RoleID   string `json:"role_id" validate:"required" example:"a1bf3d66-e4ae-4d73-89c6-917f0f301003"`
	IsActive bool   `json:"is_active"`
}
