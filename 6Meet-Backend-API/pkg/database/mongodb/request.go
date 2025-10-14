package mongodb

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPaginationOptions creates MongoDB options for pagination
func GetPaginationOptions(p *dto.PaginationOptions) *options.FindOptions {
	skip := int64((p.Page - 1) * p.PageSize)
	limit := int64(p.PageSize)

	return &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
}
