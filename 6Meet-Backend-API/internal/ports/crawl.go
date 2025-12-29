package ports

import "context"

// CrawlService defines the interface for crawling operations
type CrawlService interface {
	StartCrawling(ctx context.Context, filename string) error
	Shutdown()
}
