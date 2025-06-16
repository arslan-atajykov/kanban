package main

import (
	"log"
	"net/http"

	"github.com/arslan-atajykov/kanban/internal/api"
	"github.com/arslan-atajykov/kanban/internal/db"
)

func main() {
	conn := db.Init()

	router := api.SetupRouter(conn)

	log.Println("Server is running on localhost : 8888")
	http.ListenAndServe(":8888", router)
}
