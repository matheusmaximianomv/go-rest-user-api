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

		idParsed, err := uuid.Parse(id)
		if err != nil {
			slog.Error("UUID inválido", "error", err)
			utils.SendJSON(w, utils.Response{Error: "UUID inválido"}, http.StatusUnprocessableEntity)
			return
		}

		if user := db.FindById(database.ID(idParsed)); user != nil {
			utils.SendJSON(w, utils.Response{Data: user}, http.StatusOK)
			return
		}

		utils.SendJSON(w, utils.Response{Data: new(any)}, http.StatusOK)
	}
}

func handlePostUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser database.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			utils.SendJSON(w, utils.Response{Error: "Não foi possível converter dados"}, http.StatusBadRequest)
			slog.Error("Não foi possível converter dados", "erro", err)
			return
		}
		defer r.Body.Close()

		invalidFields := newUser.HasAnyFieldInvalid()
		if len(invalidFields) != 0 {
			utils.SendJSON(w, utils.Response{Error: fmt.Sprintf("Campo(s) inválido(s): %v", invalidFields)}, http.StatusUnprocessableEntity)
			slog.Error("Campo inválido", slog.Any("campos", invalidFields))
			return
		}

		uuidUser, err := db.Insert(newUser)
		if uuidUser == nil {
			utils.SendJSON(w, utils.Response{Error: "ID não identificado"}, http.StatusInternalServerError)
			slog.Error("Usuário não criado")
			return
		}

		if err != nil {
			utils.SendJSON(w, utils.Response{Error: "Usuário não criado"}, http.StatusInternalServerError)
			slog.Error("Usuário não criado", slog.Any("error", err))
			return
		}

		newUser.ID = *uuidUser
		utils.SendJSON(w, utils.Response{Data: newUser}, http.StatusCreated)
	}
}

func handlePutUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idParsed, err := uuid.Parse(id)
		if err != nil {
			slog.Error("UUID inválido", "error", err)
			utils.SendJSON(w, utils.Response{Error: "UUID inválido"}, http.StatusUnprocessableEntity)
			return
		}

		uuidUser := database.ID(idParsed)
		if userExists := db.FindById(uuidUser); userExists == nil {
			slog.Error("Usuário não existe", "uuid", uuidUser)
			utils.SendJSON(w, utils.Response{Error: "Usuário não existe"}, http.StatusUnprocessableEntity)
			return
		}

		var userUpdated database.User
		if err := json.NewDecoder(r.Body).Decode(&userUpdated); err != nil {
			utils.SendJSON(w, utils.Response{Error: "Não foi possível converter dados"}, http.StatusBadRequest)
			slog.Error("Não foi possível converter dados", "erro", err)
			return
		}
		defer r.Body.Close()

		invalidFields := userUpdated.HasAnyFieldInvalid()
		if len(invalidFields) != 0 {
			utils.SendJSON(w, utils.Response{Error: fmt.Sprintf("Campo(s) inválido(s): %v", invalidFields)}, http.StatusUnprocessableEntity)
			slog.Error("Campo inválido", slog.Any("campos", invalidFields))
			return
		}

		userUpdated.ID = uuidUser
		if err := db.UpdateUser(uuidUser, userUpdated); err != nil {
			utils.SendJSON(w, utils.Response{Error: "Usuário não atualizado"}, http.StatusInternalServerError)
			slog.Error("Usuário não atualizado", slog.Any("error", err))
			return
		}

		utils.SendJSON(w, utils.Response{Data: userUpdated}, http.StatusAccepted)
	}
}

func handleDeleteUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idParsed, err := uuid.Parse(id)
		if err != nil {
			slog.Error("UUID inválido", "error", err)
			utils.SendJSON(w, utils.Response{Error: "UUID inválido"}, http.StatusUnprocessableEntity)
			return
		}

		uuidUser := database.ID(idParsed)
		if userExists := db.FindById(uuidUser); userExists == nil {
			slog.Error("Usuário não existe", "uuid", uuidUser)
			utils.SendJSON(w, utils.Response{Error: "Usuário não existe"}, http.StatusUnprocessableEntity)
			return
		}

		if err := db.DeleteUser(uuidUser); err != nil {
			utils.SendJSON(w, utils.Response{Error: "Usuário não removido"}, http.StatusInternalServerError)
			slog.Error("Usuário não removido", slog.Any("error", err))
			return
		}

		utils.SendJSON(w, utils.Response{}, http.StatusNoContent)
	}
}
