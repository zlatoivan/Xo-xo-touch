package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"redditclone/pkg/posts"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
)

func (h *PostHandler) Upvote(w http.ResponseWriter, r *http.Request) {
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

	votesRepo := post.Votes

	err = votesRepo.Upvote(sess.UserID)
	if err != nil {
		http.Error(w, "failed to upvote post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

func (h *PostHandler) Downvote(w http.ResponseWriter, r *http.Request) {
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

	votesRepo := post.Votes
	err = votesRepo.Downvote(sess.UserID)
	if err != nil {
		http.Error(w, "failed to downvote post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

func (h *PostHandler) Unvote(w http.ResponseWriter, r *http.Request) {
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

	votesRepo := post.Votes
	err = votesRepo.Unvote(sess.UserID)
	if err != nil {
		http.Error(w, "failed to unvote post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
