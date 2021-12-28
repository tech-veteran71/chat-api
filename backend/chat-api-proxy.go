package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ChatAPIProxy struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
	Token string
}

func (api *ChatAPIProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set Host header to the Chat-API hostname.
	r.Host = api.URL.Host

	// Let Chat-API know it's us.
	r.Header.Set("User-Agent", "https://github.com/andre-luiz-dos-santos/chat-api")

	// Add Chat-API token to the URL.
	uq := r.URL.Query()
	uq.Set("token", api.Token)
	r.URL.RawQuery = uq.Encode()

	// Contact Chat-API service.
	api.Proxy.ServeHTTP(w, r)
}
