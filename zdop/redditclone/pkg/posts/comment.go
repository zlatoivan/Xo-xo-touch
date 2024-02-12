package posts

import "time"

type Comment struct {
	ID      string    `json:"id"`
	Author  Author    `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

type CommentsRepo interface {
	GetAll() ([]*Comment, error)
	GetByID(id string) (*Comment, error)
	Add(comment *Comment) (string, error)
	Delete(id string) error
}
