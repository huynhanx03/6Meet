package service

import (
	"context"
	"errors"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/mapper"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"

	d "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService struct {
	userRepo ports.IUserRepository
}

func NewUserService(
	userRepo ports.IUserRepository,
) ports.IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

// ListUsers lists users with pagination and sorting
func (s *userService) ListUsers(ctx context.Context, opts *d.ListOptions) (*d.ListResult[dto.UserResponse], error) {
	users, err := s.userRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(*users.Records))
	for i, user := range *users.Records {
		userResponses[i] = *mapper.ModelToUserResponse(&user)
	}

	return &d.ListResult[dto.UserResponse]{
		Records:    &userResponses,
		Pagination: users.Pagination,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid student ID format")
	}

	user, err := s.userRepo.GetByID(ctx, objectID)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return mapper.ModelToUserResponse(user), nil
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := mapper.CreateUserRequestToModel(req)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return mapper.ModelToUserResponse(user), nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user = mapper.UpdateUserRequestToModel(req, user)
	if err := s.userRepo.Update(ctx, objectID, user); err != nil {
		return nil, err
	}

	return mapper.ModelToUserResponse(user), nil
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	exists, err := s.userRepo.ExistsByID(ctx, objectID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(ctx, objectID)
}
