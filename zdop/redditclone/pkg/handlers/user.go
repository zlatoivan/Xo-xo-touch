package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
)

type UserHandler struct {
	UserRepo user.UserRepo
}

type NewTokenResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := &user.User{}
	err = json.Unmarshal(body, userData)
	if err != nil {
		http.Error(w, "failed to unpack payload", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.Register(userData.Username, userData.Password)
	if err != nil {
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	tokenString, err := session.CreateJWT(user.ID, user.Username)
	if err != nil {
		http.Error(w, "failed to create jwt token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := NewTokenResponse{Token: tokenString}

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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := &user.User{}
	err = json.Unmarshal(body, userData)
	if err != nil {
		http.Error(w, "failed to unpack payload", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.Authorize(userData.Username, userData.Password)
	if err != nil {
		http.Error(w, "failed to authorize", http.StatusBadRequest)
		return
	}

	tokenString, err := session.CreateJWT(user.ID, user.Username)
	if err != nil {
		http.Error(w, "failed to create jwt token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := NewTokenResponse{Token: tokenString}

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
