package main

import (
	"log"
	"net/http"

	httpx "example.com/goprac11-borisovda/internal/http"
	"example.com/goprac11-borisovda/internal/http/handlers"
	"example.com/goprac11-borisovda/internal/repo"
)

func main() {
	repo := repo.NewNoteRepoMem()
	h := &handlers.Handler{Repo: repo}
	r := httpx.NewRouter(h)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
