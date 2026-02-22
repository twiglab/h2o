FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

# 创建工作目录结构，使得 ../shared 路径可用
WORKDIR /build

# 复制 shared 模块
COPY shared ./shared

# 复制 admin-server 模块
COPY admin-server ./admin-server

# 切换到 admin-server 目录进行构建
WORKDIR /build/admin-server

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /admin-server ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /admin-server /app/admin-server
COPY admin-server/config.yaml /app/config.yaml

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["/app/admin-server", "-config", "/app/config.yaml"]
