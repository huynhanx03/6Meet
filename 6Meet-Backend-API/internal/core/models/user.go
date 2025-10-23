package models

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"
)

type User struct {
	*mongodb.BaseModel `bson:",inline"`
	Name               string   `json:"name" bson:"name"`
	Neighbors          []string `json:"neighbors" bson:"neighbors"`
}
