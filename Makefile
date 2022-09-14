.PHONY: run test start stop build deploy

run:
	swag init
	go run main.go

test:
	poetry install --no-dev

start:
	swag init
	docker compose build
	docker compose up -d

stop:
	docker compose down

build:
	swag init
	docker buildx create --use
	docker buildx build --platform linux/amd64,linux/arm64 --push --tag docker.io/avnes/test-persistence --file Dockerfile .


deploy:
	# az login
	# az acr login -n itoacr --expose-token
	# docker push itoacr.azurecr.io/ito/test-persistence
	docker push docker.io/avnes/test-persistence
