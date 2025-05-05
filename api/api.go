package api

import (
	"encoding/json"
	"fmt"
	"go-rest-user-api/database"
	"go-rest-user-api/utils"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func NewHander(db database.Database) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(middlewareContextType)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", handleGetUser(db))
		r.Get("/{id}", handleGetUser(db))
		r.Post("/", handlePostUser(db))
		r.Put("/{id}", handlePutUser(db))
		r.Delete("/{id}", handleDeleteUser(db))
	})

	return r
}

func middlewareContextType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func handleGetUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if len(id) == 0 {
			utils.SendJSON(w, utils.Response{Data: db.FindAll()}, http.StatusOK)
			return
		}

		uuidParsed, err := uuid.Parse(id)
		if err != nil {
			slog.Error("UUID inválido", "error", err)
			utils.SendJSON(w, utils.Response{Error: "UUID inválido"}, http.StatusUnprocessableEntity)
			return
		}

		if user := db.FindById(database.ID(uuidParsed)); user != nil {
			utils.SendJSON(w, utils.Response{Data: user}, http.StatusOK)
			return
		}

		utils.SendJSON(w, utils.Response{Data: nil}, http.StatusOK)
	}
}

func handlePostUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body database.User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.SendJSON(w, utils.Response{Error: "Não foi possível converter dados"}, http.StatusBadRequest)
			slog.Error("Não foi possível converter dados", "erro", err)
			return
		}

		invalidFields := body.HasAnyFieldInvalid()
		invalidFieldsLen := len(invalidFields)

		if invalidFieldsLen != 0 && (invalidFieldsLen > 1 || invalidFields[0] != database.UserFieldID) {
			utils.SendJSON(w, utils.Response{Error: fmt.Sprintf("Campo(s) inválido(s): %v", invalidFields)}, http.StatusUnprocessableEntity)
			slog.Error("Campo inválido", slog.Any("campos", invalidFields))
			return
		}

		idUser, err := db.StoreUser(body)
		if idUser == nil {
			utils.SendJSON(w, utils.Response{Error: "ID não identificado"}, http.StatusInternalServerError)
			slog.Error("Usuário não criado")
			return
		}

		if err != nil {
			utils.SendJSON(w, utils.Response{Error: "Usuário não criado"}, http.StatusInternalServerError)
			slog.Error("Usuário não criado", slog.Any("error", err))
			return
		}

		body.ID = *idUser
		utils.SendJSON(w, utils.Response{Data: body}, http.StatusCreated)
	}
}

func handlePutUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleDeleteUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
