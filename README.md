# Khalif Backend API

A production-ready Islamic content platform backend built with Go, following Clean Architecture principles.

## 🚀 Features

- **Authentication**: Register, Login, Logout with JWT tokens
- **OTP Email Verification**: Brevo integration for email verification
- **Forgot Password**: Secure password reset via email token
- **Rolling Refresh Tokens**: Secure session management with automatic token rotation
- **Google OAuth**: Login with Google account
- **Rate Limiting**: Redis-based rate limiting for auth endpoints
- **Account Lockout**: Automatic account lock after failed login attempts
- **Listening History**: Track user listening with SP-based spam prevention
- **Playlist**: Create, manage, and share audio playlists
- **Meilisearch**: Fast, typo-tolerant search across all content
- **Prayer Times**: Accurate prayer schedule with Qibla direction
- **Hadist & Doa**: Islamic content with like, bookmark, and audio support
- **Profile Picture Upload**: Multipart form upload with automatic fallback
- **Health Checks**: `/health` and `/ready` endpoints for container orchestration
- **Graceful Shutdown**: Proper cleanup of database and Redis connections

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

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
```

## 🚀 Running the Application

```bash
make dev          # Start services and run app
make stop-services # Stop all services
make status       # Check service status
make run          # Run without starting services
make build        # Build binary
make test         # Run tests
```

## 📡 API Endpoints

### Health & Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |

### Admin Auth (`/api/v1/auth`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register new admin |
| POST | `/login` | Login |
| POST | `/refresh` | Refresh token |
| GET | `/me` | Get current admin 🔒 |
| POST | `/logout` | Logout 🔒 |

### User Auth (`/api/v1/users/auth`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register (sends OTP) |
| POST | `/login` | Login |
| POST | `/google` | Login with Google |
| POST | `/verify-otp` | Verify OTP |
| POST | `/resend-otp` | Resend OTP |
| POST | `/forgot-password` | Request reset |
| POST | `/reset-password` | Reset password |
| POST | `/refresh` | Refresh token |
| GET | `/me` | Get current user 🔒 |
| POST | `/logout` | Logout 🔒 |

### User Features (`/api/v1/users`) 🔒

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/update` | Update profile |
| POST | `/audio/:id/listen` | Record listening |
| GET | `/listening-history` | Get history |
| POST | `/audio/:id/like` | Like audio |
| DELETE | `/audio/:id/like` | Unlike audio |
| GET | `/likes` | Get liked audios |

### Playlist (`/api/v1/users/playlist`) 🔒

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/` | Create playlist |
| GET | `/` | Get my playlists |
| PUT | `/:id` | Update |
| DELETE | `/:id` | Delete |
| POST | `/:id/audio/:audio_id` | Add audio |
| DELETE | `/:id/audio/:audio_id` | Remove audio |
| POST | `/:id/like` | Like playlist |
| DELETE | `/:id/like` | Unlike |

### Search (`/api/v1/search`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/?q=query` | Unified search |
| GET | `/audio?q=query` | Search audios |
| GET | `/ustadz?q=query` | Search ustadzs |
| GET | `/mood?q=query` | Search moods |
| GET | `/playlist?q=query` | Search playlists |
| GET | `/doa?q=query` | Search doa |
| GET | `/hadist?q=query` | Search hadist |

### Radio (`/api/v1/radio`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/:id?limit=20` | Generate radio queue |

### Prayer Times (`/api/v1/prayer-times`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/?lat={lat}&long={long}` | Get prayer schedule |
| GET | `/daily?lat={lat}&long={long}` | Get daily schedule |

**Response includes:**
- Prayer schedule (Imsak, Subuh, Syuruq, Dhuha, Dzuhur, Ashar, Maghrib, Isya)
- Qibla direction (degrees from North)
- Countdown to next prayer
- Auto timezone detection

### Hadist (`/api/v1/hadist`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Get all hadists |
| GET | `/:id` | Get detail |
| GET | `/random` | Get random |
| GET | `/category?category=X` | By category |
| GET | `/kitab?kitab=X` | By kitab |
| POST | `/:id/listen` | Increment count |

**User (Protected):** `POST /api/v1/users/hadist/:id/like`, `/bookmark`

**Admin (Protected):** `POST/PUT/DELETE /api/v1/admin/hadist`

### Doa (`/api/v1/doa`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Get all doa |
| GET | `/:id` | Get detail |
| GET | `/random` | Get random |
| GET | `/category?category=X` | By category |
| GET | `/hadist?hadist_id=X` | By hadist |
| POST | `/:id/listen` | Increment count |

**User (Protected):** `POST /api/v1/users/doa/:id/like`, `/bookmark`

**Admin (Protected):** `POST/PUT/DELETE /api/v1/admin/doa`

## 🔒 Security Features

- **Password Hashing**: bcrypt
- **JWT**: HS256 signed tokens
- **OTP Verification**: 6-digit code, 10 min expiry
- **Password Reset**: Hashed tokens, 30 min expiry
- **Rate Limiting**: 5 req/sec per IP
- **Account Lockout**: 5 failed → 30 min lock
- **Refresh Token Rotation**: Old tokens revoked

## 📦 Stored Procedures

- `sp_handle_login_failure`: Login throttling
- `sp_check_lock_status`: Account lock check
- `sp_revoke_user_tokens`: Revoke tokens
- `sp_record_listening`: Listening with spam prevention
- `sp_like_playlist`: Like with atomic increment
- `sp_unlike_playlist`: Unlike with atomic decrement
- `sp_record_playlist_listening`: Playlist listening
- `sp_add_audio_to_playlist`: Add audio with auto-position

## 🧪 Testing

```bash
go test -v ./...              # All tests
go test -cover ./...          # With coverage
go test -v ./tests/e2e/...    # E2E only
```

## 📄 License

MIT
