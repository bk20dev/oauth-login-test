package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"rabbit/auth"
)

func main() {
	pgConnStr := os.Getenv("PG_CONNECTION_STRING")
	db, err := openDB("postgres", pgConnStr)
	if err != nil {
		log.Fatal(err)
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	sessionStore := sessions.NewCookieStore([]byte(sessionSecret))

	router := chi.NewRouter()
	setupRoutes(router, db, sessionStore)

	server := http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start a server: %v", err)
	}
}

func openDB(driver, dataSource string) (*sql.DB, error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func setupRoutes(router *chi.Mux, db *sql.DB, sessionStore sessions.Store) {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	repo := auth.NewPostgresRepo(db)
	authRoute := auth.NewRoute(clientId, repo, sessionStore)
	authRoute.Register(router)
}
