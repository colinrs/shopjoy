# 定义项目名称
PROJECT_NAME := goleaf.api

GOCTL := goctl

# 设置go模块环境变量
export GO111MODULE := on

# 格式化API定义文件

# 生成Go API代码
gen-go-api:
	cd ./admin && make api
	cd ./shop  && make api

build:
	cd ./admin && make build
	cd ./shop  && make build

lint:
	golangci-lint run --timeout=10m
# 默认目标
all: gen-go-api

api: gen-go-api

.PHONY:  gen-go-api lint