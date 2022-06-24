package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"rabbit/auth"
)

func setupRoutes(router *chi.Mux) {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")

	authRoute := auth.NewRoute(clientId)
	authRoute.Register(router)
}

func main() {
	router := chi.NewRouter()

	setupRoutes(router)

	server := http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start a server: %v", err)
	}
}
