.PHONY: dev dev-backend dev-frontend build test test-backend test-frontend db-up db-down

db-up:
	docker compose up -d

db-down:
	docker compose down

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && pnpm dev

dev: db-up
	$(MAKE) dev-backend & $(MAKE) dev-frontend & wait

build:
	cd backend && go build -o ../bin/server ./cmd/server
	cd frontend && pnpm build

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && pnpm test --run

test: test-backend test-frontend
