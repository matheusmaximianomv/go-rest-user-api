package database

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type User struct {
	ID        ID     `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type Storage struct {
	Users map[string]User `json:"users"`
}

type Database struct {
	Data Storage
}

func (a *Database) StartStorage() error {
	db, err := a.getDataFromFile()
	if db == nil || err != nil {
		slog.Error("Não foi possível ler o banco", slog.Any("error", err))
		return fmt.Errorf("não foi possível ler o banco: %w", err)
	}

	a.Data = *db

	return nil
}

func (a *Database) getDataFromFile() (*Storage, error) {
	file, err := a.getFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var db *Storage
	if err := json.NewDecoder(file).Decode(&db); err != nil {
		return nil, err
	}

	return db, nil
}

func (a *Database) updateFile() error {
	file, err := a.getFile()
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(a.Data); err != nil {
		return err
	}

	return nil
}

func (a *Database) getFile() (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(cwd + "/database/storage.json")
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (a *Database) FindAll() []User {
	var users []User

	for _, user := range a.Data.Users {
		users = append(users, user)
	}

	return users
}

func (a *Database) FindById(id ID) *User {
	user, ok := a.Data.Users[id.ToString()]
	if !ok {
		return nil
	}

	return &user
}

func (a *Database) StoreUser(user User) ID {
	id := ID(uuid.New())
	user.ID = id

	a.Data.Users[id.ToString()] = user
	a.updateFile()

	return id
}

func (a *Database) UpdateUser(id ID, user User) {
	userExist := a.FindById(id)

	if userExist == nil {
		return
	}

	a.Data.Users[id.ToString()] = user
	a.updateFile()
}

func (a *Database) DeleteUser(id ID) {
	delete(a.Data.Users, id.ToString())
	a.updateFile()
}

func InitDatabase() (*Database, error) {
	application := &Database{}

	if err := application.StartStorage(); err != nil {
		return nil, err
	}

	return application, nil
}
