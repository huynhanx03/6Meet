package mongodb

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	db *mongo.Client
}

func NewMongoRepository(db *mongo.Client) *mongoRepository {
	return &mongoRepository{
		db: db,
	}
}

// IBaseRepository defines the common interface for all repositories
type IBaseRepository[T BaseModel] interface {
	List(ctx context.Context, opts *dto.ListOptions) (*dto.ListResult[T], error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*T, error)

	Create(ctx context.Context, model *T) error
	Update(ctx context.Context, id primitive.ObjectID, model *T) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteMany(ctx context.Context, filter bson.M) (int64, error)
	ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
}
