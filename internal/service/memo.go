package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/eugene-static/memobot/internal/entities"
	"github.com/eugene-static/memobot/pkg/random"
)

const root = "root"

type List interface {
	GetListByParent(context.Context, int64, string) ([]*entities.List, error)
}
type Folder interface {
	AddRoot(context.Context, *entities.Folder) error
	CheckRoot(context.Context, int64, string) (int, error)
	AddFolder(context.Context, *entities.Folder) (string, error)
	RenameFolder(context.Context, string, string) error
	DeleteFolder(context.Context, string) error
}

type Note interface {
	AddNote(context.Context, *entities.Note) (string, error)
	GetNoteContent(context.Context, int64, string) (string, error)
	Update(context.Context, string, string) error
	RenameNote(context.Context, string, string) error
	DeleteNote(context.Context, string) error
}

type Storage interface {
	List
	Folder
	Note
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddRoot(ctx context.Context, userID int64, dirID string, title string) error {
	found, err := s.storage.CheckRoot(ctx, userID, dirID)
	if err != nil {
		return err
	}
	if found == 0 {
		return s.storage.AddRoot(ctx,
			&entities.Folder{
				ID:     dirID,
				UserID: userID,
				Title:  title,
			})
	}
	return nil
}

func (s *Service) GetList(ctx context.Context, userID int64, id string) ([]*entities.List, error) {
	return s.storage.GetListByParent(ctx, userID, id)
}

func (s *Service) Get(ctx context.Context, userID int64, id string) (string, error) {
	return s.storage.GetNoteContent(ctx, userID, id)
}

func (s *Service) Add(ctx context.Context, userID int64, parentID string, title string, isDir bool) (string, error) {
	if isDir {
		return s.storage.AddFolder(ctx,
			&entities.Folder{
				ID:       fmt.Sprintf("1%s", random.String(15)),
				UserID:   userID,
				ParentID: parentID,
				Title:    title,
			})
	}
	return s.storage.AddNote(ctx, &entities.Note{
		ID:       fmt.Sprintf("2%s", random.String(15)),
		UserID:   userID,
		ParentID: parentID,
		Title:    title,
	})
}

func (s *Service) UpdateContent(ctx context.Context, id, content string) error {
	return s.storage.Update(ctx, id, content)
}

func (s *Service) Rename(ctx context.Context, id, title string) error {
	if strings.HasPrefix(id, "1") {
		return s.storage.RenameFolder(ctx, id, title)
	}
	return s.storage.RenameNote(ctx, id, title)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.HasPrefix(id, "1") {
		return s.storage.DeleteFolder(ctx, id)
	}
	return s.storage.DeleteNote(ctx, id)
}
