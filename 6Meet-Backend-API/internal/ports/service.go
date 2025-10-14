package ports

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/dto"
	d "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
)

// IUserService defines the interface for user service
type IUserService interface {
	ListUsers(ctx context.Context, opts *d.ListOptions) (*d.ListResult[dto.UserResponse], error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)

	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}
