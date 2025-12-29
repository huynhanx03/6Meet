package service

import (
	"context"
	"net/http"
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/constant"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/mapper"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"go.uber.org/zap"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/apperr"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/response"
	d "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/utils"
)

type userService struct {
	userRepo ports.UserRepository
}

var _ ports.UserService = (*userService)(nil)

func NewUserService(
	userRepo ports.UserRepository,
) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Find lists users with pagination and sorting
func (s *userService) Find(ctx context.Context, opts *d.QueryOptions) (*d.Paginated[*dto.UserResponse], error) {
	// Query database
	users, err := s.userRepo.Find(ctx, opts)
	if err != nil {
		return nil, apperr.Wrap(err, response.CodeDatabaseError, "Failed to list users", http.StatusInternalServerError)
	}

	if users.Records == nil {
		return &d.Paginated[*dto.UserResponse]{
			Records:    &[]*dto.UserResponse{},
			Pagination: users.Pagination,
		}, nil
	}

	// Map generic result from models -> entity
	userEntities := *users.Records
	userResponses := make([]*dto.UserResponse, len(userEntities))
	for i, user := range userEntities {
		userResponses[i] = mapper.ToUserResponse(user)
	}

	return &d.Paginated[*dto.UserResponse]{
		Records:    &userResponses,
		Pagination: users.Pagination,
	}, nil
}

// Get gets a user by ID
func (s *userService) Get(ctx context.Context, id string) (*dto.UserResponse, error) {
	var userResponse dto.UserResponse

	// Check cache
	err := utils.HandleHitCache(ctx, &userResponse, global.Redis, constant.PrefixUser+id)
	if err == nil {
		return &userResponse, nil // Cache Hit
	}

	// Cache Miss: Query Database
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, apperr.Wrap(err, response.CodeInternalServer, "Failed to get user", http.StatusInternalServerError)
	}

	// Set cache
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		e := utils.HandleSetCache(bgCtx, user, global.Redis, constant.PrefixUser+id, constant.CacheExpirationUser)
		if e != nil {
			global.Logger.Error("Failed to set cache", zap.Error(e))
		}
	}()

	// Map model -> entity
	response := *mapper.ToUserResponse(user)

	return &response, nil
}

// Create a new user
func (s *userService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Mapper RequestDTO -> Entity
	user := mapper.ToUserEntityFromReq(req)

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperr.Wrap(err, response.CodeDatabaseError, "Failed to create user", http.StatusInternalServerError)
	}

	// Map entity -> response
	resp := *mapper.ToUserResponse(user)

	return &resp, nil
}

// Update an existing user
func (s *userService) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Check existence
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, apperr.New(response.CodeNotFound, "User not found", http.StatusNotFound, err)
	}

	// Update user
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Neighbors != nil {
		user.Neighbors = *req.Neighbors
	}

	if err := s.userRepo.Update(ctx, id, user); err != nil {
		return nil, apperr.Wrap(err, response.CodeDatabaseError, "Failed to update user", http.StatusInternalServerError)
	}

	// Map entity -> response
	userResponse := *mapper.ToUserResponse(user)

	return &userResponse, nil
}

// Delete an existing user
func (s *userService) Delete(ctx context.Context, id string) error {
	// Check existence
	exists, err := s.userRepo.Exists(ctx, id)
	if err != nil {
		return apperr.Wrap(err, response.CodeInternalServer, "Failed to check user existence", http.StatusInternalServerError)
	}
	if !exists {
		return apperr.New(response.CodeNotFound, "User not found", http.StatusNotFound, nil)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return apperr.Wrap(err, response.CodeDatabaseError, "Failed to delete user", http.StatusInternalServerError)
	}
	return nil
}
