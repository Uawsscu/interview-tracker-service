package common_body

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnResponse(c *gin.Context, res interface{}, err interface{}) {
	if e, isError := err.(error); isError {
		errorResponse(e.Error(), c)
		return
	} else if err != nil {
		errorResponse(err, c)
		return
	} else if res != nil {
		successResponse(res, c)
	} else {
		successResponseStatus(c)
	}
}

type CommonStatusResponse struct {
	Status string `json:"status"`
}
type CommonSuccessResponse struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}

type CommonErrorResponse struct {
	Status       string      `json:"status"`
	ErrorMessage interface{} `json:"errorMessage"`
}

func errorResponse(errRes any, c *gin.Context) {
	errorResponse := CommonErrorResponse{
		Status:       "error",
		ErrorMessage: errRes,
	}

	// logs.Logger.Printf(string_util.ToJsonString(errorResponse))
	c.JSON(http.StatusBadRequest, errorResponse)
}

func successResponse(resp any, c *gin.Context) {
	successResponse := CommonSuccessResponse{
		Status:   "success",
		Response: &resp,
	}
	// logs.Logger.Printf(string_util.ToJsonString(successResponse))
	c.JSON(http.StatusOK, successResponse)
}

func successResponseStatus(c *gin.Context) {
	successResponse := CommonStatusResponse{
		Status: "success",
	}
	// logs.Logger.Printf(string_util.ToJsonString(successResponse))
	c.JSON(http.StatusOK, successResponse)
}
