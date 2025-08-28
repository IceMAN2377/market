# Subscription Management Service

A RESTful API service for managing user subscriptions to various services. This application allows creating, retrieving, updating, and deleting subscription records, as well as calculating the total cost of subscriptions over a specified period.

## Features

- CRUD operations for subscription management
- Filtering subscriptions by user ID and service name
- Calculating total subscription costs for a specified period
- Pagination support for listing subscriptions
- Containerized deployment with Docker and Docker Compose
- Database migrations using golang-migrate

## Technologies Used

- Go 1.24
- PostgreSQL 15
- Docker & Docker Compose
- RESTful API design
- JSON for data interchange
- Structured logging with slog
- Database migrations with golang-migrate
- Environment-based configuration

## Project Structure

```
market/
├── app/                  # Application initialization and setup
├── cmd/                  # Command-line entry points
│   └── market/           # Main application entry point
├── db/                   # Database related files
│   └── migrations/       # SQL migration files
├── internal/             # Internal application code
│   ├── config/           # Configuration handling
│   ├── errors/           # Custom error definitions
│   ├── models/           # Data models and DTOs
│   ├── repository/       # Data access layer
│   │   └── postgres/     # PostgreSQL implementation
│   ├── service/          # Business logic layer
│   │   └── subscription/ # Subscription service implementation
│   └── transport/        # Transport layer
│       └── http/         # HTTP handlers and routing
├── .env                  # Environment variables for local development
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile            # Docker image definition
├── go.mod                # Go module definition
└── go.sum                # Go module checksums
```

## Installation and Setup

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (for containerized deployment)

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/IceMAN2377/market.git
   cd market
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the environment variables (or use the provided .env file):
   ```
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=myuser
   POSTGRES_PASSWORD=mypassword
   POSTGRES_DB=market
   POSTGRES_SSL_MODE=disable
   POSTGRES_MIGRATE=true
   HTTP_PORT=8080
   LOG_LEVEL=info
   ```

4. Start PostgreSQL (if not using Docker):
   ```bash
   # Configure your PostgreSQL instance with the credentials from .env
   ```

5. Run the application:
   ```bash
   go run cmd/market/main.go
   ```

### Docker Deployment

The easiest way to run the application is using Docker Compose:

```bash
docker-compose up -d
```

This will start both the PostgreSQL database and the application. The API will be available at http://localhost:8080.

To stop the services:

```bash
docker-compose stop
```

To stop and remove the containers:

```bash
docker-compose down
```

## API Endpoints

### Subscriptions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/subscriptions` | Create a new subscription |
| GET | `/api/v1/subscriptions` | List subscriptions with optional filtering |
| GET | `/api/v1/subscriptions/{id}` | Get a specific subscription by ID |
| PUT | `/api/v1/subscriptions/{id}` | Update a subscription |
| DELETE | `/api/v1/subscriptions/{id}` | Delete a subscription |
| POST | `/api/v1/subscriptions/cost-calculation` | Calculate subscription costs for a period |

### Request/Response Examples

#### Create Subscription

Request:
```
POST /api/v1/subscriptions

{
  "service_name": "Netflix",
  "price": 1499,
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "start_date": "2023-01-01"
}
```

Response:
```
Status: 201 Created

{
  "id": 1,
  "service_name": "Netflix",
  "price": 1499,
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "start_date": "2023-01-01",
  "end_date": null,
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

#### Calculate Cost

Request:
```
POST /api/v1/subscriptions/cost-calculation

{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "start_date": "2023-01-01",
  "end_date": "2023-12-31"
}
```

Response:
```
Status: 200 OK

{
  "total_cost": 17988,
  "start_date": "2023-01-01",
  "end_date": "2023-12-31",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "service_name": null
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| POSTGRES_HOST | PostgreSQL host | localhost |
| POSTGRES_PORT | PostgreSQL port | 5432 |
| POSTGRES_USER | PostgreSQL username | myuser |
| POSTGRES_PASSWORD | PostgreSQL password | mypassword |
| POSTGRES_DB | PostgreSQL database name | market |
| POSTGRES_SSL_MODE | PostgreSQL SSL mode | disable |
| POSTGRES_MIGRATE | Whether to run migrations on startup | true |
| HTTP_PORT | HTTP server port | 8080 |
| LOG_LEVEL | Logging level (debug, info, warn, error) | info |

## Development

### Database Migrations

Migrations are automatically applied on application startup if `POSTGRES_MIGRATE=true`. 

To manually apply migrations:

```bash
# Install golang-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply migrations
migrate -path db/migrations -database "postgres://myuser:mypassword@localhost:5432/market?sslmode=disable" up
```

To create a new migration:

```bash
migrate create -ext sql -dir db/migrations -seq migration_name
```