package entity

import "time"

// User represents the domain entity for a user
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Neighbors []string  `json:"neighbors"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
