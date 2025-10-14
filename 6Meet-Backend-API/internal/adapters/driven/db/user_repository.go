package driven

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/models"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	*mongodb.BaseRepository[models.User]
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(mongo *mongo.Database) ports.IUserRepository {
	collection := mongo.Collection("users")

	return &userRepository{
		BaseRepository: mongodb.NewBaseRepository[models.User](collection),
	}
}
