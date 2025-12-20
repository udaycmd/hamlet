.PHONY: all build vet clean

BUILD_DIR ?= build
LDFLAGS :=	-s -w

all: build

fmt:
	go mod tidy && go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o $(BUILD_DIR)/ -ldflags="$(LDFLAGS)"

clean:
	rm -r $(BUILD_DIR)
