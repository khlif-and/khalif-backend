# Khalif Backend Auth Service

A production-ready authentication microservice built with Go, following Clean Architecture principles.

## 🚀 Features

- **Authentication**: Register, Login, Logout with JWT tokens
- **Rolling Refresh Tokens**: Secure session management with automatic token rotation
- **Rate Limiting**: Redis-based rate limiting for auth endpoints
- **Account Lockout**: Automatic account lock after failed login attempts
- **Profile Picture Upload**: Multipart form upload with automatic fallback to initials avatar
- **Health Checks**: `/health` and `/ready` endpoints for container orchestration
- **Graceful Shutdown**: Proper cleanup of database and Redis connections
- **CORS**: Cross-origin resource sharing enabled

## 📁 Project Structure

```
├── cmd/api/              # Application entrypoint
├── internal/
│   ├── adapters/
│   │   ├── handlers/     # HTTP handlers
│   │   ├── http/         # Router configuration
│   │   └── repositories/ # Data access layer
│   ├── core/
│   │   ├── domain/       # Entities & DTOs
│   │   ├── ports/        # Interfaces (contracts)
│   │   └── services/     # Business logic
│   └── platform/
│       ├── config/       # Configuration
│       ├── database/     # DB & Redis connections
│       └── logger/       # Zap logger
├── migrations/sql/       # SQL migrations & stored procedures
├── pkg/
│   ├── messages/         # Error & success messages
│   ├── middleware/       # Auth, CORS, Rate limiting
│   └── utils/            # JWT, Password, Upload utilities
├── tests/
│   ├── e2e/              # End-to-end tests
│   └── integration/      # Integration tests
└── uploads/              # Uploaded files
```

## 🛠️ Tech Stack

- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Auth**: JWT (Access + Refresh tokens)
- **Logging**: Zap
- **Testing**: testify

## ⚙️ Configuration

Create a `.env` file:

```env
SERVER_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=khalif_backend
DB_PORT=5432
DB_SSLMODE=disable

# JWT
JWT_SECRET=your_super_secret_key
JWT_EXP_HOURS=24
REFRESH_TOKEN_SECRET=your_refresh_secret_key
REFRESH_TOKEN_EXP_DAYS=7

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## 🚀 Running the Application

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go

# Run tests
go test ./...

# Run specific test suites
go test ./tests/e2e/...
go test ./tests/integration/...
```

## 📡 API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check (DB & Redis) |
| POST | `/api/v1/auth/register` | Register new admin (multipart) |
| POST | `/api/v1/auth/login` | Login (multipart) |
| POST | `/api/v1/auth/refresh` | Refresh access token (JSON) |

### Protected Endpoints (Bearer Token Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/me` | Get current user |
| POST | `/api/v1/auth/logout` | Logout |
| PUT | `/api/v1/admin/update` | Update profile |

## 🔒 Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT**: HS256 signed tokens
- **Rate Limiting**: 5 requests/second per IP on auth endpoints
- **Account Lockout**: 5 failed attempts → 30 minute lock
- **Refresh Token Rotation**: Old tokens revoked on refresh

## 📦 Database Migrations

Migrations run automatically on startup via GORM AutoMigrate and SQL files in `migrations/sql/`.

### Stored Procedures

- `sp_handle_login_failure`: Handles login throttling
- `sp_check_lock_status`: Checks account lock status
- `sp_revoke_user_tokens`: Revokes all user refresh tokens

## 🧪 Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run E2E tests only
go test -v ./tests/e2e/...

# Run integration tests only
go test -v ./tests/integration/...
```

## 📄 License

MIT
