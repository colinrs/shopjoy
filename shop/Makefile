# 定义项目名称
PROJECT_NAME := ./desc/shop.api

GOCTL := goctl

# 设置go模块环境变量
export GO111MODULE := on

# 格式化API定义文件
format:
	$(GOCTL) api format --dir ./desc

# 生成Go API代码
gen-go-api:
	$(GOCTL) api go --api $(PROJECT_NAME) --dir ./ --style=go_zero --home ../.go-zero-template
	mkdir -p swagger
	$(GOCTL) api plugin -plugin goctl-swagger="swagger -filename shop.json" -api $(PROJECT_NAME) -dir swagger


build:
	mkdir -p bin
	rm -rf bin/*
	go build -o bin/shop shop.go

lint:
	golangci-lint run --timeout=10m
# 默认目标
all: format gen-go-api

api: format gen-go-api

.PHONY: format gen-go-api lint