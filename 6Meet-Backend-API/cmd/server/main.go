package main

import (
	"log"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/infrastructure"
)

func main() {
	if err := infrastructure.Run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
