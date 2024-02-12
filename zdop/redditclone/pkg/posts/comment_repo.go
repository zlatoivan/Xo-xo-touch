package posts

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoComment = errors.New("no comment found")
)

type CommentMemoryRepository struct {
	data  []Comment
	mutex *sync.Mutex
}

func NewCommentMemoryRepo() *CommentMemoryRepository {
	return &CommentMemoryRepository{
		data:  make([]Comment, 0, 10),
		mutex: &sync.Mutex{},
	}
}

func (repo *CommentMemoryRepository) GetAll() ([]*Comment, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	comments := make([]*Comment, 0, len(repo.data))
	for _, comment := range repo.data {
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (repo *CommentMemoryRepository) GetByID(id string) (*Comment, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	idxToGet := -1
	for idx, comment := range repo.data {
		if comment.ID == id {
			idxToGet = idx
			break
		}
	}

	if idxToGet == -1 {
		return nil, ErrNoComment
	}

	comment := repo.data[idxToGet]

	return &comment, nil
}

func (repo *CommentMemoryRepository) Add(comment *Comment) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	newID := uuid.New()
	comment.ID = newID.String()
	comment.Created = time.Now()
	repo.data = append(repo.data, *comment)

	return comment.ID, nil
}

func (repo *CommentMemoryRepository) Delete(id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	idxToDelete := -1
	for idx, comment := range repo.data {
		if comment.ID == id {
			idxToDelete = idx
			break
		}
	}

	if idxToDelete == -1 {
		return ErrNoComment
	}

	repo.data = append(repo.data[:idxToDelete], repo.data[idxToDelete+1:]...)

	return nil
}
