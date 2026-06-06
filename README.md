# Wallet API

A RESTful wallet service built with Go and Gin. Supports OTP-based passwordless authentication, per-user wallets, and peer-to-peer transactions.

## Stack

- **Go** (1.25) + **Gin** — HTTP framework
- **PostgreSQL** — primary database (via GORM + pgx)
- **golang-migrate** — SQL migrations
- **JWT** — stateless auth tokens
- **gomail** — OTP delivery via SMTP
- **Docker Compose** — local database setup
- **Air** — live reload for development

## Project Structure

```
.
├── bootstrap/          # App wiring (DB, router, server)
├── database/
│   └── migrations/     # SQL migration files
├── internal/
│   ├── handler/        # HTTP handlers (users, wallets, transactions)
│   ├── middleware/     # JWT auth middleware
│   ├── service/        # Business logic layer
│   └── store/          # Data access layer (GORM)
├── postman/            # Postman collection + environments
├── provider/
│   ├── jwt/            # Token generation & validation
│   └── mailer/         # SMTP email provider
└── route/              # Route registration
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### 1. Start the database

```bash
docker compose up -d
```

### 2. Configure environment

Copy `.env` and adjust values as needed:

```bash
cp .env .env.local
```

Key variables:

| Variable | Description |
|---|---|
| `ADDR` | Server listen address (default `:8030`) |
| `DB_ADDR` | PostgreSQL connection string |
| `JWT_SECRET_KEY` | Secret used to sign JWT tokens |
| `AUTH_EXPIRATION` | Token TTL (e.g. `30d`) |
| `SMTP_HOST` / `SMTP_PORT` | SMTP server for OTP emails |
| `SMTP_USERNAME` / `SMTP_PASSWORD` | SMTP credentials |

### 3. Run migrations

```bash
make migrate-up
```

### 4. Start the server

```bash
# With live reload (requires Air)
air

# Or directly
go run bootstrap/main.go
```

The server starts on `http://localhost:8030`.

## API Reference

All protected routes require the header:

```
Authorization: Bearer <token>
```

### Auth (Users)

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/users/otp` | No | Request OTP for a username |
| `POST` | `/users/otp/verify` | No | Verify OTP and receive JWT |
| `GET` | `/users/me` | Yes | Get current user profile |

**Request OTP**
```json
POST /users/otp
{ "username": "alice" }
```

**Verify OTP**
```json
POST /users/otp/verify
{ "username": "alice", "code": "123456" }
```
Response: `{ "token": "<jwt>" }`

### Wallets

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/wallets/me` | Yes | Get authenticated user's wallet |

### Transactions

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/transactions/` | Yes | Send funds to another user |
| `GET` | `/transactions/:id` | Yes | Get transaction by ID |
| `GET` | `/transactions/user/:userID` | Yes | List user transactions (paginated) |

**Create Transaction**
```json
POST /transactions/
{
  "receiver_id": 2,
  "amount": 50.00,
  "currency": "USD"
}
```

**List Transactions** — query params: `page`, `pageSize`, `phone`, `email`, `username`

```
GET /transactions/user/1?page=1&pageSize=10
```

### Health Check

```
GET /ping   →  { "message": "pong" }   (requires auth)
```

## Database Migrations

```bash
# Create a new migration
make migrate-create name=add_column_to_users

# Apply all pending migrations
make migrate-up

# Roll back N migrations
make migrate-down 1
```

## Postman

Import the collection and environment from the `postman/` directory:

- `wallet.postman_collection.json` — all endpoints
- `wallet.local.postman_environment.json` — local environment
- `wallet.staging.postman_environment.json` — staging environment

## Database Schema

```
users
  id, username (unique), email, phone, status, timestamps

wallets
  id, user_id (FK → users), balance, currency, timestamps

transactions
  id, sender_id (FK → users), receiver_id (FK → users), amount, status, currency, timestamps

otps
  (stores pending OTP codes tied to users)
```
