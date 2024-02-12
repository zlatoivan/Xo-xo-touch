package posts

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoPost = errors.New("no post found")
)

type PostMemoryRepository struct {
	data  map[string]Post
	mutex *sync.Mutex
}

func NewPostMemoryRepo() *PostMemoryRepository {
	return &PostMemoryRepository{
		data:  map[string]Post{},
		mutex: &sync.Mutex{},
	}
}

func (repo *PostMemoryRepository) GetAll() ([]*Post, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	posts := make([]*Post, 0, len(repo.data))
	for _, post := range repo.data {
		posts = append(posts, &post)
	}
	return posts, nil
}

func (repo *PostMemoryRepository) GetByID(id string) (*Post, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if post, ok := repo.data[id]; ok {
		return &post, nil
	}
	return nil, ErrNoPost
}

func (repo *PostMemoryRepository) GetByAuthor(authorUsername string) ([]*Post, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	posts := []*Post{}
	for _, post := range repo.data {
		if post.Author.Username == authorUsername {
			posts = append(posts, &post)
		}
	}
	return posts, nil
}

func (repo *PostMemoryRepository) GetByCategory(category string) ([]*Post, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	posts := []*Post{}
	for _, post := range repo.data {
		if post.Category == category {
			posts = append(posts, &post)
		}
	}
	return posts, nil
}

func (repo *PostMemoryRepository) Add(post *Post) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	newID := uuid.New()
	post.ID = newID.String()
	post.Created = time.Now()
	post.Comments = NewCommentMemoryRepo()
	post.Votes = NewVoteMemoryRepo()
	repo.data[post.ID] = *post

	return post.ID, nil
}

func (repo *PostMemoryRepository) UpdateViews(id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if post, ok := repo.data[id]; ok {
		post.Views++
		return nil
	}
	return ErrNoPost
}

func (repo *PostMemoryRepository) Delete(id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if _, ok := repo.data[id]; !ok {
		return ErrNoPost
	}
	delete(repo.data, id)

	return nil
}
