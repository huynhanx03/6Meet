package models

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	*mongodb.BaseModel `bson:",inline"`
	Name               string   `json:"name" bson:"name"`
	Neighbors          []string `json:"neighbors" bson:"neighbors"`
}

// ToEntity converts DB Model to Domain Entity
func (m *User) ToEntity() *entity.User {
	return &entity.User{
		ID:        m.BaseModel.ID.Hex(),
		Name:      m.Name,
		Neighbors: m.Neighbors,
		CreatedAt: m.BaseModel.CreatedAt,
		UpdatedAt: m.BaseModel.UpdatedAt,
	}
}

// FromEntity converts Domain Entity to DB Model
func FromEntity(e *entity.User) *User {
	id, _ := primitive.ObjectIDFromHex(e.ID)

	return &User{
		BaseModel: &mongodb.BaseModel{
			ID:        id,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		Name:      e.Name,
		Neighbors: e.Neighbors,
	}
}
