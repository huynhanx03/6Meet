package search

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/search/models"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/elasticsearch"
	d "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
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

func (s *userSearch) Search(ctx context.Context, opts *d.QueryOptions) (*d.Paginated[*entity.User], error) {
	if opts == nil {
		opts = &d.QueryOptions{}
	}
	if opts.Pagination == nil {
		opts.Pagination = &d.PaginationOptions{}
	}
	opts.Pagination.SetDefaults()

	// Build generic query
	queryMap := elasticsearch.BuildSearchQuery(opts)

	// Marshal to JSON
	queryJSON, err := json.Marshal(queryMap)
	if err != nil {
		return nil, err
	}

	// Execute search
	docs, err := s.repo.Search(ctx, strings.NewReader(string(queryJSON)))
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(docs))
	for i, doc := range docs {
		users[i] = doc.ToEntity()
	}

	// Calculate pagination
	totalItems := int64(len(docs)) // Approx

	// Pagination options are already defaulted in BuildSearchQuery -> BuildPagination
	return &d.Paginated[*entity.User]{
		Records:    &users,
		Pagination: d.CalculatePagination(opts.Pagination.Page, opts.Pagination.PageSize, totalItems),
	}, nil
}
