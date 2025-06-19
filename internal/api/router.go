package api

import (
	"fmt"
	"net/http"

	"github.com/arslan-atajykov/kanban/internal/auth"
	"github.com/arslan-atajykov/kanban/internal/board"
	"github.com/arslan-atajykov/kanban/internal/task"
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
		r.Route("/boards", func(r chi.Router) {
			r.Post("/", board.CreateBoardHandler(db))
			r.Get("/", board.ListBoardsHandler(db))
			r.Put("/{id}", board.UpdateBoardHandler(db))
			r.Delete("/{id}", board.DeleteBoardHandler(db))
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Post("/", task.CreateTaskHandler(db))
			r.Get("/", task.ListTasksHandler(db))
			r.Get("/{taskID}", task.GetTaskByIDHandler(db))
			r.Put("/{taskID}", task.UpdateTaskHandler(db))
			r.Delete("/{taskID}", task.DeleteTaskHandler(db))
		})
	})

	return r
}
