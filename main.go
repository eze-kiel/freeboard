package main

import (
	"github.com/caddyserver/certmagic"
	log "github.com/sirupsen/logrus"

	"github.com/eze-kiel/freeboard/handlers"
)

func main() {
	// srv := &http.Server{
	// 	Addr:         ":8080",
	// 	Handler:      handlers.HandleFunc(),
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// read and agree to your CA's legal documents
	certmagic.DefaultACME.Agreed = true

	// provide an email address
	certmagic.DefaultACME.Email = "hugoblanc@fastmail.com"

	log.Info("Server is starting, wish me luck boys")
	// err := certmagic.HTTPS([]string{"example.com", "www.example.com"}, mux)

	certmagic.HTTPS([]string{"freeboard.tech", "www.freeboard.tech"}, handlers.HandleFunc())
	// log.Println(srv.ListenAndServe())
}
