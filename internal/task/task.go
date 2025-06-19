package task

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/arslan-atajykov/kanban/internal/auth"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	ID          int64     `db:"id" json:"id"`
	BoardID     int64     `db:"board_id" json:"board_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type CreateTaskRequest struct {
	BoardID     int64  `json:"board_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey).(int64)

		var req CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		log.Printf("User %d is creating a task for board %d", userID, req.BoardID)
		res, err := db.Exec(`INSERT INTO tasks (board_id, title, description, status, created_at)
			VALUES (?, ?, ?, 'todo', ?)`, req.BoardID, req.Title, req.Description, time.Now())
		if err != nil {
			http.Error(w, "failed to create task", http.StatusInternalServerError)
			return
		}

		id, _ := res.LastInsertId()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": id,
		})
	}
}

func GetTaskByIDHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, _ := strconv.Atoi(chi.URLParam(r, "taskID"))

		var task Task
		if err := db.Get(&task, "SELECT * FROM tasks WHERE id = ?", taskID); err != nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(task)
	}
}

func ListTasksHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		boardID, _ := strconv.Atoi(r.URL.Query().Get("board_id"))

		var tasks []Task
		if err := db.Select(&tasks, "SELECT * FROM tasks WHERE board_id = ?", boardID); err != nil {
			http.Error(w, "could not get tasks", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(tasks)
	}
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func UpdateTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskID")

		var req UpdateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ?`,
			req.Title, req.Description, req.Status, taskID)
		if err != nil {
			http.Error(w, "failed to update task", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("task updated"))
	}
}

func DeleteTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskID")

		_, err := db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
		if err != nil {
			http.Error(w, "failed to delete task", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("task deleted"))
	}
}
