# 定义项目名称
PROJECT_NAME := ./desc/admin.api

GOCTL := goctl
GOFMT := gofmt

# 设置go模块环境变量
export GO111MODULE := on

# 格式化API定义文件
format:
	$(GOCTL) api format --dir ./desc

# 生成Go API代码
gen-go-api:
	$(GOCTL) api go --api $(PROJECT_NAME) --dir ./ --style=go_zero --home ../.go-zero-template
	mkdir -p swagger
	$(GOCTL) api plugin -plugin goctl-swagger="swagger -filename admin.json" -api $(PROJECT_NAME) -dir swagger


build:
	mkdir -p bin
	rm -rf bin/*
	go build -o bin/admin admin.go

lint:
	golangci-lint run --timeout=10m
# 默认目标
all: format gen-go-api

api: format gen-go-api

## fmt: Format Go code
fmt:
	$(GOFMT) -w -s .
	@echo "Go code formatted"

run:
	go run admin.go -f etc/admin-api.yaml

dev:
	go run admin.go -f etc/admin-api.yaml

.PHONY: format gen-go-api lint fmt run dev