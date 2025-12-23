.PHONY: all build vet clean test

BUILD_DIR ?= build
LDFLAGS :=	-s -w

all: build

fmt:
	go mod tidy && go fmt ./...

vet: fmt
	go vet ./...

build: vet test
	go build -o $(BUILD_DIR)/ -ldflags="$(LDFLAGS)"

test:
	go test ./...

clean:
	rm -r $(BUILD_DIR)
