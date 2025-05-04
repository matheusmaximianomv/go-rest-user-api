package api

import (
	"go-rest-user-api/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHander(database database.Database) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", handleGetUser)
		r.Get("/{id}", handleGetUser)
		r.Post("/", handlePostUser)
		r.Put("/{id}", handlePutUser)
		r.Delete("/{id}", handleDeleteUser)
	})

	return r
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {

}

func handlePostUser(w http.ResponseWriter, r *http.Request) {

}

func handlePutUser(w http.ResponseWriter, r *http.Request) {

}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {

}
