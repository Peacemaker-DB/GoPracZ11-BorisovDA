package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/goprac11-borisovda/internal/core"
	httpx "example.com/goprac11-borisovda/internal/http"
	"example.com/goprac11-borisovda/internal/http/handlers"
	"example.com/goprac11-borisovda/internal/repo"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (string, func()) {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("notes_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("Не удалось запустить контейнер PostgreSQL: %v", err)
	}

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://test:test@%s:%s/notes_test?sslmode=disable", host, port.Port())

	db, _ := sql.Open("postgres", dsn)
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)

	cleanup := func() {
		db.Exec("TRUNCATE TABLE notes RESTART IDENTITY")
		pgContainer.Terminate(ctx)
	}

	return dsn, cleanup
}

func createTestServer(t *testing.T, dsn string) *httptest.Server {
	t.Helper()

	repo, _ := repo.NewNoteRepoPostgres(dsn)
	db := repo.DB
	db.Exec("TRUNCATE TABLE notes RESTART IDENTITY")

	h := &handlers.Handler{Repo: repo}
	r := httpx.NewRouter(h)

	return httptest.NewServer(r)
}

func testCase(name string, t *testing.T, testFunc func(*testing.T, string)) {
	fmt.Printf("\n=== %s ===\n", name)
	dsn, cleanup := setupTestDB(t)
	defer cleanup()
	server := createTestServer(t, dsn)
	defer server.Close()
	testFunc(t, server.URL)
	fmt.Printf("✅ Успешно\n")
}

func TestAllCases(t *testing.T) {

	testCase("Создание заметки (POST /api/v1/notes) → 201 Created", t, func(t *testing.T, url string) {
		noteData := map[string]string{
			"title":   "Тестовая заметка",
			"content": "содержание",
		}
		jsonData, _ := json.Marshal(noteData)

		resp, _ := http.Post(url+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Ожидался статус 201, получен: %d", resp.StatusCode)
		}

		var createdNote core.Note
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &createdNote)

		fmt.Printf("Создана заметка %d\n", createdNote.ID)
		fmt.Printf("Статус: %d Created\n", resp.StatusCode)
	})

	testCase("Получение заметки (GET /api/v1/notes/{id}) → 200 OK", t, func(t *testing.T, url string) {
		noteData := map[string]string{
			"title":   "Заметка для получения",
			"content": "cодержание",
		}
		jsonData, _ := json.Marshal(noteData)
		resp, _ := http.Post(url+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))

		var createdNote core.Note
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &createdNote)
		resp.Body.Close()

		resp, _ = http.Get(fmt.Sprintf("%s/api/v1/notes/%d", url, createdNote.ID))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
		}

		fmt.Printf("Получена заметка %d\n", createdNote.ID)
		fmt.Printf("Статус %d OK\n", resp.StatusCode)
	})

	testCase("Получение всех заметок (GET /api/v1/notes) → 200 OK", t, func(t *testing.T, url string) {
		notes := []map[string]string{
			{"title": "Заметка 1", "content": "Содержание 1"},
			{"title": "Заметка 2", "content": "Содержание 2"},
		}

		for i, note := range notes {
			jsonData, _ := json.Marshal(note)
			resp, _ := http.Post(url+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
			resp.Body.Close()
			fmt.Printf("Создана заметка %d\n", i+1)
		}

		resp, _ := http.Get(url + "/api/v1/notes")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен: %d", resp.StatusCode)
		}

		var allNotes []core.Note
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &allNotes)

		fmt.Printf("Получено заметок: %d\n", len(allNotes))
		fmt.Printf("Статус: %d OK\n", resp.StatusCode)
	})

	testCase("Обновление заметки (PATCH /api/v1/notes/{id}) → 200 OK", t, func(t *testing.T, url string) {
		noteData := map[string]string{
			"title":   "Исходная заметка",
			"content": "Исходное содержание",
		}
		jsonData, _ := json.Marshal(noteData)
		resp, _ := http.Post(url+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))

		var createdNote core.Note
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &createdNote)
		resp.Body.Close()

		updateData := map[string]string{
			"title":   "Обновлённая заметка",
			"content": "Обновлённое содержание",
		}
		updateJson, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PATCH",
			fmt.Sprintf("%s/api/v1/notes/%d", url, createdNote.ID),
			bytes.NewBuffer(updateJson))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, _ = client.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен: %d", resp.StatusCode)
		}

		fmt.Printf("Обновлена заметка ID: %d\n", createdNote.ID)
		fmt.Printf("Статус: %d OK\n", resp.StatusCode)
	})

	testCase("Удаление заметки (DELETE /api/v1/notes/{id}) → 204 No Content", t, func(t *testing.T, url string) {
		noteData := map[string]string{
			"title":   "Заметка для удаления",
			"content": "Содержание для удаления",
		}
		jsonData, _ := json.Marshal(noteData)
		resp, _ := http.Post(url+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))

		var createdNote core.Note
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &createdNote)
		resp.Body.Close()

		req, _ := http.NewRequest("DELETE",
			fmt.Sprintf("%s/api/v1/notes/%d", url, createdNote.ID), nil)

		client := &http.Client{}
		resp, _ = client.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Ожидался статус 204, получен: %d", resp.StatusCode)
		}

		fmt.Printf("Удалена заметка ID: %d\n", createdNote.ID)
		fmt.Printf("Статус: %d No Content\n", resp.StatusCode)
	})

	testCase("Несуществующая заметка → 404 Not Found", t, func(t *testing.T, url string) {
		resp, _ := http.Get(url + "/api/v1/notes/999")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Ожидался статус 404, получен: %d", resp.StatusCode)
		}

		fmt.Printf("Запрос несуществующей заметки ID: 999\n")
		fmt.Printf("Статус: %d Not Found\n", resp.StatusCode)
	})

}
