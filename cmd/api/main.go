package main

import (
	"database/sql"
	"log"
	"net/http"

	httpx "example.com/goprac11-borisovda/internal/http"
	"example.com/goprac11-borisovda/internal/http/handlers"
	"example.com/goprac11-borisovda/internal/repo"
)

func main() {
	dsn := "postgres://user:password@localhost:5432/notes_db?sslmode=disable"

	repo, err := repo.NewNoteRepoPostgres(dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}
	defer repo.Close()

	if err := createTableIfNotExists(repo.DB); err != nil {
		log.Fatalf("Ошибка создания таблицы %v", err)
	}

	h := &handlers.Handler{Repo: repo}
	r := httpx.NewRouter(h)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTableIfNotExists(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)
	return err
}
