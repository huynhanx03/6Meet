package ports

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/models"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"
)

type IUserRepository interface {
	mongodb.IBaseRepository[models.User]
}