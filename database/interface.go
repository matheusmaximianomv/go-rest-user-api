package database

import "go-rest-user-api/entities"

type Database interface {
	StartStorage() error
	FindAll() []entities.User
	FindById(id entities.ID) *entities.User
	Insert(user entities.User) (*entities.ID, error)
	UpdateUser(id entities.ID, user entities.User) error
	DeleteUser(id entities.ID) error
}
