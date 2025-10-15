# Messages Worker SDK Makefile

.PHONY: test build clean examples install

run-pipeline:
	@echo ******RUNNING BUILD******
	go build
	@echo ******MAKING SURE LINT IS CORRECT******
#	go get -u golang.org/x/lint/golint
#	golint -set_exit_status api/... iplocation/... utils/... shieldio/... ./
	@echo ******STARTING TESTS******
	go test -gcflags="all=-l -N" -v ./...
	@echo ******DONE******

# Run tests
test:
	go test -v ./...

# Build the SDK
build:
	go build ./...

# Clean build artifacts
clean:
	go clean ./...

# Run examples (requires messages-worker service running)
examples:
	cd examples && go run main.go

# Install dependencies
install:
	go mod tidy
	go mod download

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate go.sum
deps:
	go mod tidy
