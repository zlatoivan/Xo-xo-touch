package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/posts"
	"redditclone/pkg/user"
)

func SetHeader(header, value string, handle http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(header, value)
		handle.ServeHTTP(w, req)
	}
}

func main() {
	userRepo := user.NewMemoryRepo()
	postRepo := posts.NewPostMemoryRepo()

	userHandler := &handlers.UserHandler{
		UserRepo: userRepo,
	}

	postsHandler := &handlers.PostHandler{
		PostRepo: postRepo,
	}

	r := mux.NewRouter()

	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./static")),
	)
	r.PathPrefix("/static/").Handler(staticHandler)

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/posts/", postsHandler.List).Methods("GET")
	r.HandleFunc("/api/posts/{category_name}", postsHandler.ListByCategory).Methods("GET")
	r.HandleFunc("/api/post/{post_id}", postsHandler.GetByID).Methods("GET")
	r.HandleFunc("/api/user/{user_login}", postsHandler.ListUserPosts).Methods("GET")
	r.HandleFunc("/api/posts", middleware.Auth(postsHandler.AddPost)).Methods("POST")
	r.HandleFunc("/api/post/{post_id}", middleware.Auth(postsHandler.AddComment)).Methods("POST")
	r.HandleFunc("/api/post/{post_id}/{comment_id}", middleware.Auth(postsHandler.DeleteComment)).Methods("DELETE")
	r.HandleFunc("/api/post/{post_id}/upvote", middleware.Auth(postsHandler.Upvote)).Methods("GET")
	r.HandleFunc("/api/post/{post_id}/downvote", middleware.Auth(postsHandler.Downvote)).Methods("GET")
	r.HandleFunc("/api/post/{post_id}/unvote", middleware.Auth(postsHandler.Unvote)).Methods("GET")
	r.HandleFunc("/api/post/{post_id}", middleware.Auth(postsHandler.DeletePost)).Methods("DELETE")

	indexData, err := os.ReadFile("static/html/index.html")
	if err != nil {
		log.Fatal("failed to read html")
	}
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, string(indexData))
	}).Methods("GET")

	addr := ":8090"
	log.Printf("starting server at addr %s", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
