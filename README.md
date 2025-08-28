# RUN THE PROJECT WITH FOLLOWING COMMAND
```bash
docker-compose up
```

### Subscriptions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/subscriptions` | Create a new subscription |
| GET | `/api/v1/subscriptions` | List subscriptions with optional filtering |
| GET | `/api/v1/subscriptions/{id}` | Get a specific subscription by ID |
| PUT | `/api/v1/subscriptions/{id}` | Update a subscription |
| DELETE | `/api/v1/subscriptions/{id}` | Delete a subscription |
| POST | `/api/v1/subscriptions/cost-calculation` | Calculate subscription costs for a period |

### Swagger endpoint
| Method | Endpoint | Description           |
|--------|---------|-----------------------|
| POST | `/docs` | Swagger documentation |


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
   

