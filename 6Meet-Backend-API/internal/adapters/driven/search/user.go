package search

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/search/models"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/elasticsearch"
)

const (
	userIndex = "users"
)

type userSearch struct {
	repo *elasticsearch.BaseRepository[models.UserDoc]
}

// NewUser creates a new user search repository
func NewUser(client elasticsearch.ElasticClient) ports.UserSearchRepository {
	return &userSearch{
		repo: elasticsearch.NewBaseRepository[models.UserDoc](client, userIndex),
	}
}

func (s *userSearch) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userSearch) BatchDelete(ctx context.Context, ids []string) error {
	return s.repo.BatchDelete(ctx, ids)
}

func (s *userSearch) Index(ctx context.Context, user *entity.User) error {
	doc := models.FromEntity(user)
	return s.repo.Index(ctx, doc)
}

func (s *userSearch) BatchIndex(ctx context.Context, users []*entity.User) error {
	docs := make([]*models.UserDoc, len(users))
	for i, user := range users {
		docs[i] = models.FromEntity(user)
	}
	return s.repo.BatchIndex(ctx, docs)
}
