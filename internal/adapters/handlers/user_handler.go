package handlers

import (
	"net/http"

	"interview-tracker/internal/models/user_models"
	"interview-tracker/internal/pkg/logs"
	"interview-tracker/internal/pkg/utils"
	"interview-tracker/internal/usecases"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase *usecases.UserUsecase
}

func NewUserHandler(usecase *usecases.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

// @Summary      Create user
// @Description  Create user for
// @Accept       json
// @Produce      json
// @Tags         users
// @Param        request  body  user_models.CreateUserRequest  true  "user create JSON"
// @Success      201      {object} user_models.CreateUserResponse  "Success response"
// @Failure      400      {object} map[string]string "Bad Request"
// @Router       /interview-tracker/internal/v1/users/create [post]
func (h *UserHandler) Create(c *gin.Context) {
	logs.Logger.Printf("[user] create start...")

	var request user_models.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.Logger.Printf("[user] create bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if mapErr := utils.ValidateRequest(request); mapErr != nil {
		logs.Logger.Printf("[user] create validate error: %+v", mapErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": mapErr})
		return
	}

	res, err := h.usecase.CreateUser(request)
	if err != nil {
		logs.Logger.Printf("[user] create failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logs.Logger.Printf("[user] create success")
	c.JSON(http.StatusCreated, res)
}

// @Summary Get user By ID
// @Description  Get user By ID
// @Accept json
// @Produce json
// @Tags users
// @Param userId path string true "User ID" default(888f2c6b-cc1a-4e94-bd6e-d8ba0ac36fc3)
// @Success 200 {object} entities.User "Successful response"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /interview-tracker/internal/v1/users/details/{userId} [get]
func (h *UserHandler) GetById(c *gin.Context) {
	logs.Logger.Printf("[user] GetById start....")

	userId := c.Param("userId")
	resp, err := h.usecase.GetUserById(userId)
	if err != nil {
		logs.Logger.Printf("[user] GetById failed: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	logs.Logger.Printf("[user] GetById success....")
	c.JSON(http.StatusOK, resp)
}

// @Summary Get role
// @Description  Get role
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {array} entities.Role "Successful response"
// @Failure 500 {object} map[string]string "Server error"
// @Router /interview-tracker/internal/v1/users/role-list [get]
func (h *UserHandler) GetListRole(c *gin.Context) {
	logs.Logger.Printf("[user] ListRole start....")

	resp, err := h.usecase.GetRoleList()
	if err != nil {
		logs.Logger.Printf("[user] ListRole failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logs.Logger.Printf("[user] ListRole success....")
	c.JSON(http.StatusOK, resp)
}
