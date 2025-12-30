package event

import "time"

// UserCrawled represents the event when a user page is crawled
type UserCrawled struct {
	Name      string    `json:"name"`
	Neighbors []string  `json:"neighbors"`
	Timestamp time.Time `json:"timestamp"`
}
