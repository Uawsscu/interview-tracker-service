package handlers

import (
	"interview-tracker/internal/models/user_models"
	common_body "interview-tracker/internal/pkg/body"
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
// @Success      200      {object} user_models.CreateUserResponse  "Success response"
// @Router       /interview-tracker/internal/v1/users/create [post]
func (h *UserHandler) Create(c *gin.Context) {
	logs.Logger.Printf("[user] create star...")

	var request user_models.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logs.Logger.Printf("[user] create bind error: %v", err)
		common_body.ReturnResponse(c, nil, err)
		return
	}

	if mapErr := utils.ValidateRequest(request); mapErr != nil { // หรือ &request ตาม validator ของคุณ
		logs.Logger.Printf("[user] create validate error: %+v", mapErr)
		common_body.ReturnResponse(c, nil, mapErr)
		return
	}

	res, err := h.usecase.CreateUser(request)
	if err != nil {
		logs.Logger.Printf("[user] create failed: %v", err)
	} else {
		logs.Logger.Printf("[user] create success")
	}
	common_body.ReturnResponse(c, res, err)
}

// @Summary Get user By ID
// @Description  Get user By ID
// @Accept json
// @Tags users
// @Param userId path string true "User ID" default(888f2c6b-cc1a-4e94-bd6e-d8ba0ac36fc3)
// @Success 200 {object} entities.User "Successful response"
// @Router /interview-tracker/internal/v1/users/details/{userId} [get]
func (h *UserHandler) GetById(c *gin.Context) {
	logs.Logger.Printf("[user] GetById start....")

	userId := c.Param("userId")
	resp, err := h.usecase.GetUserById(userId)
	logs.Logger.Printf("[user] GetById success....")

	common_body.ReturnResponse(c, resp, err)
}

// @Summary Get role
// @Description  Get role
// @Accept json
// @Tags users
// @Success 200 {object} entities.Role "Successful response"
// @Router /interview-tracker/internal/v1/users/role-list [get]
func (h *UserHandler) GetListRole(c *gin.Context) {
	logs.Logger.Printf("[user] ListRole start....")
	resp, err := h.usecase.GetRoleList()
	logs.Logger.Printf("[user] ListRole success....")

	common_body.ReturnResponse(c, resp, err)
}
