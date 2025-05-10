package database

import (
	"encoding/json"
	"fmt"
	"go-rest-user-api/entities"
	"os"

	"github.com/google/uuid"
)

type Storage struct {
	Users map[string]entities.User `json:"users"`
}

type Database struct {
	Data Storage
}

func (a *Database) startStorage() error {
	db, err := a.getDataFromFile()
	if db == nil || err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
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

	data, err := json.MarshalIndent(a.Data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file.Name(), data, 0644)
	if err != nil {
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

func (a *Database) FindAll() []entities.User {
	users := make([]entities.User, 0)

	for _, user := range a.Data.Users {
		users = append(users, user)
	}

	return users
}

func (a *Database) FindById(id entities.ID) *entities.User {
	user, ok := a.Data.Users[id.ToString()]
	if !ok {
		return nil
	}

	return &user
}

func (a *Database) Insert(user entities.User) (*entities.ID, error) {
	id := entities.ID(uuid.New())
	user.ID = id

	a.Data.Users[id.ToString()] = user

	if err := a.updateFile(); err != nil {
		return nil, err
	}

	return &id, nil
}

func (a *Database) UpdateUser(id entities.ID, user entities.User) error {
	userExist := a.FindById(id)

	if userExist == nil {
		return nil
	}

	user.ID = id
	a.Data.Users[id.ToString()] = user

	if err := a.updateFile(); err != nil {
		return err
	}

	return nil
}

func (a *Database) DeleteUser(id entities.ID) error {
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
