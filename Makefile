APP_NAME := memento
MODULE := github.com/arrow2nd/memento
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
BUILD_TAGS := -tags prod
BUILD_FLAGS := -ldflags="-H=windowsgui -s -w -X $(MODULE)/app.appVersion=$(VERSION)"
DIST_DIR := dist
RELEASE_DIR := $(DIST_DIR)/$(APP_NAME)_$(VERSION)

.PHONY: build
build:
	@echo "Generating Windows resources..."
	go install github.com/tc-hib/go-winres@latest
	go generate

	@echo "Building $(APP_NAME) $(VERSION)..."
	@mkdir -p $(RELEASE_DIR)
	go build $(BUILD_TAGS) $(BUILD_FLAGS) -o "$(RELEASE_DIR)/$(APP_NAME).exe"

.PHONY: release
release: build
	@echo "Creating ZIP archive..."
	cd $(DIST_DIR) && zip -r $(APP_NAME)_$(VERSION).zip $(APP_NAME)_$(VERSION)
	@echo "Build completed: $(DIST_DIR)/$(APP_NAME)_$(VERSION).zip"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)
