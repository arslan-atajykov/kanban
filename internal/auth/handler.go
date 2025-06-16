package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHand(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		if len(req.Password) < 6 {
			http.Error(w, "Password is too short", http.StatusNotAcceptable)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`INSERT INTO users(email, password_hash, created_at)Values(?,?,?)`, req.Email, string(hash), time.Now())
		if err != nil {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered"))

	}

}
