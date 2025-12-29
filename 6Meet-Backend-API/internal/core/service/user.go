package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/mapper"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"go.uber.org/zap"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/cdc"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/apperr"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/response"
	d "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
)

type userService struct {
	userRepo   ports.UserRepository
	cacheRepo  ports.UserCacheRepository
	searchRepo ports.UserSearchRepository
}

func NewUser(
	userRepo ports.UserRepository,
	cacheRepo ports.UserCacheRepository,
	searchRepo ports.UserSearchRepository,
) ports.UserService {
	return &userService{
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
		searchRepo: searchRepo,
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
	// Check cache
	userResponse, err := s.cacheRepo.Get(ctx, id)
	if err == nil {
		return userResponse, nil // Cache Hit
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
		e := s.cacheRepo.Set(bgCtx, user)
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

// HandleUserChange processes CDC events for users
func (s *userService) HandleUserChange(ctx context.Context, evt *cdc.DebeziumPayload[entity.User]) error {
	switch evt.Op {
	case cdc.OpCreate:
		if evt.After == nil {
			return fmt.Errorf("operation %s requires 'after' state", evt.Op)
		}
		user := evt.After

		// Sync to Elasticsearch
		if err := s.searchRepo.Index(ctx, user); err != nil {
			global.Logger.Error("Failed to sync user to search", zap.String("id", user.ID), zap.Error(err))
		} else {
			global.Logger.Info("Synced user to search", zap.String("id", user.ID))
		}
	case cdc.OpUpdate:
		if evt.After == nil {
			return fmt.Errorf("operation %s requires 'after' state", evt.Op)
		}
		user := evt.After

		// Sync to Cache
		if _, err := s.cacheRepo.Get(ctx, user.ID); err == nil {
			if err := s.cacheRepo.Set(ctx, user); err != nil {
				global.Logger.Error("Failed to sync user to cache", zap.String("id", user.ID), zap.Error(err))
			} else {
				global.Logger.Info("Synced user to cache", zap.String("id", user.ID))
			}
		}

		// Sync to Elasticsearch
		if err := s.searchRepo.Index(ctx, user); err != nil {
			global.Logger.Error("Failed to sync user to search", zap.String("id", user.ID), zap.Error(err))
		} else {
			global.Logger.Info("Synced user to search", zap.String("id", user.ID))
		}

	case cdc.OpDelete:
		if evt.Before == nil {
			return fmt.Errorf("operation %s requires 'before' state", evt.Op)
		}
		userID := evt.Before.ID

		// Delete from Cache
		if err := s.cacheRepo.Delete(ctx, userID); err != nil {
			global.Logger.Error("Failed to delete user from cache", zap.String("id", userID), zap.Error(err))
		}

		// Delete from Elasticsearch
		if err := s.searchRepo.Delete(ctx, userID); err != nil {
			global.Logger.Error("Failed to delete user from search", zap.String("id", userID), zap.Error(err))
		}
	}

	return nil
}

// HandleUserBatchChange processes a batch of CDC events for users
func (s *userService) HandleUserBatchChange(ctx context.Context, evts []*cdc.DebeziumPayload[entity.User]) error {
	var usersToCache []*entity.User
	var usersToIndex []*entity.User
	var idsToDelete []string

	for _, evt := range evts {
		switch evt.Op {
		case cdc.OpCreate, cdc.OpUpdate:
			if evt.After == nil {
				continue
			}
			user := evt.After
			if _, err := s.cacheRepo.Get(ctx, user.ID); err == nil {
				usersToCache = append(usersToCache, user)
			}

			usersToIndex = append(usersToIndex, user)
		case cdc.OpDelete:
			if evt.Before != nil {
				idsToDelete = append(idsToDelete, evt.Before.ID)
			}
		}
	}

	// Bulk Sync to Cache
	if len(usersToCache) > 0 {
		if err := s.cacheRepo.BatchSet(ctx, usersToCache); err != nil {
			global.Logger.Error("Failed to batch sync users to cache", zap.Error(err))
		} else {
			// global.Logger.Info("Batch synced users to cache", zap.Int("count", len(usersToCache)))
		}
	}

	if len(idsToDelete) > 0 {
		if err := s.cacheRepo.BatchDelete(ctx, idsToDelete); err != nil {
			global.Logger.Error("Failed to batch delete users from cache", zap.Error(err))
		} else {
			// global.Logger.Info("Batch deleted users from cache", zap.Int("count", len(idsToDelete)))
		}
	}

	// Bulk Sync to Elasticsearch
	if len(usersToIndex) > 0 {
		if err := s.searchRepo.BatchIndex(ctx, usersToIndex); err != nil {
			global.Logger.Error("Failed to batch sync users to search", zap.Error(err))
		} else {
			// global.Logger.Info("Batch synced users to search", zap.Int("count", len(usersToIndex)))
		}
	}

	if len(idsToDelete) > 0 {
		if err := s.searchRepo.BatchDelete(ctx, idsToDelete); err != nil {
			global.Logger.Error("Failed to batch delete users from search", zap.Error(err))
		} else {
			// global.Logger.Info("Batch deleted users from search", zap.Int("count", len(idsToDelete)))
		}
	}

	return nil
}
