# 定义版本信息变量
GIT_TAG := $(shell git describe --tags --abbrev=0)
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S%z)
GO_VERSION := $(shell go version | awk '{print $$3}')

# Go 源文件和二进制文件名
SRC := main.go
BINARY := harbor-img

# LDFLAGS 设置编译时注入的变量
LDFLAGS := -X 'version.GIT_TAG=$(GIT_TAG)' \
           -X 'version.GIT_COMMIT=$(GIT_COMMIT)' \
           -X 'version.GIT_BRANCH=$(GIT_BRANCH)' \
           -X 'version.BUILD_TIME=$(BUILD_TIME)' \
           -X 'version.GO_VERSION=$(GO_VERSION)'

# 默认目标
all: build

# 构建目标
build:
    go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(SRC)

# 清理目标
clean:
    rm -f $(BINARY)

# 打印版本信息目标
version:
    @echo "GIT_TAG:    $(GIT_TAG)"
    @echo "GIT_COMMIT: $(GIT_COMMIT)"
    @echo "GIT_BRANCH: $(GIT_BRANCH)"
    @echo "BUILD_TIME: $(BUILD_TIME)"
    @echo "GO_VERSION: $(GO_VERSION)"

.PHONY: all build clean version
