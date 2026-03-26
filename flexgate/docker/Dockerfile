# --- 第一阶段：编译阶段 ---
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 处理依赖（利用 Docker 缓存）
COPY flexgate/go.mod flexgate/go.sum ./
RUN go mod download

# 拷贝源码并编译
COPY flexgate/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o flexgate ./cmd/flexgate/main.go

# --- 第二阶段：运行阶段 ---
FROM alpine:latest

WORKDIR /app
# 只从编译阶段拷贝二进制文件和配置文件
COPY --from=builder /app/flexgate .
COPY --from=builder /app/configs ./configs

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./flexgate"]