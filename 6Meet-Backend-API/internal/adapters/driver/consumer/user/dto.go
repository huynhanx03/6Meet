package user

import (
	"encoding/json"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/cdc"
)

// CDCUser mimics the entity.User but with CDC-compatible field types
type CDCUser struct {
	ID        json.RawMessage `json:"_id"` // Can be string or {"$oid":...}
	Name      string          `json:"name"`
	Neighbors []string        `json:"neighbors"`
	CreatedAt *cdc.MongoDate  `json:"created_at"`
	UpdatedAt *cdc.MongoDate  `json:"updated_at"`
}

// ToEntity converts the CDC DTO to the domain entity
func (c *CDCUser) ToEntity() *entity.User {
	u := &entity.User{
		Name:      c.Name,
		Neighbors: c.Neighbors,
	}

	// Parse ID using the generic helper
	if len(c.ID) > 0 {
		u.ID = cdc.ParseMongoDBKey(c.ID)
	}

	// Parse Dates using generic helper
	if c.CreatedAt != nil {
		u.CreatedAt = c.CreatedAt.ToTime()
	}
	if c.UpdatedAt != nil {
		u.UpdatedAt = c.UpdatedAt.ToTime()
	}

	return u
}
