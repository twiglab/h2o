FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 创建工作目录结构，使得 ../shared 路径可用
WORKDIR /build

# 复制 shared 模块
COPY shared ./shared

# 复制 electric-billing 模块
COPY electric-billing ./electric-billing

# 切换到 electric-billing 目录进行构建
WORKDIR /build/electric-billing

# 下载依赖
RUN go mod download

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /electric-billing ./cmd/server

# 运行镜像
FROM alpine:3.19

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

WORKDIR /app

# 创建数据目录
RUN mkdir -p /app/data/wal /app/data/deductions /app/data/readings

# 从构建镜像复制二进制文件
COPY --from=builder /electric-billing /app/electric-billing

# 复制配置文件
COPY electric-billing/config.yaml /app/config.yaml

# 设置文件权限
RUN chown -R appuser:appuser /app

# 切换到非root用户
USER appuser

# 暴露端口（如需HTTP API可添加）
# EXPOSE 8080

# 启动命令
CMD ["/app/electric-billing", "-config", "/app/config.yaml"]
