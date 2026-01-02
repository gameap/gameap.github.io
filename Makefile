.PHONY: build run stop clean

IMAGE_NAME := registry.gameap.io/docs
CONTAINER_NAME := gameap-docs
PORT := 8080

build:
	docker buildx build -t $(IMAGE_NAME) --load --no-cache .

run: build
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):80 $(IMAGE_NAME)
	@echo "Server running at http://localhost:$(PORT)"

stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

clean: stop
	docker rmi $(IMAGE_NAME) || true

restart: stop run
