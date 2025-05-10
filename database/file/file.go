package database_file

import (
	"encoding/json"
	"fmt"
	"go-rest-user-api/entities"
	"os"
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	Users map[string]entities.User `json:"users"`
}

type DatabaseFile struct {
	Data  Storage
	Mutex sync.Mutex
}

func (df *DatabaseFile) StartStorage() error {
	db, err := df.getDataFromFile()
	if db == nil || err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	df.Data = *db

	return nil
}

func (df *DatabaseFile) getDataFromFile() (*Storage, error) {
	file, err := df.getFile()
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

func (df *DatabaseFile) updateFile() error {
	file, err := df.getFile()
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(df.Data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file.Name(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (df *DatabaseFile) getFile() (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(cwd+"/database/file/storage.json", os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (df *DatabaseFile) FindAll() []entities.User {
	users := make([]entities.User, 0)

	for _, user := range df.Data.Users {
		users = append(users, user)
	}

	return users
}

func (df *DatabaseFile) FindById(id entities.ID) *entities.User {
	user, ok := df.Data.Users[id.ToString()]
	if !ok {
		return nil
	}

	return &user
}

func (df *DatabaseFile) Insert(user entities.User) (*entities.ID, error) {
	id := entities.ID(uuid.New())
	user.ID = id

	df.Mutex.Lock()
	defer df.Mutex.Unlock()

	df.Data.Users[id.ToString()] = user

	if err := df.updateFile(); err != nil {
		return nil, err
	}

	return &id, nil
}

func (df *DatabaseFile) Update(id entities.ID, user entities.User) error {
	userExist := df.FindById(id)

	if userExist == nil {
		return nil
	}

	df.Mutex.Lock()
	defer df.Mutex.Unlock()

	user.ID = id
	df.Data.Users[id.ToString()] = user

	if err := df.updateFile(); err != nil {
		return err
	}

	return nil
}

func (df *DatabaseFile) Delete(id entities.ID) error {
	df.Mutex.Lock()
	defer df.Mutex.Unlock()

	delete(df.Data.Users, id.ToString())

	if err := df.updateFile(); err != nil {
		return err
	}

	return nil
}
