package main

import (
	"errors"
	"net/http/httputil"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// Endpoint struct containing everything needed for a new endpoint
type Endpoint struct {
	Active     bool
	Address    *url.URL
	Proxy      *httputil.ReverseProxy
	Registered string
	Available  time.Time
}

// NewEndpoint creates new endpoints to forward traffic to
func NewEndpoint(base, address string) (*Endpoint, error) {
	parsedAddress, err := url.Parse(address)
	if err != nil {
		log.WithFields(log.Fields{"url": parsedAddress}).Error("Problem parsing URL")
		return nil, errors.New("Problem parsing URL " + address)
	}

	e := &Endpoint{
		Address:    parsedAddress,
		Proxy:      httputil.NewSingleHostReverseProxy(parsedAddress),
		Active:     true,
		Available:  time.Now(),
		Registered: base,
	}

	return e, nil
}
