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

// UserHandler defines the interface for user handler
type UserHandler interface {
	Find(c *gin.Context)
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// userHandler implements UserHandler
type userHandler struct {
	handler.BaseHandler
	userService ports.UserService
}

var _ UserHandler = (*userHandler)(nil)

func NewUserHandler(userService ports.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// Find handles the HTTP request to list users with pagination and sorting
func (h *userHandler) Find(c *gin.Context) {
	opts, ok := request.ParseRequest[d.QueryOptions](c)

	if !ok {
		return
	}

	users, err := h.userService.Find(c.Request.Context(), opts)
	if err != nil {
		response.ErrorResponse(c, response.CodeInternalServer, err)
		return
	}

	response.SuccessResponse(c, response.CodeRetrieved, users)
}

// Get handles the HTTP request to get a user by ID
func (h *userHandler) Get(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.Get(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, response.CodeInternalServer, err)
		return
	}

	response.SuccessResponse(c, response.CodeRetrieved, user)
}

// Create handles the HTTP request to create a new user
func (h *userHandler) Create(c *gin.Context) {
	req, ok := request.ParseRequest[dto.CreateUserRequest](c)

	if !ok {
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req)
	if err != nil {
		response.ErrorResponse(c, response.CodeInternalServer, err)
		return
	}

	response.SuccessResponse(c, response.CodeCreated, user)
}

// Update handles the HTTP request to update an existing user
func (h *userHandler) Update(c *gin.Context) {
	req, ok := request.ParseRequest[dto.UpdateUserRequest](c)

	if !ok {
		return
	}

	id := c.Param("id")

	user, err := h.userService.Update(c.Request.Context(), id, req)
	if err != nil {
		response.ErrorResponse(c, response.CodeInternalServer, err)
		return
	}

	response.SuccessResponse(c, response.CodeUpdated, user)
}

// Delete handles the HTTP request to delete a user by ID
func (h *userHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, response.CodeInternalServer, err)
		return
	}

	response.SuccessResponse(c, response.CodeDeleted, nil)
}
