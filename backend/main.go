package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	flagConfig string
	flagProxy  = ""
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.StringVar(&flagConfig, "config", "config.json", "Path to the configuration file")
	flag.StringVar(&flagProxy, "proxy", "", "Overrides proxy in config (for testing)")
	flag.Parse()

	// Configuration.
	cf, err := ReadConfig(flagConfig)
	if err != nil {
		log.Fatal(err)
	}
	if flagProxy != "" {
		cf.Proxy = flagProxy
	}

	// Database.
	db := NewDatabase(cf.Database.Driver, cf.Database.DSN)
	err = db.Open()
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	err = db.Create()
	if err != nil {
		log.Fatalf("Cannot create database: %v", err)
	}

	// Authentication.
	auth := &Auth{
		Database: db,
	}

	// Chat-API.
	u, err := url.Parse(cf.ChatAPI.URL)
	if err != nil {
		log.Fatal(err)
	}
	chatAPI := &ChatAPI{
		URL:   u,
		Token: cf.ChatAPI.Token,
	}
	err = chatAPI.GetStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Connect Chat-API to local database.
	wadb := NewChatAPIDB(chatAPI, db)
	err = wadb.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	wadbHTTP := NewChatAPIHTTP(wadb)

	// HTTP muxes.
	mux := http.NewServeMux()
	apiMux := http.NewServeMux()

	// HTTP routes.
	mux.HandleFunc("/api/login", auth.HTTPLogin)
	mux.Handle("/api/", auth.Protect(http.StripPrefix("/api", apiMux)))
	apiMux.HandleFunc("/logout", auth.HTTPLogout)
	apiMux.HandleFunc("/messages/all", wadbHTTP.Messages)
	apiMux.HandleFunc("/messages/chat_id", wadbHTTP.MessagesByChatID)
	apiMux.HandleFunc("/chats/all", wadbHTTP.Chats)
	apiMux.HandleFunc("/chat/info", wadbHTTP.GetUserChatInfo)
	apiMux.HandleFunc("/chat/info/read", wadbHTTP.SetUserChatAsRead)

	// Webhook.
	if cf.ChatAPI.Webhook != "" {
		// Create random URL.
		u, err = url.Parse(cf.ChatAPI.Webhook)
		if err != nil {
			log.Fatal(err)
		}
		u.Path += "/" + generateToken()
		// Handle HTTP requests.
		mux.HandleFunc(u.Path, wadbHTTP.Webhook)
		// Inform Chat-API our webhook URL.
		err = chatAPI.SetWebhook(ctx, u.String())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Proxy to the Chat-API service.
	u, err = url.Parse(cf.ChatAPI.URL)
	if err != nil {
		log.Fatal(err)
	}
	chatAPIProxy := &ChatAPIProxy{
		URL:   u,
		Proxy: httputil.NewSingleHostReverseProxy(u),
		Token: cf.ChatAPI.Token,
	}
	mux.Handle("/chat-api/", auth.Protect(http.StripPrefix("/chat-api", chatAPIProxy)))

	// Proxy to the SvelteKit dev server.
	if cf.Proxy != "" {
		u, err = url.Parse(cf.Proxy)
		if err != nil {
			log.Fatal(err)
		}
		devProxy := httputil.NewSingleHostReverseProxy(u)
		mux.Handle("/", devProxy)
	} else {
		// If no proxy is set, files in /static/ will be served.
		staticFS, err := fs.Sub(staticFiles, "static")
		if err != nil {
			log.Fatal(err)
		}
		mux.Handle("/", http.FileServer(http.FS(staticFS)))
	}

	http.ListenAndServe(":8080", mux)
}
