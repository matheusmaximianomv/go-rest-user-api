package api

import (
	"encoding/json"
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
		r.Get("/", handleGetUsers(db))
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

func handleGetUsers(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := db.FindAll()
		if users == nil {
			utils.SendJSON(w, utils.Response{Message: "The users information could not be retrieved"}, http.StatusInternalServerError)
			return
		}

		utils.SendJSON(w, utils.Response{Data: users}, http.StatusOK)
	}
}

func handleGetUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idParsed, err := uuid.Parse(id)
		if err != nil {
			slog.Error("The UUID sent is not valid information", "error", err)
			utils.SendJSON(w, utils.Response{Message: "The users information could not be retrieved"}, http.StatusInternalServerError)
			return
		}

		if user := db.FindById(database.ID(idParsed)); user != nil {
			utils.SendJSON(w, utils.Response{Data: user}, http.StatusOK)
			return
		}

		utils.SendJSON(w, utils.Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
	}
}

func handlePostUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser database.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			utils.SendJSON(w, utils.Response{Message: "Please provide FirstName LastName and bio for the user"}, http.StatusBadRequest)
			slog.Error("Could not convert submitted data to expected structure", "erro", err)
			return
		}
		defer r.Body.Close()

		invalidFields := newUser.HasAnyFieldInvalid()
		if len(invalidFields) != 0 {
			utils.SendJSON(w, utils.Response{Message: "Please provide FirstName LastName and bio for the user"}, http.StatusBadRequest)
			slog.Error("Some field sent is incompatible with the specifications", "fields", invalidFields)
			return
		}

		uuidUser, err := db.Insert(newUser)
		if uuidUser == nil {
			utils.SendJSON(w, utils.Response{Message: "There was an error while saving the user to the database"}, http.StatusInternalServerError)
			slog.Error("Unable to generate an ID for the new user")
			return
		}

		if err != nil {
			utils.SendJSON(w, utils.Response{Message: "There was an error while saving the user to the database"}, http.StatusInternalServerError)
			slog.Error("Unable to create a new user due to a database error", "error", err)
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
			slog.Error("The UUID sent is not valid information", "error", err)
			utils.SendJSON(w, utils.Response{Message: "The users information could not be retrieved"}, http.StatusNotFound)
			return
		}

		uuidUser := database.ID(idParsed)
		if userExists := db.FindById(uuidUser); userExists == nil {
			slog.Error("The user with this UUID has no record", "uuid", uuidUser)
			utils.SendJSON(w, utils.Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		var userUpdated database.User
		if err := json.NewDecoder(r.Body).Decode(&userUpdated); err != nil {
			slog.Error("Could not convert submitted data to expected structure", "erro", err)
			utils.SendJSON(w, utils.Response{Message: "Please provide name and bio for the user"}, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		invalidFields := userUpdated.HasAnyFieldInvalid()
		if len(invalidFields) != 0 {
			slog.Error("Some field sent is incompatible with the specifications", "fields", invalidFields)
			utils.SendJSON(w, utils.Response{Message: "Please provide name and bio for the user"}, http.StatusBadRequest)
			return
		}

		userUpdated.ID = uuidUser
		if err := db.UpdateUser(uuidUser, userUpdated); err != nil {
			utils.SendJSON(w, utils.Response{Message: "The user information could not be modified"}, http.StatusInternalServerError)
			slog.Error("User could not be updated", "error", err)
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
			slog.Error("The UUID sent is not well formatted", "error", err)
			utils.SendJSON(w, utils.Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		uuidUser := database.ID(idParsed)
		if userExists := db.FindById(uuidUser); userExists == nil {
			slog.Error("The user with this UUID has no record", "uuid", uuidUser)
			utils.SendJSON(w, utils.Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		if err := db.DeleteUser(uuidUser); err != nil {
			utils.SendJSON(w, utils.Response{Message: "The user could not be removed"}, http.StatusInternalServerError)
			slog.Error("User could not be removed", "error", err)
			return
		}

		utils.SendJSON(w, utils.Response{}, http.StatusNoContent)
	}
}
