# Birthday Gift Helper

A mobile app to help people choose the perfect birthday gift. Register recipients with their interests, age, and budget — and get personalized gift suggestions.

## Tech Stack

- **Backend**: Go 1.24+ with Chi v5
- **Frontend**: React Native with Expo (SDK 52+)
- **Database**: PostgreSQL 16
- **Cache**: Redis 7

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose

### 1. Start infrastructure

```bash
make infra
```

### 2. Run migrations

```bash
DATABASE_URL=postgres://birthday:birthday_secret@localhost:5432/birthday_db?sslmode=disable make migrate-up
```

### 3. Seed demo data (optional)

```bash
DATABASE_URL=postgres://birthday:birthday_secret@localhost:5432/birthday_db?sslmode=disable make seed
```

Demo user: `demo@example.com` / `password123`

### 4. Start backend

```bash
cp backend/.env.example backend/.env
make dev-backend
```

### 5. Start mobile app

```bash
make mobile-install
make dev-mobile
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make infra` | Start PostgreSQL + Redis via Docker |
| `make dev-backend` | Run Go backend on :8080 |
| `make dev-mobile` | Start Expo dev server |
| `make migrate-up` | Run all pending migrations |
| `make migrate-down` | Rollback all migrations |
| `make migrate-create name=xxx` | Create a new migration |
| `make seed` | Seed database with demo data |
| `make test-backend` | Run all backend tests |
| `make stop` | Stop all Docker services |

## API Endpoints

### Auth (Public)
- `POST /api/auth/register` — Register with email + password
- `POST /api/auth/login` — Login with email + password
- `POST /api/auth/google` — Google Sign-In
- `POST /api/auth/apple` — Apple Sign-In
- `POST /api/auth/refresh` — Refresh access token

### Protected (Bearer JWT)
- `GET /api/auth/me` — Get current user
- `POST /api/recipients` — Create recipient
- `GET /api/recipients` — List all recipients
- `GET /api/recipients/:id` — Get recipient
- `PUT /api/recipients/:id` — Update recipient
- `DELETE /api/recipients/:id` — Delete recipient
- `DELETE /api/recipients` — Bulk delete recipients
