# 设置变量
GO := go
DIST_DIR := bin
BINARY_NAME := ${DIST_DIR}/sycode
TEST_DIRS := ./cmd ./sync
FORMAT_DIRS := ./cmd ./sync


OS := windows
ARCH := amd64

# 设置变量
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

PLATFORMS := linux/amd64 darwin/amd64 windows/amd64

# 默认目标
all: build_windows build_linux build_darwin build_darwin_arm64 compress_windows compress_linux compress_darwin compress_darwin_arm64


build_windows:
	@echo "Building the binary of Windows..."
	CGO_ENABLED=0  GOOS=windows GOARCH=amd64 $(GO) build -a -ldflags="-s -w" -gcflags="all=-l -B" -o $(BINARY_NAME)-windows-amd64.exe

compress_windows:
	@echo "compress the binary of windows"
	upx --best --lzma $(BINARY_NAME)-windows-amd64.exe



build_linux:
	@echo "Building the binary of Linux..."
	GO_ENABLED=0  GOOS=linux  GOARCH=amd64 $(GO) build -ldflags="-s -w" -gcflags="all=-l -B" -o $(BINARY_NAME)-linux-amd64

compress_linux:
	@echo "compress the binary of linux"
	upx --best --lzma $(BINARY_NAME)-linux-amd64



##
build_darwin:
	@echo "Building the binary of darwin..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="-s -w" -gcflags="all=-l -B" -o $(BINARY_NAME)-darwin-amd64


compress_darwin:
	@echo "compress the binary of darwin"
	upx --best --lzma --force-macos $(BINARY_NAME)-darwin-amd64

build_darwin_arm64:
	@echo "Building the binary of darwin..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="-s -w" -gcflags="all=-l -B" -o $(BINARY_NAME)-darwin-arm64


compress_darwin_arm64:
	@echo "compress the binary of darwin arm64"
	upx --best --lzma --force-macos $(BINARY_NAME)-darwin-arm64


# 运行测试
test:
	@echo "Running tests..."
	$(GO) test -v $(TEST_DIRS)

# 格式化代码
fmt:
	@echo "Formatting code..."
	$(GO) fmt $(FORMAT_DIRS)

# 安装依赖
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy
	$(GO) mod download


# 清理生成的文件
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf ./dist
	rm -rf ./coverage.out
	rm -rf ./coverprofile.out



