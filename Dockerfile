# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装基本工具
RUN apk add --no-cache gcc musl-dev

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -o shopee_tool ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装基本工具和时区数据
RUN apk add --no-cache ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D -g '' appuser

# 创建必要的目录
RUN mkdir -p /app/configs /app/logs
RUN chown -R appuser:appuser /app

# 切换到工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/shopee_tool .
COPY --from=builder /app/configs/dev.conf ./configs/

# 使用非 root 用户运行
USER appuser

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./shopee_tool"] 