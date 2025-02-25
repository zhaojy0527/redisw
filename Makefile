# 定义变量
BINARY_NAME=redisw
GO=go

# 默认目标
.PHONY: all
all: build

# 构建项目
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME)

# 运行测试
.PHONY: test
test:
	$(GO) test -v ./...

# 清理构建文件
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# 安装到系统
.PHONY: install
install:
	cp $(BINARY_NAME) /usr/local/bin/

# 帮助信息
.PHONY: help
help:
	@echo "可用的 make 命令："
	@echo "make build    - 构建项目"
	@echo "make test     - 运行测试"
	@echo "make clean    - 清理构建文件"
	@echo "make install  - 安装到系统"