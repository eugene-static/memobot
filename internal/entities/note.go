package entities

type Note struct {
	ID       string
	UserID   int64
	ParentID string
	Title    string
	Content  string
}
