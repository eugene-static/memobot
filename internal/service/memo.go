package service

import (
	"fmt"
	"strings"

	"github.com/eugene-static/memobot/internal/entities"
	"github.com/eugene-static/memobot/pkg/random"
)

const root = "root"

type List interface {
	GetListByParent(int64, string) ([]*entities.List, error)
}
type Folder interface {
	AddRoot(*entities.Folder) error
	CheckRoot(int64, string) (int, error)
	AddFolder(*entities.Folder) (string, error)
	RenameFolder(string, string) error
	DeleteFolder(string) error
}

type Note interface {
	AddNote(*entities.Note) (string, error)
	GetNoteContent(int64, string) (string, error)
	Update(string, string) error
	RenameNote(string, string) error
	DeleteNote(string) error
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

func (s *Service) AddRoot(userID int64, dirID string, title string) error {
	found, err := s.storage.CheckRoot(userID, dirID)
	if err != nil {
		return err
	}
	if found == 0 {
		return s.storage.AddRoot(&entities.Folder{
			ID:     dirID,
			UserID: userID,
			Title:  title,
		})
	}
	return nil
}

func (s *Service) GetList(userID int64, id string) ([]*entities.List, error) {
	return s.storage.GetListByParent(userID, id)
}

func (s *Service) Get(userID int64, id string) (string, error) {
	return s.storage.GetNoteContent(userID, id)
}

func (s *Service) Add(userID int64, parentID string, title string, isDir bool) (string, error) {
	if isDir {
		return s.storage.AddFolder(&entities.Folder{
			ID:       fmt.Sprintf("1%s", random.String(15)),
			UserID:   userID,
			ParentID: parentID,
			Title:    title,
		})
	}
	return s.storage.AddNote(&entities.Note{
		ID:       fmt.Sprintf("2%s", random.String(15)),
		UserID:   userID,
		ParentID: parentID,
		Title:    title,
	})
}

func (s *Service) UpdateContent(id, content string) error {
	return s.storage.Update(id, content)
}

func (s *Service) Rename(id, title string) error {
	if strings.HasPrefix(id, "1") {
		return s.storage.RenameFolder(id, title)
	}
	return s.storage.RenameNote(id, title)
}

func (s *Service) Delete(id string) error {
	if strings.HasPrefix(id, "1") {
		return s.storage.DeleteFolder(id)
	}
	return s.storage.DeleteNote(id)
}
