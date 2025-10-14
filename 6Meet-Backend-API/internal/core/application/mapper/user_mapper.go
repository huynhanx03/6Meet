package mapper

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/models"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/model"
)

// CreateUserRequestToModel maps CreateUserRequest DTO to User models
func CreateUserRequestToModel(req *dto.CreateUserRequest) *models.User {
	return &models.User{
		BaseModel: model.NewBaseModel(),
		Name:      req.Name,
		Neighbors: req.Neighbors,
	}
}

// UpdateUserRequestToModel maps UpdateUserRequest DTO to User models
func UpdateUserRequestToModel(req *dto.UpdateUserRequest, user *models.User) *models.User {
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Neighbors != nil {
		user.Neighbors = *req.Neighbors
	}
	return user
}

// ModelToUserResponse maps User models to UserResponse DTO
func ModelToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Neighbors: user.Neighbors,
	}
}
