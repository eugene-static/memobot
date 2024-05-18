package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/eugene-static/memobot/internal/entities"
	"github.com/eugene-static/memobot/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(ctx context.Context, cfg *config.Database, l *slog.Logger) (*Storage, error) {
	err := os.MkdirAll(filepath.Dir(cfg.Path), 0750)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(cfg.Driver, cfg.Path)
	if err != nil {
		return nil, err
	}
	query := `
			PRAGMA foreign_keys = ON;
    		CREATE TABLE IF NOT EXISTS folders(
        	id VARCHAR(16) PRIMARY KEY UNIQUE,
    		user_id INT,
    		parent_id VARCHAR(16),
    		is_dir BOOL DEFAULT TRUE,
    		title TEXT,
    		FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE CASCADE
    		);
			CREATE TABLE IF NOT EXISTS notes(
			id VARCHAR(16) PRIMARY KEY UNIQUE,
			user_id INT,
			parent_id VARCHAR(16),
			is_dir BOOL DEFAULT FALSE,
			title TEXT,
			content TEXT DEFAULT '' NOT NULL,
			FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE CASCADE
			);
			CREATE TABLE IF NOT EXISTS media(
			id VARCHAR(16) PRIMARY KEY UNIQUE,
			parent_id VARCHAR(16),
			bytes BLOB,
			ext VARCHAR(4),
			FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE CASCADE
			);`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
		l:  l,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) GetListByParent(ctx context.Context, userID int64, parentID string) ([]*entities.List, error) {
	list := make([]*entities.List, 0)
	query := `
			SELECT * FROM (
				    SELECT id, title, is_dir
			        FROM folders
			        WHERE user_id = ? AND parent_id = ? 
			        UNION ALL SELECT id, title, is_dir
			              FROM notes
			              WHERE user_id = ? AND parent_id = ?)
			ORDER BY is_dir DESC, title ASC
			`
	rows, err := s.db.QueryContext(ctx, query, userID, parentID, userID, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		elem := &entities.List{}
		if err = rows.Scan(&elem.ID, &elem.Title, &elem.IsDir); err != nil {
			return nil, err
		}
		list = append(list, elem)
	}
	return list, nil
}
