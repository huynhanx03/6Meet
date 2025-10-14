package main

import (
	"log"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/initialize"
)

func main() {
	if err := initialize.Run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
