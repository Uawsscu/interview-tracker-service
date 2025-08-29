package auth_models

type LoginReq struct {
	Email    string `json:"email" binding:"required,email" example:"example@example.com"`
	Password string `json:"password" binding:"required" example:"P@ssw0rd"`
}
