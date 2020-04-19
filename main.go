package main

import (
	"log"
	"net/http"

	"github.com/eze-kiel/freeboard/handlers"
)

func main() {
	router := handlers.HandleFunc()
	log.Fatal(http.ListenAndServe(":8080", router))
}
