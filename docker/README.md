# Docker Development Environment

This directory contains Docker configurations for local development of the Calendar API.

## Quick Start

### Start the Database Only
```bash
# Start just the PostgreSQL database
task db:start

# Or with docker-compose directly
docker-compose -f docker/docker-compose.yml up -d db
```

### Start Full Development Environment
```bash
# Start both database and API
task dev

# Or with docker-compose directly
docker-compose -f docker/docker-compose.yml up --build
```

## Database Configuration

### Ports
- **Main Database**: `localhost:5484` (external) → `5432` (internal)
- **Test Database**: `localhost:5485` (external) → `5432` (internal) (when using `--profile test`)
- **API Server**: `localhost:8012` (external) → `8000` (internal)

### Default Credentials
- **Main DB**: `calendar_user` / `calendar_pass` / `calendar_db`
- **Test DB**: `test_user` / `test_pass` / `calendar_test_db`

### Environment Variables
You can override any of these in a `.env` file:

```bash
# Database ports
DB_PORT=5484
TEST_DB_PORT=5485
API_PORT=8012

# Database credentials
POSTGRES_USER=calendar_user
POSTGRES_PASSWORD=calendar_pass
POSTGRES_DB=calendar_db

# Test database credentials
TEST_POSTGRES_USER=test_user
TEST_POSTGRES_PASSWORD=test_pass
TEST_POSTGRES_DB=calendar_test_db

# Application settings
DEBUG=true
ENVIRONMENT=development
BOOTSTRAP_ADMIN_KEY=dev-test-key-123
```

## Database Management

### Available Task Commands
```bash
task db:start          # Start main database
task db:start:test      # Start both main and test databases
task db:stop            # Stop all services
task db:logs            # Show database logs
task db:connect         # Connect to main database with psql
task db:connect:test    # Connect to test database with psql
task db:reset           # Reset database (removes all data!)
task db:migrate:local   # Run migrations against local database
```

### Manual Database Connection
```bash
# Connect to main database
psql -h localhost -p 5484 -U calendar_user -d calendar_db

# Connect to test database (when running)
psql -h localhost -p 5485 -U test_user -d calendar_test_db

# From inside the container
docker exec -it calendar-api-postgres psql -U calendar_user -d calendar_db
```

## Development Workflow

### 1. Start Database
```bash
task db:start
```

### 2. Run Migrations
```bash
task db:migrate:local
```

### 3. Start API
```bash
# With Doppler (staging environment)
task dev:local

# Or locally without Doppler
cd src && go run .
```

### 4. Test the API
```bash
# Health check
curl http://localhost:8000/health

# Version check
curl http://localhost:8000/version

# API endpoints (requires X-API-Key header)
curl -H "X-API-Key: dev-test-key-123" http://localhost:8000/api/events
```

## Test Database

The test database is available but only starts when explicitly requested:

```bash
# Start with test database
docker-compose -f docker/docker-compose.yml --profile test up -d

# Or using task
task db:start:test
```

This is useful for:
- Running integration tests
- Testing migrations
- Isolating test data from development data

## Persistence

Database data is persisted in Docker volumes:
- `calendar-api-3dead8d4_postgres_data` (main database)
- `calendar-api-3dead8d4_postgres_test_data` (test database)

To completely reset and remove all data:
```bash
task db:reset
```

## Initialization Scripts

Custom SQL scripts in `./init-db/` will be executed when the database container is first created. The included `01-init.sql` sets up:
- UUID extension
- pg_stat_statements extension  
- Performance logging
- Read-only user for monitoring 