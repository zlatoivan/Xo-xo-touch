package posts

import "time"

type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type Post struct {
	ID       string
	Title    string
	Views    uint32
	Type     string
	Text     string
	Created  time.Time
	URL      string
	Votes    VotesRepo
	Category string
	Author   Author
	Comments CommentsRepo
}

type PostsRepo interface {
	GetAll() ([]*Post, error)
	GetByID(id string) (*Post, error)
	GetByAuthor(authorID string) ([]*Post, error)
	GetByCategory(category string) ([]*Post, error)
	Add(post *Post) (string, error)
	UpdateViews(id string) error
	Delete(id string) error
}
