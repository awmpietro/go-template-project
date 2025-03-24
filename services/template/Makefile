run-dev:
	@echo "Running development server..."
	air

swag:
	@echo "Generating Swagger docs..."
	@swag init --parseDependency --parseInternal -g cmd/main.go -o docs

tests:
	@echo "Running tests..."
	@go test ./... 

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -v -cover

test-bench:
	@echo "Running benchmark tests..."
	@go test ./... -v -bench .

test-cover-bench:
	@echo "Running benchmark tests with coverage..."
	@go test ./... -v -cover -bench .

test-cover-profile:
	@echo "Running coverage profiling..."
	@go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out