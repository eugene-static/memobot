package storage

import (
	"context"
	"errors"

	"github.com/eugene-static/memobot/internal/entities"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) AddNote(ctx context.Context, note *entities.Note) (string, error) {
	query := `INSERT INTO notes(id, user_id, parent_id, title) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, note.ID, note.UserID, note.ParentID, note.Title)
	if errors.Is(err, sqlite3.ErrConstraintUnique) {
		err = nil
	}
	return note.ID, err
}

func (s *Storage) GetNoteContent(ctx context.Context, userID int64, id string) (string, error) {
	query := `SELECT content FROM notes WHERE user_id = ? AND id = ?`
	var row string
	err := s.db.QueryRowContext(ctx, query, userID, id).Scan(&row)
	if err != nil {
		return "", err
	}
	return row, nil
}

func (s *Storage) Update(ctx context.Context, id string, content string) error {
	query := `UPDATE notes SET content = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, content, id)
	return err
}

func (s *Storage) RenameNote(ctx context.Context, id, title string) error {
	query := `UPDATE notes SET title = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, title, id)
	return err
}

func (s *Storage) DeleteNote(ctx context.Context, id string) error {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
