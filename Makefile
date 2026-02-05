.PHONY: all build vet clean test

BUILD_DIR ?= build
LDFLAGS :=	-s -w -extldflags="-static"

all: build

fmt:
	go mod tidy && go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o $(BUILD_DIR)/hamlet -ldflags="$(LDFLAGS)" ./src/

debug: vet
	go build -o $(BUILD_DIR)/hamlet.debug -race ./src/

test:
	-go test -v -race ./...

clean:
	rm -r $(BUILD_DIR)
