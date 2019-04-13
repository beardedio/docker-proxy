package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/kcmerrill/shutdown.go"
	log "github.com/sirupsen/logrus"
)

// Store all of our endpoints
var endpoints map[string]*Endpoint
var endpointkeys sort.StringSlice

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// passThrough takes in traffic on specific port and passes it through to the appropriate endpoint
func passThrough(w http.ResponseWriter, r *http.Request, defaultEndpoint string) {
	// remove www.
	if strings.HasPrefix(r.Host, "www.") {
		r.Host = strings.Replace(r.Host, "www.", "", 1)
	}

	endpoint := siteKey(r.Host, defaultEndpoint)

	// One quick sanity check before sending it on it's way
	if _, exists := endpoints[endpoint]; exists {
		log.WithFields(
			log.Fields{
				"Host":       r.Host,
				"From":       r.RemoteAddr,
				"To":         endpoints[endpoint].Address,
				"RequestURI": r.RequestURI,
				"Forwarded":  endpoint,
			}).Info("New Request")

		endpoints[endpoint].Proxy.ServeHTTP(w, r)
	} else {
		log.WithFields(
			log.Fields{
				"Request":   r.Host,
				"From":      r.RemoteAddr,
				"Forwarded": endpoint,
			}).Info("Bad Request")

		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Error 502 - Bad Gateway"))
	}
}

// DProxyStart creates and starts the proxy
func DProxyStart(httpPort int, defaultEndpoint string) {
	log.WithFields(
		log.Fields{
			"port": httpPort,
		}).Info("Starting fetch proxy")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		passThrough(w, r, defaultEndpoint)
	})

	// Not secured, so lets just start a simple webserver
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
		log.Fatal(err.Error())
		shutdown.Now()
	}

}

// AddSite adds a new website to the proxy to be forwarded
func AddSite(base, address string) error {
	// Check if endpoint already exists
	for _, item := range endpoints {
		if item.Registered == base && item.Address.String() == address {
			return nil
		}
	}

	// Construct the key so that you can sort by url base and time added
	urlbase := base

	// Remove any thing after the _ from the url
	if strings.Contains(urlbase, "_") {
		urlbase = urlbase[0:strings.Index(urlbase, "_")]
	}

	key := urlbase

	// Add new endpoint
	ep, err := NewEndpoint(base, address)
	if err == nil {
		// If it doesn't exist ...
		log.WithFields(log.Fields{
			"url":        address,
			"registered": base,
			"urlbase":    urlbase,
		}).Info("Registered endpoint")

		endpoints[key] = ep
		if !stringInSlice(key, endpointkeys) {
			endpointkeys = append(endpointkeys, key)
		}

		sort.Sort(sort.Reverse(endpointkeys))

		return nil
	}
	return err
}

// Site key determines the endpoint to use based on the host
func siteKey(host, defaultEndpoint string) string {
	registered := ""
	// Grab the first key in the list that matches
	for _, key := range endpointkeys {
		b := endpoints[key].Registered

		// Allow for multiple containers with the same url
		if strings.Contains(b, "_") {
			b = b[0:strings.Index(b, "_")]
		}

		if strings.HasPrefix(defaultEndpoint, b) && endpoints[key].Active {
			defaultEndpoint = key
		}

		if strings.HasPrefix(host, b) && endpoints[key].Active {
			registered = key
			break
		}
	}

	if registered == "" {
		return defaultEndpoint
	}

	return registered
}

// init our maps
func init() {
	endpoints = make(map[string]*Endpoint)
}
