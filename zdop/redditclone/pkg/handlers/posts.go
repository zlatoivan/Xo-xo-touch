package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"redditclone/pkg/posts"
	"redditclone/pkg/session"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

const JSONContent string = "application/json"

type PostHandler struct {
	PostRepo posts.PostsRepo
}

type PostResponse struct {
	ID               string           `json:"id"`
	Title            string           `json:"title"`
	Views            uint32           `json:"views"`
	Type             string           `json:"type"`
	Text             string           `json:"text,omitempty"`
	Created          time.Time        `json:"created"`
	URL              string           `json:"url,omitempty"`
	Votes            []*posts.Vote    `json:"votes"`
	Score            int              `json:"score"`
	UpvotePercentage int              `json:"upvotePercentage"`
	Category         string           `json:"category"`
	Author           posts.Author     `json:"author"`
	Comments         []*posts.Comment `json:"comments"`
}

func NewPostResponse(post *posts.Post) (*PostResponse, error) {
	votes, err := post.Votes.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get votes")
	}

	score, err := post.Votes.GetScore()
	if err != nil {
		return nil, fmt.Errorf("failed to get score")
	}

	upvotePercentage, err := post.Votes.GetUpvotePercentage()
	if err != nil {
		return nil, fmt.Errorf("failed to get upvote percentage")
	}

	comments, err := post.Comments.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get comments")
	}

	resp := &PostResponse{
		ID:               post.ID,
		Title:            post.Title,
		Views:            post.Views,
		Type:             post.Type,
		Text:             post.Text,
		Created:          post.Created,
		URL:              post.URL,
		Votes:            votes,
		Score:            score,
		UpvotePercentage: upvotePercentage,
		Category:         post.Category,
		Author:           post.Author,
		Comments:         comments,
	}

	return resp, nil
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	postsData, err := h.PostRepo.GetAll()
	if err != nil {
		http.Error(w, "failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", JSONContent)
	w.WriteHeader(http.StatusOK)

	resp := make([]PostResponse, len(postsData))
	for i, post := range postsData {
		respCurrent, errNewResp := NewPostResponse(post)
		if errNewResp != nil {
			http.Error(w, errNewResp.Error(), http.StatusInternalServerError)
			return
		}
		resp[i] = *respCurrent
	}
	sort.Slice(resp, func(i, j int) bool { return resp[i].Score > resp[j].Score })

	bytes, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to write data: %s", err.Error())
	}
}

type PostCreateRequest struct {
	Category string `json:"category"`
	Text     string `json:"text,omitempty"`
	URL      string `json:"url,omitempty"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

func (h *PostHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != JSONContent {
		http.Error(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postData := &PostCreateRequest{}
	err = json.Unmarshal(body, postData)
	if err != nil {
		http.Error(w, "failed to unpack payload", http.StatusBadRequest)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		if err == session.ErrNoAuth {
			http.Error(w, "no authorization", http.StatusBadRequest)
		} else {
			http.Error(w, "unknown session", http.StatusInternalServerError)
		}
		return
	}

	newPostID, err := h.PostRepo.Add(&posts.Post{
		Category: postData.Category,
		Text:     postData.Text,
		URL:      postData.URL,
		Title:    postData.Title,
		Type:     postData.Type,
		Author: posts.Author{
			ID:       sess.UserID,
			Username: sess.Username,
		},
	})
	if err != nil {
		http.Error(w, "failed to add post", http.StatusInternalServerError)
		return
	}

	newPost, err := h.PostRepo.GetByID(newPostID)
	if err != nil {
		http.Error(w, "failed to get post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", JSONContent)
	w.WriteHeader(http.StatusCreated)

	resp, err := NewPostResponse(newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to write data: %s", err.Error())
	}
}

func (h *PostHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryName := vars["category_name"]

	postsData, err := h.PostRepo.GetByCategory(categoryName)
	if err != nil {
		http.Error(w, "failed to get post by category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", JSONContent)
	w.WriteHeader(http.StatusOK)

	resp := make([]PostResponse, len(postsData))
	for i, post := range postsData {
		respCurrent, errNewResp := NewPostResponse(post)
		if errNewResp != nil {
			http.Error(w, errNewResp.Error(), http.StatusInternalServerError)
			return
		}
		resp[i] = *respCurrent
	}
	sort.Slice(resp, func(i, j int) bool { return resp[i].Score > resp[j].Score })

	bytes, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to write data: %s", err.Error())
	}
}

func (h *PostHandler) ListUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["user_login"]

	postsData, err := h.PostRepo.GetByAuthor(username)
	if err != nil {
		http.Error(w, "failed to get posts by author", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", JSONContent)
	w.WriteHeader(http.StatusOK)

	resp := make([]PostResponse, len(postsData))
	for i, post := range postsData {
		respCurrent, errNewResp := NewPostResponse(post)
		if errNewResp != nil {
			http.Error(w, errNewResp.Error(), http.StatusInternalServerError)
			return
		}
		resp[i] = *respCurrent
	}
	sort.Slice(resp, func(i, j int) bool { return resp[i].Score > resp[j].Score })

	bytes, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to write data: %s", err.Error())
	}
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		if err == posts.ErrNoPost {
			http.Error(w, "there is no post with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to get post by id", http.StatusInternalServerError)
		}
		return
	}

	err = h.PostRepo.UpdateViews(postID)
	if err != nil {
		if err == posts.ErrNoPost {
			http.Error(w, "there is no post with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to update post views", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", JSONContent)
	w.WriteHeader(http.StatusOK)

	resp, err := NewPostResponse(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to write data: %s", err.Error())
	}
}

type PostDeleteResponse struct {
	Message string `json:"message"`
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		if err == session.ErrNoAuth {
			http.Error(w, "no authorization", http.StatusBadRequest)
		} else {
			http.Error(w, "unknown session", http.StatusInternalServerError)
		}
		return
	}

	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		if err == posts.ErrNoPost {
			http.Error(w, "there is no post with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}

	if sess.UserID != post.Author.ID {
		http.Error(w, "you are not the author of the post", http.StatusBadRequest)
		return
	}

	err = h.PostRepo.Delete(postID)
	if err != nil {
		if err == posts.ErrNoPost {
			http.Error(w, "there is no post with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to delete post", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", JSONContent)
	w.WriteHeader(http.StatusOK)

	resp := &PostDeleteResponse{Message: "success"}
	bytes, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to write data: %s", err.Error())
	}
}
