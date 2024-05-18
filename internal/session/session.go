package session

import (
	"fmt"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	folder = "📁"
	note   = "📄"
)

type Dir struct {
	ID    string
	Title string
}

type User struct {
	ID            int64
	Username      string
	Action        string
	LastMessageID int
	Dir           []*Dir
	Timer         *time.Timer
}

type Manager struct {
	storage map[int64]*User
	mu      sync.Mutex
}

func New() *Manager {
	return &Manager{
		storage: make(map[int64]*User),
		mu:      sync.Mutex{},
	}
}

func (m *Manager) timer(userID int64) *time.Timer {
	return time.AfterFunc(time.Minute, func() {
		m.mu.Lock()
		delete(m.storage, userID)
		m.mu.Unlock()
	})
}

func (m *Manager) updateTime(userID int64) {
	if !m.storage[userID].Timer.Stop() {
		<-m.storage[userID].Timer.C
	}
	m.storage[userID].Timer = m.timer(userID)
}

func (m *Manager) GetUser(userID int64, name string) *User {
	c, ok := m.storage[userID]
	if !ok {
		return nil
	}
	return c
}

func (m *Manager) AddUser(userID int64, name string) *User {
	m.storage[userID] = &User{
		ID:       userID,
		Username: name,
		Timer:    m.timer(userID),
	}
	m.storage[userID].Root(name)
	return m.storage[userID]
}

func (u *User) Root(name string) *User {
	u.Dir = []*Dir{{ID: fmt.Sprintf("0%s", name), Title: name}}
	return u
}

func (u *User) Down(id, title string) string {
	u.Dir = append(u.Dir, &Dir{ID: id, Title: title})
	return u.Path()
}

func (u *User) Up() string {
	if len(u.Dir) < 2 {
		return ""
	}
	u.Dir = u.Dir[:len(u.Dir)-1]
	return u.Path()
}

func (u *User) CurrentDir() *Dir {
	if len(u.Dir) == 0 {
		return nil
	}
	return u.Dir[len(u.Dir)-1]
}

func (u *User) ParentDir() *Dir {
	switch len(u.Dir) {
	case 0:
		return nil
	case 1:
		return u.Dir[0]
	}
	return u.Dir[len(u.Dir)-2]
}

func (u *User) Path() string {
	d := make([]string, len(u.Dir))
	symbol := folder
	for i := range u.Dir {
		d[i] = u.Dir[i].Title
		if strings.HasPrefix(u.Dir[i].ID, "2") {
			symbol = note
		}
	}
	return symbol + path.Join(d...)
}

func (u *User) NewAction(action string) {
	u.Action = action
}
