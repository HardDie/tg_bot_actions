.PHONY: docker-build
docker-build:
	docker build -t bot -f deployments/Dockerfile .

.PHONY: docker-run
docker-run:
	docker run --rm -d --env-file .env --name bot bot
