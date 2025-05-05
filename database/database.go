package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Storage struct {
	Users map[string]User `json:"users"`
}

type Database struct {
	Data Storage
}

func (a *Database) startStorage() error {
	db, err := a.getDataFromFile()
	if db == nil || err != nil {
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

	file, err := os.OpenFile(cwd+"/database/storage.json", os.O_RDWR, 0644)
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

func (a *Database) StoreUser(user User) (*ID, error) {
	id := ID(uuid.New())
	user.ID = id

	a.Data.Users[id.ToString()] = user

	if err := a.updateFile(); err != nil {
		return nil, err
	}

	return &id, nil
}

func (a *Database) UpdateUser(id ID, user User) error {
	userExist := a.FindById(id)

	if userExist == nil {
		return errors.New("usuário não encontrado")
	}

	a.Data.Users[id.ToString()] = user

	if err := a.updateFile(); err != nil {
		return err
	}

	return nil
}

func (a *Database) DeleteUser(id ID) error {
	delete(a.Data.Users, id.ToString())

	if err := a.updateFile(); err != nil {
		return err
	}

	return nil
}

func InitDatabase() (Database, error) {
	db := Database{}

	if err := db.startStorage(); err != nil {
		return db, err
	}

	return db, nil
}
