.PHONY: all build clean debug test
BIN_FILE=gimg
all: build
build:
	@go build -o "${BIN_FILE}" cmd/gimg/main.go

clean:
	@go clean
	@rm "${BIN_FILE}"

test:
	@go test -v ./...

debug:
	@go run cmd/gimg/main.go