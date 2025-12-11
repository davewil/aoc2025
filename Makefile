.PHONY: test bench lint run

# Run all tests
test:
	go test ./...

# Run all benchmarks
bench:
	go test -bench=. -benchmem ./...

# Run a specific day (usage: make run day=10)
run:
	go run ./day$(day)

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run
