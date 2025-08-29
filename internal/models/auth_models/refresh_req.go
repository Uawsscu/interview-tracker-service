package auth_models

type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
