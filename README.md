# Docker-Proxy

Docker-Proxy is a simple reverse proxy that routes web traffic to running docker containers to host ports. It is designed to run in dev/ci environments.

* Response timeouts
* Automagically maps new containers
* Ideal for dev environments
* Custom configured dns names and ports

## Installation

`docker run -d -p 80:80 -v /var/run/docker.sock:/var/run/docker.sock --restart=always --name docker-proxy beardedio/docker-proxy --containerized`

## Container Setup

By default Docker-Proxy with automatically route traffic based on the containers name when using port 80 inside the container.

## Custom DNS Name/Port

To override the domain or port for routing to a container, add these env vars `VIRTUAL_HOST` &/or `VIRTUAL_PORT` to the list of environment variables to your container.


Here is a docker-compose.yml example
```yaml
version: '3'
services:
  apache:
    image: beardedio/php-apache:php7
    ports:
     - 12345
    environment:
      - VIRTUAL_HOST=cic.test
      - VIRTUAL_PORT=12345
      - server_env=dev
```
