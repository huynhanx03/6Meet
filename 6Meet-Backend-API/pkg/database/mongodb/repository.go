package mongodb

import (
	"context"
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModelEntity interface that all models must implement
type ModelEntity interface {
	GetID() primitive.ObjectID
	UpdateTimestamp()
}

// BaseRepository provides common database operations using generics
type BaseRepository[T ModelEntity] struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T ModelEntity](collection *mongo.Collection) *BaseRepository[T] {
	return &BaseRepository[T]{
		collection: collection,
		timeout:    30 * time.Second,
	}
}

// GetContext creates a context with timeout
func (r *BaseRepository[T]) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.timeout)
}

// GetCollection returns the MongoDB collection
func (r *BaseRepository[T]) GetCollection() *mongo.Collection {
	return r.collection
}

// GetByID retrieves a document by ID
func (r *BaseRepository[T]) GetByID(ctx context.Context, id primitive.ObjectID) (*T, error) {
	var model T

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// Create inserts a new document
func (r *BaseRepository[T]) Create(ctx context.Context, model *T) error {
	_, err := r.collection.InsertOne(ctx, model)
	return err
}

// Update updates a document by ID
func (r *BaseRepository[T]) Update(ctx context.Context, id primitive.ObjectID, model *T) error {
	(*model).UpdateTimestamp()

	update := bson.M{"$set": model}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete removes a document by ID
func (r *BaseRepository[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteMany removes multiple documents based on filter
func (r *BaseRepository[T]) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// ExistsByID checks whether a document exists by its ID
func (r *BaseRepository[T]) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// BuildFilter creates MongoDB filter from SearchFilter slice
func BuildFilter(filters []dto.SearchFilter) bson.M {
	filter := bson.M{}

	for _, f := range filters {
		if f.Key == "" || f.Value == nil {
			continue
		}

		switch f.Type {
		case "search":
			// Text search using regex
			if str, ok := f.Value.(string); ok && str != "" {
				filter[f.Key] = bson.M{"$regex": str, "$options": "i"}
			}
		case "exact":
			filter[f.Key] = f.Value
		case "filter":
			if str, ok := f.Value.(string); ok {
				// Convert string ID to ObjectID
				if objectID, err := primitive.ObjectIDFromHex(str); err == nil {
					filter[f.Key] = bson.M{"$in": []primitive.ObjectID{objectID}}
				}
			} else {
				filter[f.Key] = f.Value
			}
		default:
			// Default to exact match
			filter[f.Key] = f.Value
		}
	}

	return filter
}

// BuildSort creates MongoDB sort from SortOption slice
func BuildSort(sorts []dto.SortOption) bson.M {
	sort := bson.M{}

	for _, s := range sorts {
		if s.Key == "" {
			continue
		}

		order := s.Order
		if order != 1 && order != -1 {
			order = -1 // Default to descending
		}
		sort[s.Key] = order
	}

	// Default sort if no sort specified
	if len(sort) == 0 {
		sort["created_at"] = -1
	}

	return sort
}

// List retrieves documents with pagination, search/filter, and sorting
func (r *BaseRepository[T]) List(ctx context.Context, opts *dto.ListOptions) (*dto.ListResult[T], error) {
	if opts == nil {
		opts = &dto.ListOptions{}
	}
	if opts.Pagination == nil {
		opts.Pagination = &dto.PaginationOptions{}
	}
	opts.Pagination.SetDefaults()

	// Build filter from search/filter options
	filter := BuildFilter(opts.Filters)

	// Count total documents
	totalItems, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	pagination := dto.CalculatePagination(opts.Pagination.Page, opts.Pagination.PageSize, totalItems)

	// Build sort from sort options
	sort := BuildSort(opts.Sort)

	// Find documents with pagination and sorting
	findOpts := GetPaginationOptions(opts.Pagination)
	findOpts.SetSort(sort)

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode documents
	var results []T
	for cursor.Next(ctx) {
		var model T

		if err := cursor.Decode(&model); err != nil {
			return nil, err
		}

		results = append(results, model)
	}

	return &dto.ListResult[T]{
		Records:    &results,
		Pagination: pagination,
	}, nil
}
