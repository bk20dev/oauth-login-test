package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"google.golang.org/api/idtoken"
	"net/http"
)

type Route struct {
	clientId     string
	repo         Repo
	sessionStore sessions.Store
}

func NewRoute(clientId string, repo Repo, sessionStore sessions.Store) *Route {
	return &Route{
		clientId:     clientId,
		repo:         repo,
		sessionStore: sessionStore,
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
		user, err := rt.repo.GetUser(entry.UserId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = rt.login(w, r, user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		user := &User{Name: payload.Claims["name"].(string)}
		id, err := rt.repo.CreateUser(user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := rt.repo.CreateOAuthEntry(GoogleProvider, payload.Subject, id); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		user, err = rt.repo.GetUser(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := rt.login(w, r, user); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (rt *Route) login(w http.ResponseWriter, r *http.Request, user *User) error {
	session, err := rt.sessionStore.New(r, "sid")
	if err != nil {
		return err
	}
	session.Values["user_id"] = user.Id
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}
