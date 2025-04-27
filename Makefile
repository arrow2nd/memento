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
	@if not exist $(subst /,\\,$(RELEASE_DIR)) mkdir $(subst /,\\,$(RELEASE_DIR))
	@go build $(BUILD_TAGS) $(BUILD_FLAGS) -o "$(RELEASE_DIR)/$(APP_NAME).exe"

.PHONY: release
release: build
	@echo "Creating ZIP archive..."
	@powershell -Command "if (Get-Command Compress-Archive -ErrorAction SilentlyContinue) { \
		Compress-Archive -Path '$(subst /,\\,$(RELEASE_DIR))' -DestinationPath '$(subst /,\\,$(DIST_DIR))/$(APP_NAME)_$(VERSION).zip' -Force; \
	} else { \
		echo 'PowerShell Compress-Archive not available, trying 7-Zip...'; \
		if (Get-Command 7z -ErrorAction SilentlyContinue) { \
			7z a -tzip '$(subst /,\\,$(DIST_DIR))/$(APP_NAME)_$(VERSION).zip' '$(subst /,\\,$(RELEASE_DIR))'; \
		} else { \
			echo 'ERROR: Neither Compress-Archive nor 7z are available. Please install one of them.'; \
			exit 1; \
		} \
	}"
	@echo "Build completed: $(DIST_DIR)/$(APP_NAME)_$(VERSION).zip"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@if exist $(subst /,\\,$(DIST_DIR)) rmdir /S /Q $(subst /,\\,$(DIST_DIR))
