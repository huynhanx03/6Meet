# 6Meet Backend API

## Tech Stack

- **Language**: [Go](https://go.dev/) (Golang)
- **Framework**: [Gin](https://gin-gonic.com/)
- **Database**: [MongoDB](https://www.mongodb.com/)
- **Caching**: [Redis](https://redis.io/)
- **Message Queue**: [Kafka](https://kafka.apache.org/)

## Architecture & Project Structure

This project follows **Hexagonal Architecture** (Ports and Adapters) to ensure domain logic independence.

```
├── cmd/                # Application entry points
├── config/             # Configuration files
├── internal/
│   ├── adapters/       # Adapters Layer
│   │   ├── driver/     # Primary adapters (e.g., HTTP Handlers)
│   │   └── driven/     # Secondary adapters (e.g., MongoDB, Redis)
│   ├── core/           # Core Domain Layer
│   │   ├── service/    # Application Business Rules
│   │   └── dto/        # Data Transfer Objects
│   │   └── entity/     # Enterprise Business Rules (Pure Entities)
│   ├── ports/          # Interfaces (Ports) defining boundaries
│   ├── initialize/     # Bootstrap logic & Dependency Injection
│   ├── mapper/         # Object Mappers
│   └── constant/       # Internal constants
├── pkg/                # Shared utilities and libraries
└── Makefile            # Build and run commands
```

## API Documentation

API endpoints are available at `/api/v1`.
- **User APIs**: `/api/v1/users`
