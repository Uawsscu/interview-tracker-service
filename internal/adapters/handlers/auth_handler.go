package handlers

import (
	"net/http"

	"interview-tracker/internal/models/auth_models"
	"interview-tracker/internal/pkg/logs"
	"interview-tracker/internal/usecases"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ auth usecases.AuthUsecase }

func NewAuthHandler(a usecases.AuthUsecase) *AuthHandler { return &AuthHandler{auth: a} }

// @Summary      Login
// @Description  Authenticate user and issue access/refresh tokens (JWT access token has only sub/exp/iat/iss)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  auth_models.LoginReq  true  "login JSON"
// @Success      200      {object} auth_models.TokenResponse  "Success"
// @Router       /interview-tracker/internal/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	logs.Logger.Printf("[auth] login start...")
	var req auth_models.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logs.Logger.Printf("[auth] login bind error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid request"})
		return
	}
	access, refresh, refID, err := h.auth.Login(c, req.Email, req.Password)
	if err != nil {
		logs.Logger.Printf("[auth] login error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	logs.Logger.Printf("[auth] login success...")
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh, "ref_id": refID})
}

// @Summary      Refresh token
// @Description  Rotate refresh token and issue a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  auth_models.RefreshReq  true  "refresh JSON"
// @Success      200      {object} auth_models.TokenResponse  "Success"
// @Router       /interview-tracker/internal/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	logs.Logger.Printf("[auth] refresh start...")
	var req auth_models.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAccess, newRefresh, refID, err := h.auth.Refresh(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	logs.Logger.Printf("[auth] refresh success...")
	c.JSON(http.StatusOK, gin.H{"access_token": newAccess, "refresh_token": newRefresh, "ref_id": refID})
}

// @Summary      Logout
// @Description  Revoke current session using access token (refID is taken from JWT `sub`)
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]string  "ok"
// @Failure      401  {object}  map[string]string  "unauthorized"
// @Failure      500  {object}  map[string]string  "logout failed"
// @Router       /interview-tracker/internal/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refID := c.GetString("refID") // set โดย auth middleware หลัง verify JWT
	if refID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.auth.Logout(c, refID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
