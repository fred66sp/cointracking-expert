BINARY := dist/cointracking-mcp.exe

.PHONY: build test vet fmt clean smoketest

build:
	go build -o $(BINARY) ./cmd/cointracking-mcp

test:
	go vet ./...
	go test ./...

fmt:
	gofmt -w .

smoketest: build
	go run ./cmd/smoketest

clean:
	rm -rf dist cache cache_smoketest
