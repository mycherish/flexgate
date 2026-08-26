.PHONY: dev-docker stop-docker logs

dev-docker:
	docker compose -f deployments/docker-compose.yml up --build -d

stop-docker:
	docker compose -f deployments/docker-compose.yml down

logs:
	docker compose -f deployments/docker-compose.yml logs -f