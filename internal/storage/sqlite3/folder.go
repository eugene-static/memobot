package storage

import (
	"context"

	"github.com/eugene-static/memobot/internal/entities"
)

func (s *Storage) CheckRoot(ctx context.Context, userID int64, rootID string) (int, error) {
	var found int
	query := `SELECT EXISTS(SELECT * FROM folders where user_id = ? AND id = ?)`
	err := s.db.QueryRowContext(ctx, query, userID, rootID).Scan(&found)
	return found, err
}

func (s *Storage) AddRoot(ctx context.Context, folder *entities.Folder) error {
	query := `INSERT INTO folders(id, user_id, title) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, folder.ID, folder.UserID, folder.Title)
	return err
}

func (s *Storage) AddFolder(ctx context.Context, folder *entities.Folder) (string, error) {
	query := `INSERT INTO folders(id, user_id, parent_id, title) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, folder.ID, folder.UserID, folder.ParentID, folder.Title)
	return folder.ID, err
}

func (s *Storage) RenameFolder(ctx context.Context, id, title string) error {
	query := `UPDATE folders SET title = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, title, id)
	return err
}

func (s *Storage) DeleteFolder(ctx context.Context, id string) error {
	query := `DELETE FROM folders WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
