APP_NAME := memento
MODULE := github.com/arrow2nd/memento
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
BUILD_TAGS := -tags prod
BUILD_FLAGS := -ldflags="-H=windowsgui -s -w -X $(MODULE)/app.appVersion=$(VERSION)"
DIST_DIR := dist

.PHONY: build
build:
	@echo "Generating Windows resources..."
	go install github.com/tc-hib/go-winres@latest
	go generate
	@echo "Generating Windows resources..."
	go install github.com/tc-hib/go-winres@latest
	go generate
	@echo "Building $(APP_NAME) $(VERSION)..."
	@mkdir -p $(DIST_DIR)
	go build $(BUILD_TAGS) $(BUILD_FLAGS) -o "$(DIST_DIR)/$(APP_NAME)_$(VERSION).exe"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)
