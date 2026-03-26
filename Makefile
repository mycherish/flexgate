.PHONY: dev-docker stop-docker logs

# 一键启动容器化开发环境
dev-docker:
	cd deployments && docker-compose up --build -d

# 停止所有服务
stop-docker:
	cd deployments && docker-compose down

# 查看所有服务日志
logs:
	cd deployments && docker-compose logs -f