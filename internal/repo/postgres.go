package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/goprac11-borisovda/internal/core"
	_ "github.com/lib/pq"
)

type NoteRepoPostgres struct {
	DB *sql.DB
}

func NewNoteRepoPostgres(dsn string) (*NoteRepoPostgres, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("нет подключения к БД %w", err)
	}

	return &NoteRepoPostgres{DB: db}, nil
}

func (r *NoteRepoPostgres) Create(ctx context.Context, n *core.Note) error {
	query := `
		INSERT INTO notes (title, content, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`

	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = &now

	err := r.DB.QueryRowContext(ctx, query,
		n.Title, n.Content, n.CreatedAt, n.UpdatedAt,
	).Scan(&n.ID)

	return err
}

func (r *NoteRepoPostgres) GetAll(ctx context.Context) ([]core.Note, error) {
	query := `
		SELECT id, title, content, created_at, updated_at 
		FROM notes 
		ORDER BY created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []core.Note
	for rows.Next() {
		var n core.Note
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, nil
}

func (r *NoteRepoPostgres) Get(ctx context.Context, id int64) (*core.Note, error) {
	query := `
		SELECT id, title, content, created_at, updated_at 
		FROM notes 
		WHERE id = $1
	`

	var n core.Note
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("заметка %d не найдена", id)
	}

	return &n, err
}

func (r *NoteRepoPostgres) Update(ctx context.Context, id int64, upd core.Note) (*core.Note, error) {
	query := `
		UPDATE notes 
		SET title = $1, content = $2, updated_at = $3 
		WHERE id = $4 
		RETURNING id, title, content, created_at, updated_at
	`

	now := time.Now()
	var n core.Note
	err := r.DB.QueryRowContext(ctx, query,
		upd.Title, upd.Content, now, id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("заметка %d не найдена", id)
	}

	return &n, err
}

func (r *NoteRepoPostgres) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notes WHERE id = $1`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("заметка %d не найдена", id)
	}

	return nil
}

func (r *NoteRepoPostgres) Close() error {
	return r.DB.Close()
}
