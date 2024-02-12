package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"redditclone/pkg/posts"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
)

type CommentCreationRequest struct {
	Comment string `json:"comment"`
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	if r.Header.Get("Content-Type") != JSONContent {
		http.Error(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	commentData := &CommentCreationRequest{}
	err = json.Unmarshal(body, commentData)
	if err != nil {
		http.Error(w, "failed to unpack payload", http.StatusBadRequest)
		return
	}
	comment := commentData.Comment

	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		if err == posts.ErrNoPost {
			http.Error(w, "there is no post with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}

	commentsRepo := post.Comments

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		if err == session.ErrNoAuth {
			http.Error(w, "no authorization", http.StatusBadRequest)
		} else {
			http.Error(w, "unknown session", http.StatusInternalServerError)
		}
		return
	}

	_, err = commentsRepo.Add(&posts.Comment{
		Author: posts.Author{
			ID:       sess.UserID,
			Username: sess.Username,
		},
		Body: comment,
	})
	if err != nil {
		http.Error(w, "failed to add comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", JSONContent)
	w.WriteHeader(http.StatusCreated)

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

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]
	commentID := vars["comment_id"]

	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		if err == posts.ErrNoPost {
			http.Error(w, "there is no post with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}

	commentsRepo := post.Comments

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		if err == session.ErrNoAuth {
			http.Error(w, "no authorization", http.StatusBadRequest)
		} else {
			http.Error(w, "unknown session", http.StatusInternalServerError)
		}
		return
	}

	comment, err := commentsRepo.GetByID(commentID)
	if err != nil {
		if err == posts.ErrNoComment {
			http.Error(w, "there is no comment with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to get comment", http.StatusInternalServerError)
		}
		return
	}

	if sess.UserID != comment.Author.ID {
		http.Error(w, "you are not the author of the comment", http.StatusBadRequest)
		return
	}

	err = commentsRepo.Delete(commentID)
	if err != nil {
		if err == posts.ErrNoComment {
			http.Error(w, "there is no comment with this id", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to delete comment", http.StatusInternalServerError)
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
