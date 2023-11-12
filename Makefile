.PHONY:

# ==============================================================================
# local

run:
	go run cmd\cli\main.go -path=$(path) -cpucount=4 -countWorker=10 -maxprocessurldurationmsec=1000 -maxtotaldurationsecond=600

test:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Linters

run-linter:
	echo "Starting linters"
	golangci-lint run ./...


