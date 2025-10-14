package request

import (
	"github.com/gin-gonic/gin"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/response"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/validation"
)

type RequestParamSetter interface {
	SetParams(params gin.Params)
}

func ParseRequest[T any](c *gin.Context) (*T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamInvalid, response.ToErrorResponse(err))
		return nil, false
	}

	if ok, msg := validation.IsRequestValid(req); !ok {
		response.ErrorResponse(c, response.ErrCodeValidationFailed, response.ToErrorResponse(msg))
		return nil, false
	}

	return &req, true
}
