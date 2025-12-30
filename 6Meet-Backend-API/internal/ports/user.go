package ports

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/cdc"
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

// UserCacheRepository defines the interface for user cache operations
type UserCacheRepository interface {
	Set(ctx context.Context, user *entity.User) error
	BatchSet(ctx context.Context, users []*entity.User) error
	BatchDelete(ctx context.Context, ids []string) error
	Get(ctx context.Context, id string) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
}

// UserService defines the interface for user service
type UserService interface {
	Find(ctx context.Context, opts *d.QueryOptions) (*d.Paginated[*dto.UserResponse], error)
	Get(ctx context.Context, id string) (*dto.UserResponse, error)

	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error

	HandleUserChange(ctx context.Context, evt *cdc.DebeziumPayload[entity.User]) error
	HandleUserBatchChange(ctx context.Context, evts []*cdc.DebeziumPayload[entity.User]) error

	Search(ctx context.Context, opts *d.QueryOptions) (*d.Paginated[*dto.UserResponse], error)
}

// UserSearchRepository defines the interface for user search engine operations
type UserSearchRepository interface {
	Index(ctx context.Context, user *entity.User) error
	BatchIndex(ctx context.Context, users []*entity.User) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Search(ctx context.Context, opts *d.QueryOptions) (*d.Paginated[*entity.User], error)
}

// UserProducer defines the interface for user related events
type UserProducer interface {
	Publish(ctx context.Context, event any) error
}

// UserConsumer defines the interface for user event consumer
type UserConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}
