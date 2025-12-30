package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/event"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	commonHttp "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/workerpool"
	"go.uber.org/zap"
)

const (
	// CrawlWorkerPoolSize is the size of the worker pool for crawling
	CrawlWorkerPoolSize = 32

	// HTTP Client Configuration
	HTTPClientTimeout         = 30 * time.Minute
	HTTPClientMaxIdleConns    = 100
	HTTPClientIdleConnTimeout = 90 * time.Second
	HTTPClientMaxConnsPerHost = 64

	// Crawling Configuration
	CrawlRetryCount = 3
)

// CrawlJob represents a crawling task
type CrawlJob struct {
	PageName string
}

// CrawlService handles crawling operations
type CrawlService struct {
	producer    ports.UserProducer
	consumer    ports.UserConsumer
	httpPool    *commonHttp.HTTPClientPool
	crawlerPool *workerpool.GenericPool[CrawlJob]
}

// NewCrawlService creates a new crawl service
func NewCrawlService(producer ports.UserProducer, consumer ports.UserConsumer) (*CrawlService, error) {
	// Initialize HTTP Client Pool
	httpConfig := &commonHttp.HTTPClientConfig{
		Timeout:         HTTPClientTimeout,
		MaxIdleConns:    HTTPClientMaxIdleConns,
		IdleConnTimeout: HTTPClientIdleConnTimeout,
		MaxConnsPerHost: HTTPClientMaxConnsPerHost,
	}

	httpPool := commonHttp.NewHTTPClientPool(httpConfig)

	s := &CrawlService{
		producer: producer,
		consumer: consumer,
		httpPool: httpPool,
	}

	// Initialize Worker Pool with specialized task function
	pool, err := workerpool.NewGenericPool(CrawlWorkerPoolSize, s.processCrawlJob)
	if err != nil {
		return nil, fmt.Errorf("failed to create worker pool: %w", err)
	}
	s.crawlerPool = pool

	return s, nil
}

// processCrawlJob handles the logic for a single crawl job
func (s *CrawlService) processCrawlJob(job CrawlJob) {
	ctx := context.Background() // Or pass context in job if cancellation needed per task
	pageName := job.PageName

	links, err := s.fetchLinks(ctx, pageName)
	if err != nil {
		global.Logger.Error("Failed to fetch links", zap.String("page", pageName), zap.Error(err))
		return
	}

	// Create Event
	evt := &event.UserCrawled{
		Name:      pageName,
		Neighbors: links,
		Timestamp: time.Now(),
	}

	// Publish Event
	if err := s.producer.Publish(ctx, evt); err != nil {
		global.Logger.Error("Failed to publish event", zap.String("page", pageName), zap.Error(err))
	} else {
		// global.Logger.Info("Published user crawled event", zap.String("page", pageName), zap.Int("neighbors", len(links)))
	}
}

// StartCrawling starts the consumer and initiates the crawling process
func (s *CrawlService) StartCrawling(ctx context.Context, filename string) error {
	global.Logger.Info("Starting crawling process orchestration")

	// Start Consumer
	if err := s.consumer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}
	// defer s.consumer.Stop()

	// Perform Crawl
	if err := s.CrawlFromFile(ctx, filename); err != nil {
		return fmt.Errorf("crawl failed: %w", err)
	}

	return nil
}

// ApiResponse represents the structure of the Wikipedia API response
type ApiResponse struct {
	Query struct {
		Pages map[string]struct {
			Links []struct {
				Title string `json:"title"`
				Ns    int    `json:"ns"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
	Continue struct {
		Plcontinue string `json:"plcontinue"`
	} `json:"continue"`
}

func (s *CrawlService) fetchLinks(ctx context.Context, pageTitle string) ([]string, error) {
	baseURL := "https://en.wikipedia.org/w/api.php?action=query&prop=links&format=json&plnamespace=0&pllimit=max&titles=%s"
	var allLinks []string
	plcontinue := ""

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		encodedTitle := url.QueryEscape(pageTitle)
		requestURL := fmt.Sprintf(baseURL, encodedTitle)
		if plcontinue != "" {
			requestURL += "&plcontinue=" + url.QueryEscape(plcontinue)
		}

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("User-Agent", "6MeetBot/1.0")

		resp, err := s.httpPool.RequestWithRetry(ctx, req, CrawlRetryCount)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		var result ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		}

		for _, page := range result.Query.Pages {
			for _, link := range page.Links {
				if link.Ns == 0 {
					allLinks = append(allLinks, link.Title)
				}
			}
		}

		plcontinue = result.Continue.Plcontinue
		if plcontinue == "" {
			return allLinks, nil
		}
	}
}

// CrawlFromFile reads a file and crawls pages concurrently
func (s *CrawlService) CrawlFromFile(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	pagesChan := make(chan string, 1000)
	go func() {
		defer close(pagesChan)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case pagesChan <- scanner.Text():
			}
		}
		if err := scanner.Err(); err != nil {
			global.Logger.Error("Error reading file", zap.Error(err))
		}
	}()

	for page := range pagesChan {
		if page == "" {
			continue
		}

		err := s.crawlerPool.Invoke(CrawlJob{
			PageName: page,
		})

		if err != nil {
			global.Logger.Error("Failed to submit crawl task", zap.String("page", page), zap.Error(err))
		}
	}

	global.Logger.Info("Successfully submitted all crawl tasks from file", zap.String("file", filename))
	return nil
}

func (s *CrawlService) Shutdown() {

}
