package main

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// ContainerWatch checks for new containers and if they exist, add the sites and it's endpoints
func ContainerWatch(containerized bool, myPort int) {
	address := "127.0.0.1"
	if containerized {
		address = containerizedIP()
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.22"))
	if err != nil {
		log.Error(err.Error())
		return
	}

	for {
		if containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{}); err == nil {
			for _, container := range containers {
				// Use container name as default host
				vHost := container.Names[0][1:]
				vPort := 0
				vPriPort := 0

				cJSON, _ := cli.ContainerInspect(context.Background(), container.ID)
				for _, conEnv := range cJSON.Config.Env {
					splits := strings.Split(conEnv, "=")
					key := splits[0]
					val := splits[1]

					// Check for overrides to hostname
					if key == "VIRTUAL_HOST" {
						vHost = val
					}

					// Check for overrides to port
					if key == "VIRTUAL_PORT" {
						vPriPort, _ = strconv.Atoi(val)
					}
				}

				// Check for correct port
				for _, port := range container.Ports {
					if int(port.PrivatePort) == vPriPort || port.PrivatePort == 80 {
						vPort = int(port.PublicPort)
						break
					}
				}

				if vPort != 0 && vPort != myPort {
					AddSite(vHost, fmt.Sprintf("http://%s:%d", address, vPort))
				}
			}
		} else {
			log.Error("Unable to connect to docker")
			log.Error(err.Error())
		}

		// Every 5 seconds, check for new containers
		<-time.After(5 * time.Second)
	}
}

// containerizedIP returns a string with the ip address of the docker host
func containerizedIP() string {
	cmd := exec.Command("sh", "-c", "/sbin/ip route|awk '/default/ { print $3 }'")

	if output, err := cmd.Output(); err == nil {
		log.WithFields(log.Fields{
			"IP": strings.TrimSpace(string(output)),
		}).Info("Auto detecting docker host IP Address")

		return fmt.Sprintf("%s", strings.TrimSpace(string(output)))
	}

	return "127.0.0.1"
}
