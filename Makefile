.PHONY: dev-db dev-backend dev-mobile migrate-up migrate-down test-backend dev stop infra

# Start all infrastructure (PostgreSQL + Redis)
infra:
	docker compose up -d

# Start PostgreSQL only
dev-db:
	docker compose up postgres -d

# Run backend locally
dev-backend:
	cd backend && go run ./cmd/api

# Run mobile app
dev-mobile:
	cd mobile && npx expo start

# Database migrations
migrate-up:
	cd backend && go run ./cmd/migrate up

migrate-down:
	cd backend && go run ./cmd/migrate down

# Run backend tests
test-backend:
	cd backend && go test ./...

# Stop all services
stop:
	docker compose down
