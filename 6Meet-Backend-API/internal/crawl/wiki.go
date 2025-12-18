package crawl

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/goroutine"
	commonHttp "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http"
	"go.uber.org/zap"
)

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

// CrawlTask represents a single crawl operation
type CrawlTask struct {
	name     string
	httpPool *commonHttp.HTTPClientPool
}

// Process implements the Task interface for WorkerPool
func (t *CrawlTask) Process(ctx context.Context) (*dto.CreateUserRequest, error) {
	links, err := t.fetchLinks(ctx, t.name)
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserRequest{
		Name:      t.name,
		Neighbors: links,
	}, nil
}

// fetchLinks fetches all links from a Wikipedia page
func (t *CrawlTask) fetchLinks(ctx context.Context, pageTitle string) ([]string, error) {
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

		resp, err := t.httpPool.RequestWithRetry(ctx, req, 3)
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
			return allLinks, nil // Return here when no more pages
		}

		// time.Sleep(100 * time.Millisecond)
	}
}

// createHTTPPool creates and configures the HTTP client pool
func createHTTPPool() *commonHttp.HTTPClientPool {
	config := &commonHttp.HTTPClientConfig{
		Timeout:         30 * time.Minute,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		MaxConnsPerHost: 64,
	}
	return commonHttp.NewHTTPClientPool(config)
}

// createWorkerPool creates and starts a worker pool
func createWorkerPool(ctx context.Context) *goroutine.WorkerPool[*dto.CreateUserRequest] {
	pool := goroutine.NewIOExecutor[*dto.CreateUserRequest](ctx)
	return pool
}

// readPagesFromFile reads pages from a file and sends them to a channel
func readPagesFromFile(ctx context.Context, filename string, pagesChan chan<- string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		global.Logger.Error("Failed to open file",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return err
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pagesChan <- scanner.Text():
		}
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		global.Logger.Error("Error reading file",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// processResults handles the results from the worker pool
func processResults(resultsC <-chan *dto.CreateUserRequest, errorsC <-chan error, wg *sync.WaitGroup) {
	wg.Add(2)

	// Process successful results
	go func() {
		defer wg.Done()
		for result := range resultsC {
			global.Logger.Info("Processed page",
				zap.String("name", result.Name),
				zap.Int("neighbors", len(result.Neighbors)),
			)
		}
	}()

	// Process errors
	go func() {
		defer wg.Done()
		for err := range errorsC {
			global.Logger.Error("Crawl error", zap.Error(err))
		}
	}()
}

// submitTasks submits pages to the worker pool for processing
func submitTasks(pagesChan <-chan string, workerPool *goroutine.WorkerPool[*dto.CreateUserRequest], httpPool *commonHttp.HTTPClientPool) error {
	for page := range pagesChan {
		if page == "" {
			continue
		}

		task := &CrawlTask{
			name:     page,
			httpPool: httpPool,
		}

		// Submit blocks until a worker receives the task or context is done
		if err := workerPool.Submit(task); err != nil {
			global.Logger.Error("Failed to submit task",
				zap.String("page", page),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}

// CrawlData processes pages from a file concurrently
func CrawlData(filename string) error {
	start := time.Now()
	global.Logger.Info("Crawling started")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Initialize components
	httpPool := createHTTPPool()
	workerPool := createWorkerPool(ctx)
	workerPool.Start()

	// Start reading pages from file
	pagesChan := make(chan string, 1000)

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(pagesChan)
		if err := readPagesFromFile(ctx, filename, pagesChan); err != nil {
			global.Logger.Error("Failed to read pages",
				zap.String("filename", filename),
				zap.Error(err),
			)
			cancel()
		}
	}()

	// Process results
	resultsC, errorsC := workerPool.Results()
	processResults(resultsC, errorsC, &wg)

	// Wait for worker pool to finish
	defer func() {
		workerPool.Shutdown()
		global.Logger.Info(fmt.Sprintf("Crawling completed in %s", time.Since(start).Round(time.Second)))
	}()

	// Submit tasks to the worker pool
	if err := submitTasks(pagesChan, workerPool, httpPool); err != nil {
		return fmt.Errorf("error submitting tasks: %w", err)
	}

	return nil
}
