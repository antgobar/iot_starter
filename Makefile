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

run-test-device-monolith:
	@echo "Running test device..."
	go fmt ./...
	go run cmd/device/main.go -m

fmt:
	@echo "Formatting code..."
	go fmt ./...