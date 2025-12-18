.PHONY: all clean

BUILD_DIR ?= build

clean:
	go clean -i ./... && rm -r $(BUILD_DIR)