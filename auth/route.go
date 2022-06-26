package auth

import (
	"context"
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
	m.Route("/callback", func(r chi.Router) {
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
	fmt.Println(payload.Claims)
}
