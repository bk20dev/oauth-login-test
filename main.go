package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	server := http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start a server: %v", err)
	}
}
