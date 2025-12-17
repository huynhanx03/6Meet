package mapper

import (
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/db/models"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToUserEntity converts DB Model to Domain Entity
func ToUserEntity(m *models.User) *entity.User {
	return &entity.User{
		ID:        m.BaseModel.ID.Hex(),
		Name:      m.Name,
		Neighbors: m.Neighbors,
		CreatedAt: m.BaseModel.CreatedAt,
		UpdatedAt: m.BaseModel.UpdatedAt,
	}
}

// ToUserModel converts Domain Entity to DB Model
func ToUserModel(e *entity.User) *models.User {
	id, _ := primitive.ObjectIDFromHex(e.ID)
	
	return &models.User{
		BaseModel: &mongodb.BaseModel{
            ID: id,
            CreatedAt: e.CreatedAt,
            UpdatedAt: e.UpdatedAt,
        },
		Name:      e.Name,
		Neighbors: e.Neighbors,
	}
}

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
