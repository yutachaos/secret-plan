
.PHONY: build
build:
	go build -o secret-plan ./cmd/secret-plan/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...
	gofumpt -w -s ./..

.PHONY: lint
lint:
	golangci-lint run