# ShopJoy 多阶段构建
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建 Admin 服务
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/admin ./admin/admin.go

# 构建 Shop 服务
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/shop ./shop/shop.go

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从 builder 复制二进制文件
COPY --from=builder /app/bin/admin /usr/local/bin/admin
COPY --from=builder /app/bin/shop /usr/local/bin/shop

# 复制配置文件
COPY --from=builder /app/admin/etc /etc/admin
COPY --from=builder /app/shop/etc /etc/shop

EXPOSE 8888 8889

CMD ["admin", "-f", "/etc/admin/admin-api.yaml"]
