package ports

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	d "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
)

// UserRepository defines the interface for user repository
type UserRepository interface {
	Find(ctx context.Context, opts *d.QueryOptions) (*d.Paginated[*entity.User], error)
	Get(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, id string, user *entity.User) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// UserService defines the interface for user service
type UserService interface {
	Find(ctx context.Context, opts *d.QueryOptions) (*d.Paginated[*dto.UserResponse], error)
	Get(ctx context.Context, id string) (*dto.UserResponse, error)

	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
}
