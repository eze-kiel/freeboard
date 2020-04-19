package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eze-kiel/freeboard/handlers"
)

func main() {
	router := handlers.HandleFunc()
	timeoutRouter := http.TimeoutHandler(router, time.Second*3, "Timeout!")
	log.Fatal(http.ListenAndServe(":8080", timeoutRouter))
}
