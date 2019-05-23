.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clear any built containers
	docker rm docker-proxy-local

.PHONY: build
build: ## Build the docker-proxy docker container
	docker build --pull --no-cache -t docker-proxy-local .

.PHONY: push
push: ## Push the docker container to docker hub
	docker tag docker-proxy-local beardedio/docker-proxy
	docker push beardedio/docker-proxy

.PHONY: run
run: ## Run docker-proxy localy from the repo
	- docker rm -f docker-proxy-local
	docker run --rm -p 80:80 -v /var/run/docker.sock:/var/run/docker.sock --name docker-proxy-local docker-proxy-local --containerized
