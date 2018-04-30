package main

import (
	"fmt"
	"./services"
	"golang.org/x/net/http2"
	"log"
	"net/http"
)

func main() {
	// redirect every http request to https
	go http.ListenAndServe(services.Config.Address, http.HandlerFunc(services.Redirect))
	mux := http.NewServeMux()

	// serving style and js files(only if no Bootstrap CDM):
	statics := http.FileServer(http.Dir(services.Config.Static))
	mux.Handle("/cssjs/", http.StripPrefix("/cssjs/", statics))

	// defined in server_routes.go
	mux.HandleFunc("/", services.Index)
	mux.HandleFunc("/err", services.Err)
	mux.HandleFunc("/about", services.About)

	// defined in handler_functions.go
	mux.HandleFunc("/signup_account", services.SignupAccount)
	mux.HandleFunc("/authenticate", services.Authenticate)
	mux.HandleFunc("/logout", services.Logout)
	mux.HandleFunc("/delete_account", services.DelAccount)

	// javascript calls:
	mux.HandleFunc("/bruse", services.Bruse)

	server := http.Server{
		Addr:    services.Config.AddressSSL,
		Handler: mux,
	}
	fmt.Println(services.Config.Address + " and " + services.Config.AddressSSL)
	http2.ConfigureServer(&server, &http2.Server{})

	log.Fatal(server.ListenAndServeTLS(services.Config.Encryptcl1, services.Config.Encryptcl2))

	log.Fatal(server.ListenAndServeTLS("/etc/letsencrypt/live/datapenetration.de/fullchain.pem", "/etc/letsencrypt/live/datapenetration.de/privkey.pem"))

}
