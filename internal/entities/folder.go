package entities

type Folder struct {
	ID       string
	UserID   int64
	ParentID string
	Title    string
}
