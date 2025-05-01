APP_NAME=labrador
BUILD_DIR=bin
VERSION=v0.2.1

.PHONY: all build clean package

all: build package

build:
	mkdir -p $(BUILD_DIR)
	# macOS (amd64)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-macos-amd64 ./cmd/labrador
	# macOS (arm64)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-macos-arm64 ./cmd/labrador
	# Linux (amd64)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/labrador
	# Windows (amd64)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/labrador

package:
	# For Mac (amd64)
	cd $(BUILD_DIR) && mkdir macos-amd64 && cp $(APP_NAME)-macos-amd64 macos-amd64/$(APP_NAME) && tar -czvf $(APP_NAME)-macos-amd64.tar.gz -C macos-amd64 $(APP_NAME) && rm -r macos-amd64
	# For Mac (arm64)
	cd $(BUILD_DIR) && mkdir macos-arm64 && cp $(APP_NAME)-macos-arm64 macos-arm64/$(APP_NAME) && tar -czvf $(APP_NAME)-macos-arm64.tar.gz -C macos-arm64 $(APP_NAME) && rm -r macos-arm64
	# For Linux
	cd $(BUILD_DIR) && mkdir linux-amd64 && cp $(APP_NAME)-linux-amd64 linux-amd64/$(APP_NAME) && tar -czvf $(APP_NAME)-linux-amd64.tar.gz -C linux-amd64 $(APP_NAME) && rm -r linux-amd64
	# For Windows
	cd $(BUILD_DIR) && mkdir windows-amd64 && cp $(APP_NAME)-windows-amd64.exe windows-amd64/$(APP_NAME).exe && zip -j $(APP_NAME)-windows-amd64.zip windows-amd64/$(APP_NAME).exe && rm -r windows-amd64


clean:
	rm -rf $(BUILD_DIR)
