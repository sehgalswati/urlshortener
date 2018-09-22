# Usage:
DOCKER_ARGS=--build-arg HTTP_PROXY --build-arg HTTPS_PROXY --build-arg NO_PROXY --build-arg http_proxy --build-arg https_proxy --build-arg no_proxy --pull

release: clean build

build:
	docker build $(DOCKER_ARGS) -t urlshortener-builder .
	docker run urlshortener-builder | docker build $(DOCKER_ARGS) -t urlshortener -


run-urlshortener:
	docker-compose up
 
# remove previous images and containers
clean:
	docker rm -f urlshortener-builder 2> /dev/null || true
	docker rmi -f urlshortener-builder || true
	docker rm -f urlshortener 2> /dev/null || true
	docker rmi -f urlshortener || true
