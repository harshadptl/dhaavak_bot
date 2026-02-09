BINARY   := dhaavak
BUILD_DIR := bin
GO       := go
GOFLAGS  := -trimpath

.PHONY: build run clean test tidy

build:
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/dhaavak

run: build
	./$(BUILD_DIR)/$(BINARY) --config dhaavak.yaml

test:
	$(GO) test ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BUILD_DIR)
