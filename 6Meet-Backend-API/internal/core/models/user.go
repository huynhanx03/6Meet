package models

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/model"
)

type User struct {
	*model.BaseModel `bson:",inline"`
	Name            string   `json:"name" bson:"name"`
	Neighbors       []string `json:"neighbors" bson:"neighbors"`
}
