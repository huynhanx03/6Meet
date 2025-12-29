package models

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/elasticsearch"
)

// UserDoc represents the user document in Elasticsearch
type UserDoc struct {
	*elasticsearch.BaseDocument
	Name      string   `json:"name"`
	Neighbors []string `json:"neighbors"`
}

// ToEntity converts Elasticsearch document to domain entity
func (d *UserDoc) ToEntity() *entity.User {
	return &entity.User{
		ID:        d.ID,
		Name:      d.Name,
		Neighbors: d.Neighbors,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// FromEntity converts domain entity to Elasticsearch document
func FromEntity(user *entity.User) *UserDoc {
	return &UserDoc{
		BaseDocument: &elasticsearch.BaseDocument{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Name:      user.Name,
		Neighbors: user.Neighbors,
	}
}
