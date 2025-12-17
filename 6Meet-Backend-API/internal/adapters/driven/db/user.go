package db

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/db/models"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/mapper"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "users"
)

type userRepository struct {
	repo *mongodb.BaseRepository[models.User]
}

var _ ports.UserRepository = (*userRepository)(nil)

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *mongo.Database) ports.UserRepository {
	collection := db.Collection(userCollection)
	return &userRepository{
		repo: mongodb.NewBaseRepository[models.User](collection),
	}
}

// Find users by query options
func (r *userRepository) Find(ctx context.Context, opts *dto.QueryOptions) (*dto.Paginated[*entity.User], error) {
	// Get models
	res, err := r.repo.Find(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Map generic result from models -> entity
    if res.Records == nil {
         return &dto.Paginated[*entity.User]{
            Records:    &[]*entity.User{},
            Pagination: res.Pagination,
        }, nil
    }

	models := *res.Records
	entities := make([]*entity.User, len(models))
	for i := range models {
		entities[i] = mapper.ToUserEntity(&models[i])
	}

	return &dto.Paginated[*entity.User]{
		Records:    &entities,
		Pagination: res.Pagination,
	}, nil
}

// Get user by ID
func (r *userRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	// Convert string ID to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	model, err := r.repo.Get(ctx, oid)
	if err != nil {
		return nil, err
	}

	// Map model -> entity
	return mapper.ToUserEntity(model), nil
}

// Create a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	// Map entity -> model
	model := mapper.ToUserModel(user)

	// Create model in database
	err := r.repo.Create(ctx, model)
	if err != nil {
		return err
	}
	
	user.ID = model.ID.Hex()
	
	return nil
}

// Update user by ID
func (r *userRepository) Update(ctx context.Context, id string, user *entity.User) error {
	// Convert string ID to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	// Map entity -> model
	model := mapper.ToUserModel(user)

	return r.repo.Update(ctx, oid, model)
}

// Delete user by ID
func (r *userRepository) Delete(ctx context.Context, id string) error {
	// Convert string ID to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	return r.repo.Delete(ctx, oid)
}

// Check if user exists by ID
func (r *userRepository) Exists(ctx context.Context, id string) (bool, error) {
	// Convert string ID to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	return r.repo.Exists(ctx, oid)
}
