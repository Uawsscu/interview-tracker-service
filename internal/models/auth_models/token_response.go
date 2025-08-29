package auth_models

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6Ikp.."`
	RefreshToken string `json:"refresh_token" example:"0b83305f-75b2-4e89-a444-70da68f84d4f.0224f190-0f9a-46a4-a0fc-521fcc07e2ee.."`
	RefID        string `json:"ref_id" example:"b2a8f4e5-9e3b-4c2b-9e0f-7b2c1a6d1f33"`
}
