package repo

import (
	"sync"
	"time"

	"example.com/goprac11-borisovda/internal/core"
)

type NoteRepoMem struct {
	mu    sync.Mutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	n.ID = r.next
	r.notes[n.ID] = &n
	return n.ID, nil
}

func (r *NoteRepoMem) GetAll() ([]core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := []core.Note{}
	for _, v := range r.notes {
		res = append(res, *v)
	}
	return res, nil
}

func (r *NoteRepoMem) Get(id int64) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	v, ok := r.notes[id]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (r *NoteRepoMem) Update(id int64, upd core.Note) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notes[id]
	if !ok {
		return nil, nil
	}
	now := time.Now()
	n.Title = upd.Title
	n.Content = upd.Content
	n.UpdatedAt = &now

	return n, nil
}

func (r *NoteRepoMem) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.notes, id)
	return nil
}
