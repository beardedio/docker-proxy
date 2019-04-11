package main

import (
	"flag"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/kcmerrill/shutdown.go"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	//log.SetLevel(log.InfoLevel)
}

func main() {
	// Setup some command line arguments
	port := flag.Int("port", 80, "The port in which the proxy will listen on")
	containerized := flag.Bool("containerized", false, "Is fetch-proxy running in a container?")
	timeout := flag.Int("response-timeout", 10, "The response timeout for the proxy")
	defaultEndpoint := flag.String("default", "__default", "The default endpoint fetch-proxy uses when requested endpoing isn't found")

	flag.Parse()

	// Set a global timeout
	http.DefaultClient.Timeout = time.Duration(*timeout) * time.Second

	// Start our proxy on the specified port
	go DProxyStart(*port, *defaultEndpoint)

	go ContainerWatch(*containerized, *port)

	// No need to shutdown the application _UNLESS_ we catch it
	shutdown.WaitFor(syscall.SIGINT, syscall.SIGTERM)
	log.Info("Shutting down ... ")
}
