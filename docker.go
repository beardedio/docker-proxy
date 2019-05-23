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
				vDefaultPort := 80
				vPubPort := 0
				vPriPort := 0

				// Check for default port 80
				vPriPort = convertPrivatePortToPublic(container.Ports, vDefaultPort)

				// Override default port with VIRTUAL_PORT
				cJSON1, _ := cli.ContainerInspect(context.Background(), container.ID)
				for _, conEnv1 := range cJSON1.Config.Env {
					splits := strings.Split(conEnv1, "=")
					key := splits[0]
					val := splits[1]

					// Check for overrides to port
					if strings.HasPrefix(key, "VIRTUAL_PORT") {
						vDefaultPort, _ := strconv.Atoi(val)
						vPriPort = convertPrivatePortToPublic(container.Ports, vDefaultPort)
					}
				}

				// Add default host (container name) and port
				if vPriPort != 0 && vPriPort != myPort {
					AddSite(vHost, fmt.Sprintf("http://%s:%d", address, vPriPort))
				}


				cJSON2, _ := cli.ContainerInspect(context.Background(), container.ID)
				for _, conEnv2 := range cJSON2.Config.Env {
					splits := strings.Split(conEnv2, "=")
					key := splits[0]
					val := splits[1]

					// Check for overrides to hostname
					if strings.HasPrefix(key, "VIRTUAL_HOST") {

						if strings.Contains(val, ":") {
							splits = strings.Split(val, ":")
							vHost = splits[0]
							vPubPort, _ = strconv.Atoi(splits[1])
						} else {
							vHost = val
							vPubPort = vDefaultPort
						}

						vPriPort = convertPrivatePortToPublic(container.Ports, vPubPort)

						// Add extra hosts
						if vPriPort != 0 && vPriPort != myPort {
							AddSite(vHost, fmt.Sprintf("http://%s:%d", address, vPriPort))
						}
					}

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

// Convert private port to public port
func convertPrivatePortToPublic(PortList []types.Port, PriPort int) int {
	for _, port := range PortList {
		if int(port.PrivatePort) == PriPort {
			return int(port.PublicPort)
		}
	}

	return 0
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
