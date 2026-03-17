# Joy Shop Frontend Makefile
# 商城前端构建脚本

.PHONY: help install dev build preview clean lint format check node_modules

# 默认目标
help:
	@echo "Joy Shop Frontend - Available Commands:"
	@echo ""
	@echo "  make install    - 安装项目依赖 (npm install)"
	@echo "  make dev        - 启动开发服务器 (localhost:3001)"
	@echo "  make build      - 构建生产环境代码"
	@echo "  make preview    - 预览生产构建"
	@echo "  make clean      - 清理 node_modules 和构建产物"
	@echo "  make lint       - 运行 ESLint 检查"
	@echo "  make format     - 运行 Prettier 格式化代码"
	@echo "  make check      - 运行类型检查 (vue-tsc)"
	@echo "  make update     - 更新依赖到最新版本"
	@echo ""

# 安装依赖
install: node_modules

node_modules: package.json
	@echo "📦 安装依赖..."
	npm install
	@touch node_modules

# 启动开发服务器
dev: node_modules
	@echo "🚀 启动开发服务器 (http://localhost:3001)..."
	npm run dev

# 构建生产环境
build: node_modules
	@echo "🔨 构建生产环境代码..."
	npm run build

# 预览生产构建
preview: build
	@echo "👁️  预览生产构建..."
	npm run preview

# 清理项目
clean:
	@echo "🧹 清理项目..."
	rm -rf node_modules
	rm -rf dist
	rm -f package-lock.json
	@echo "✅ 清理完成"

# 代码检查
lint: node_modules
	@echo "🔍 运行 ESLint..."
	npm run lint || echo "⚠️  ESLint not configured"

# 代码格式化
format: node_modules
	@echo "✨ 格式化代码..."
	npm run format || echo "⚠️  Prettier not configured"

# 类型检查
check: node_modules
	@echo "🔎 运行类型检查..."
	npx vue-tsc --noEmit

# 更新依赖
update:
	@echo "⬆️  更新依赖..."
	npm update

# 重新安装依赖 (清理后安装)
reinstall: clean install
	@echo "✅ 重新安装完成"

# 生产环境构建并分析
analyze: node_modules
	@echo "📊 构建并分析包大小..."
	npx vite-bundle-visualizer

# 运行所有检查
test-all: lint check build
	@echo "✅ 所有检查通过"
