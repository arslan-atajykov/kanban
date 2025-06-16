package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Init() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "kanban.db")
	if err != nil {
		log.Fatalln("Failed to connect to database: ", err)
	}
	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver : %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}
	err = m.Up()

	if err != nil && err.Error() != "no change" {
		log.Fatalf("Migration failed :%v ", err)
	}
	log.Println("Migration applied successfully")
	return db
}
