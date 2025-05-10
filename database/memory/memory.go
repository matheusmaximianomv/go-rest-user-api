package database_memory

import (
	"go-rest-user-api/entities"
	"sync"

	"github.com/google/uuid"
)

type DatabaseMemory struct {
	Data  map[string]entities.User
	Mutex sync.Mutex
}

func (dm *DatabaseMemory) StartStorage() error {
	dm.Data = make(map[string]entities.User, 0)

	return nil
}

func (dm *DatabaseMemory) FindAll() []entities.User {
	users := make([]entities.User, 0)

	for _, user := range dm.Data {
		users = append(users, user)
	}

	return users
}

func (dm *DatabaseMemory) FindById(id entities.ID) *entities.User {
	if user, ok := dm.Data[id.ToString()]; ok {
		return &user
	}

	return nil
}

func (dm *DatabaseMemory) Insert(user entities.User) (*entities.ID, error) {
	id := entities.ID(uuid.New())
	user.ID = id

	dm.Mutex.Lock()
	defer dm.Mutex.Unlock()

	dm.Data[id.ToString()] = user

	return &id, nil
}

func (dm *DatabaseMemory) Update(id entities.ID, user entities.User) error {
	userExist := dm.FindById(id)

	if userExist == nil {
		return nil
	}

	user.ID = id

	dm.Mutex.Lock()
	defer dm.Mutex.Unlock()

	dm.Data[id.ToString()] = user

	return nil
}

func (dm *DatabaseMemory) Delete(id entities.ID) error {
	dm.Mutex.Lock()
	defer dm.Mutex.Unlock()

	delete(dm.Data, id.ToString())

	return nil
}
