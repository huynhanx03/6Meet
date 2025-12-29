package mapper

import (
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
)

// ToUserResponse converts Domain Entity to Response DTO
func ToUserResponse(e *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        e.ID,
		Name:      e.Name,
		Neighbors: e.Neighbors,
	}
}

// ToUserEntityFromReq converts Request DTO to Domain Entity
func ToUserEntityFromReq(req *dto.CreateUserRequest) *entity.User {
	return &entity.User{
		Name:      req.Name,
		Neighbors: req.Neighbors,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
