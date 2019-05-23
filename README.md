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

## Additional Domains Names

To add additional domain names to a container, add this env var `VIRTUAL_HOST` to the list of environment variables to your container. To add more then one domain add more env vars like `VIRTUAL_HOST_1` &/or `VIRTUAL_HOST_blah`

You can also specify a port along with the domain name by using the format `VIRTUAL_HOST=mycoolhost.test:443`

## Override default port

To override the default port looked for inside a container add the env var `VIRTUAL_PORT`

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
