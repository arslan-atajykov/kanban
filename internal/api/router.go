package api

import (
	"net/http"

	"github.com/arslan-atajykov/kanban/internal/auth"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Post("/register", auth.RegisterHand(db))
	r.Post("/login", auth.LoginHandler(db))
	return r
}
