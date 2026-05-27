backend:
	cd backend && go run ./cmd/server

frontend:
	cd frontend && npm run dev

test:
	cd backend && go test ./...

build:
	cd backend && go build ./cmd/server
	cd frontend && npm run build

compose-up:
	docker compose up --build
