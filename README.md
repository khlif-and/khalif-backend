# Khalif Backend Auth Service

A production-ready authentication microservice built with Go, following Clean Architecture principles.

## 🚀 Features

- **Authentication**: Register, Login, Logout with JWT tokens
- **OTP Email Verification**: Brevo integration for email verification before login
- **Forgot Password**: Secure password reset via email token
- **Rolling Refresh Tokens**: Secure session management with automatic token rotation
- **Rate Limiting**: Redis-based rate limiting for auth endpoints
- **Account Lockout**: Automatic account lock after failed login attempts
- **Listening History**: Track user listening with SP-based spam prevention
- **Playlist**: Create, manage, and share audio playlists
- **Meilisearch**: Fast, typo-tolerant search across all content
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
│   ├── infrastructure/
│   │   ├── email/        # Brevo email service
│   │   └── search/       # Meilisearch integration
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
- **Search**: Meilisearch
- **Auth**: JWT (Access + Refresh tokens)
- **Email**: Brevo (Sendinblue)
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

# Brevo Email Service
BREVO_API_KEY=your_brevo_api_key
BREVO_SENDER_EMAIL=your@email.com
BREVO_SENDER_NAME=Khalif App

# Meilisearch
MEILISEARCH_HOST=http://localhost:7700
MEILISEARCH_API_KEY=your_master_key
```

## 🚀 Running the Application

```bash
# Start all services (PostgreSQL, Redis, Meilisearch) and run the app
make dev

# Stop all services
make stop-services

# Check service status
make status

# Run without starting services
make run

# Build binary
make build

# Run tests
make test
```

## 📡 API Endpoints

### Health & Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check (DB & Redis) |

### Admin Auth (`/api/v1/auth`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register new admin (multipart) |
| POST | `/login` | Login (multipart) |
| POST | `/refresh` | Refresh access token |
| GET | `/me` | Get current admin 🔒 |
| POST | `/logout` | Logout 🔒 |

### User Auth (`/api/v1/users/auth`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register (sends OTP email) |
| POST | `/login` | Login (requires verified account) |
| POST | `/verify-otp` | Verify OTP & activate account |
| POST | `/resend-otp` | Resend OTP to email |
| POST | `/forgot-password` | Request password reset |
| POST | `/reset-password` | Reset password with token |
| POST | `/refresh` | Refresh access token |
| GET | `/me` | Get current user 🔒 |
| POST | `/logout` | Logout 🔒 |

### User Features (`/api/v1/users`) 🔒

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/update` | Update profile |
| POST | `/audio/:id/listen` | Record listening (SP) |
| GET | `/listening-history` | Get listening history |
| POST | `/audio/:id/like` | Like audio |
| DELETE | `/audio/:id/like` | Unlike audio |
| GET | `/audio/:id/is-liked` | Check if liked |
| GET | `/likes` | Get user's liked audios |

### Playlist (`/api/v1/users/playlist`) 🔒

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/` | Create playlist (multipart) |
| GET | `/` | Get my playlists |
| PUT | `/:id` | Update playlist |
| DELETE | `/:id` | Delete playlist |
| POST | `/:id/audio/:audio_id` | Add audio to playlist |
| DELETE | `/:id/audio/:audio_id` | Remove audio from playlist |
| POST | `/:id/like` | Like playlist |
| DELETE | `/:id/like` | Unlike playlist |
| GET | `/:id/is-liked` | Check if liked |

### Playlist Public (`/api/v1/playlist`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Get all public playlists |
| GET | `/:id` | Get playlist detail with audios |
| POST | `/:id/listen` | Increment listening count |

### Search (`/api/v1/search`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/?q=query` | Unified search (all content) |
| GET | `/audio?q=query` | Search audios only |
| GET | `/ustadz?q=query` | Search ustadzs only |
| GET | `/mood?q=query` | Search mood categories only |
| GET | `/playlist?q=query` | Search playlists only |

### Radio (`/api/v1/radio`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/:id?limit=20` | Generate radio queue based on seed audio |

**Radio Algorithm**: Returns similar audios scored by same Ustadz (+3), same Mood (+2), popularity (+1).

### Audio, Mood Categories, Ustadz

Public GET endpoints and admin-only CUD operations. See router for details.

### Hadist (`/api/v1/hadist`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Get all hadists |
| GET | `/:id` | Get hadist detail |
| GET | `/random` | Get random hadist |
| GET | `/category?category=X` | Get by category |
| GET | `/kitab?kitab=X` | Get by kitab |
| POST | `/:id/listen` | Increment listening count |

**User Engagement (Protected):**
- `POST /api/v1/users/hadist/:id/like`
- `POST /api/v1/users/hadist/:id/bookmark`

**Admin (Protected):**
- `POST /api/v1/admin/hadist` (Create)
- `PUT /api/v1/admin/hadist/:id` (Update)
- `DELETE /api/v1/admin/hadist/:id` (Delete)

### Prayer Times (`/api/v1/prayer-times`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/?lat={lat}&long={long}` | Get daily prayer schedule + countdown |

**Response features:**
- Accurate calculation using Kemenag RI method
- Auto timezone detection based on coordinates
- Realtime `time_remaining` countdown to next prayer
- Handles day transition (Next prayer 'Subuh' tomorrow)

## 🔒 Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT**: HS256 signed tokens
- **OTP Verification**: 6-digit code, 10 min expiry
- **Password Reset**: Hashed tokens, 30 min expiry, single-use
- **Rate Limiting**: 5 requests/second per IP on auth endpoints
- **Account Lockout**: 5 failed attempts → 30 minute lock
- **Refresh Token Rotation**: Old tokens revoked on refresh
- **Session Revocation**: All sessions revoked after password reset

## 📦 Database Migrations

Migrations run automatically on startup via GORM AutoMigrate and SQL files in `migrations/sql/`.

### Stored Procedures

- `sp_handle_login_failure`: Handles login throttling
- `sp_check_lock_status`: Checks account lock status
- `sp_revoke_user_tokens`: Revokes all user refresh tokens
- `sp_record_listening`: Records listening with spam prevention
- `sp_like_playlist`: Like playlist with atomic increment
- `sp_unlike_playlist`: Unlike playlist with atomic decrement
- `sp_record_playlist_listening`: Record playlist listening
- `sp_add_audio_to_playlist`: Add audio with auto-position

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
