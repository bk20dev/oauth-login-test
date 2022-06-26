package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"google.golang.org/api/idtoken"
	"net/http"
)

type Route struct {
	clientId string
	repo     Repo
}

func NewRoute(clientId string, repo Repo) *Route {
	return &Route{
		clientId: clientId,
		repo:     repo,
	}
}

func (rt *Route) Register(m *chi.Mux) {
	m.Route("/m", func(r chi.Router) {
		r.Post("/google", rt.googleCallback)
	})
}

func (rt *Route) googleCallback(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("credential")
	payload, err := idtoken.Validate(context.Background(), token, rt.clientId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	entry, err := rt.repo.GetOAuthEntry(GoogleProvider, payload.Subject)
	if entry != nil {
		// TODO: Login
		return
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		user := User{Name: payload.Claims["name"].(string)}
		id, err := rt.repo.CreateUser(&user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := rt.repo.CreateOAuthEntry(GoogleProvider, payload.Subject, id); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Println("Created an account", id)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
