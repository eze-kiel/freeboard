package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/eze-kiel/freeboard/handlers"
)

func main() {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlers.HandleFunc(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info("Server is starting, wish me luck boys")
	log.Println(srv.ListenAndServe())
}
