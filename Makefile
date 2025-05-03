run-monolith:
	@echo "Running monolith..."
	pg_isready -q || (echo "PostgreSQL is not running. Please start it manually." && exit 1)
	go fmt ./...
	go run cmd/monolith/main.go

run-dashboard:
	@echo "Running dashboard..."
	pg_isready -q || (echo "PostgreSQL is not running. Please start it manually." && exit 1)
	go fmt ./...
	go run cmd/dashboard/main.go
