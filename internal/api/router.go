package api

import (
	"fmt"
	"net/http"

	"github.com/arslan-atajykov/kanban/internal/auth"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()

	// Public endpoints
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Post("/register", auth.RegisterHand(db))
	r.Post("/login", auth.LoginHandler(db))

	// Protected endpoints (require JWT)
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware)

		// Пример защищённого эндпоинта /me
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(auth.UserIDKey).(int64)
			w.Write([]byte("Your user ID is: " + fmt.Sprint(userID)))
		})
	})

	return r
}
