package http

import (
	"github.com/gin-gonic/gin"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/handler"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/request"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/response"

	d "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"

)

type userHandler struct {
	handler.BaseHandler
	userService ports.IUserService
}

func NewUserHandler(userService ports.IUserService) ports.IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

// ListUsers handles the HTTP request to list users with pagination and sorting
func (h *userHandler) ListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		opts, ok := request.ParseRequest[d.ListOptions](c)

		if !ok {
			return
		}

		users, err := h.userService.ListUsers(c.Request.Context(), opts)
		if err != nil {
			response.ErrorResponse(c, response.ErrCodeInternalServer, response.ToErrorResponse(err))
			return
		}

		response.SuccessResponse(c, response.ErrCodeRetrieved, users)
	}
}

// GetUserByID handles the HTTP request to get a user by ID
func (h *userHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, err := h.userService.GetUserByID(c.Request.Context(), id)
		if err != nil {
			response.ErrorResponse(c, response.ErrCodeInternalServer, response.ToErrorResponse(err))
			return
		}

		response.SuccessResponse(c, response.ErrCodeRetrieved, user)
	}
}

// CreateUser handles the HTTP request to create a new user
func (h *userHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, ok := request.ParseRequest[dto.CreateUserRequest](c)

		if !ok {
			return
		}

		user, err := h.userService.CreateUser(c.Request.Context(), req)
		if err != nil {
			response.ErrorResponse(c, response.ErrCodeInternalServer, response.ToErrorResponse(err))
		}

		response.SuccessResponse(c, response.ErrCodeCreated, user)
	}
}

// UpdateUser handles the HTTP request to update an existing user
func (h *userHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, ok := request.ParseRequest[dto.UpdateUserRequest](c)

		if !ok {
			return
		}

		id := c.Param("id")

		user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
		if err != nil {
			response.ErrorResponse(c, response.ErrCodeInternalServer, response.ToErrorResponse(err))
		}

		response.SuccessResponse(c, response.ErrCodeUpdated, user)
	}
}

// DeleteUser handles the HTTP request to delete a user by ID
func (h *userHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := h.userService.DeleteUser(c.Request.Context(), id)
		if err != nil {
			response.ErrorResponse(c, response.ErrCodeInternalServer, response.ToErrorResponse(err))
		}

		response.SuccessResponse(c, response.ErrCodeDeleted, nil)
	}
}
