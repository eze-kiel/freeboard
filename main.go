package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/eze-kiel/freeboard/handlers"
)

func main() {
	router := handlers.HandleFunc()
	timeoutRouter := http.TimeoutHandler(router, time.Second*3, "Timeout!")
	log.Info("Server is starting, wish me luck boys")
	log.Fatal(http.ListenAndServe(":8080", timeoutRouter))
}
