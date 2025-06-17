package board

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

type CreateBoardRequest struct {
	Title string `json:"title"`
}

type UpdateBoardRequest struct {
	Title string `json:"title"`
}

func CreateBoardHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth.UserIDKey).(int64)
		log.Println("Creating board for user:", userID)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var req CreateBoardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		res, err := db.Exec("INSERT INTO boards(title,owner_id,created_at) VALUES(?,?,?)", req.Title, userID, time.Now())
		if err != nil {
			http.Error(w, "Failed to create a board", http.StatusInternalServerError)
			return
		}

		boardID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to get board ID", http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"id":    boardID,
			"title": req.Title,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	}
}

func ListBoardsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth.UserIDKey).(int64)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var boards []struct {
			ID    int64  `db:"id" json:"id"`
			Title string `db:"title" json:"title"`
		}
		err := db.Select(&boards, `SELECT id, title FROM boards WHERE owner_id =?`, userID)
		if err != nil {
			http.Error(w, "error fetching boards", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(boards)
	}
}

func UpdateBoardHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth.UserIDKey).(int64)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		boardIDStr := chi.URLParam(r, "id")
		boardID, err := strconv.ParseInt(boardIDStr, 10, 64)

		if err != nil {
			http.Error(w, "invalid board id", http.StatusBadRequest)
			return
		}

		var req UpdateBoardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		res, err := db.Exec(`UPDATE boards SET title = ? WHERE owner_id = ? AND id = ?`, req.Title, userID, boardID)
		if err != nil {
			http.Error(w, "Failed to update a board", http.StatusInternalServerError)
			return
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "board not found or not owned by you", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

}
func DeleteBoardHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth.UserIDKey).(int64)
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		boardIDStr := chi.URLParam(r, "id")
		boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Failed to get board id ", http.StatusBadRequest)
			return
		}

		res, err := db.Exec(`DELETE FROM boards WHERE id = ? AND owner_id = ?`, boardID, userID)
		if err != nil {
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "board not found or not owned by you", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
