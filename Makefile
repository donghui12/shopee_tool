.PHONY: build run clean docker-build docker-run docker-compose-up

# 构建应用
build:
	go build -o bin/shopee_tool ./cmd/main.go

# 运行应用
run: build
	./bin/shopee_tool

# 清理构建文件
clean:
	rm -rf bin/

# 构建 Docker 镜像
docker-build:
	docker build -t shopee_tool:latest .

# 运行 Docker 容器
docker-run:
	docker run -d --name shopee_tool -p 8080:8080 shopee_tool:latest

# 使用 docker-compose 启动所有服务
docker-compose-up:
	docker-compose up -d

# 停止并删除所有服务
docker-compose-down:
	docker-compose down

# 查看服务日志
logs:
	docker-compose logs -f 