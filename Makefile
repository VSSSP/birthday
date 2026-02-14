.PHONY: infra dev-db dev-backend dev-mobile dev-web migrate-up migrate-down migrate-create seed test-backend lint-web build-web web-install mobile-install stop

# ========================
# Infrastructure
# ========================

# Start all infrastructure (PostgreSQL + Redis)
infra:
	docker compose up -d

# Start PostgreSQL only
dev-db:
	docker compose up postgres -d

# Stop all services
stop:
	docker compose down

# ========================
# Backend
# ========================

# Run backend locally
dev-backend:
	cd backend && go run ./cmd/api

# Run backend tests
test-backend:
	cd backend && go test ./... -v

# ========================
# Database
# ========================

# Run all pending migrations
migrate-up:
	cd backend && go run ./cmd/migrate up

# Rollback all migrations
migrate-down:
	cd backend && go run ./cmd/migrate down

# Create a new migration (usage: make migrate-create name=create_something)
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=create_something"; exit 1; fi
	@NEXT=$$(ls backend/migrations/*.up.sql 2>/dev/null | wc -l | tr -d ' '); \
	NEXT=$$(printf "%06d" $$((NEXT + 1))); \
	touch backend/migrations/$${NEXT}_$(name).up.sql; \
	touch backend/migrations/$${NEXT}_$(name).down.sql; \
	echo "Created: backend/migrations/$${NEXT}_$(name).up.sql"; \
	echo "Created: backend/migrations/$${NEXT}_$(name).down.sql"

# Seed database with demo data
seed:
	cd backend && go run ./cmd/seed

# ========================
# Mobile
# ========================

# Run mobile app
dev-mobile:
	cd mobile && npx expo start

# Install mobile dependencies
mobile-install:
	cd mobile && npm install

# ========================
# Web
# ========================

# Run web app locally
dev-web:
	cd web && npm run dev

# Install web dependencies
web-install:
	cd web && npm install

# Build web for production
build-web:
	cd web && npm run build

# Lint web
lint-web:
	cd web && npm run lint
