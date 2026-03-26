# --- 第一阶段：编译 (Builder) ---
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 1. 拷贝该服务的依赖文件
# 因为 context 在根目录，所以路径是 services/user-service/...
COPY services/user-service/go.mod services/user-service/go.sum ./
RUN go mod download

# 2. 拷贝该服务的源代码
COPY services/user-service/ .

# 3. 编译成静态二进制文件
# CGO_ENABLED=0 确保在 alpine 这种轻量级镜像中能直接运行
RUN CGO_ENABLED=0 GOOS=linux go build -o service-bin main.go

# --- 第二阶段：运行 (Runner) ---
FROM alpine:latest

# 安装基础库（如 tzdata 用于处理时区）
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从编译阶段拷贝二进制文件
COPY --from=builder /app/service-bin .

# 暴露端口（user-service 是 9001）
EXPOSE 9001

# 启动服务
CMD ["./service-bin"]