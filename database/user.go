package database

import (
	"github.com/google/uuid"
)

const (
	UserFieldID        = "id"
	UserFieldFirstName = "first_name"
	UserFieldLastName  = "last_name"
	UserFieldBiography = "biography"

	firstNameMin = 2
	firstNameMax = 20
	lastNameMin  = 2
	lastNameMax  = 20
	bioMin       = 20
	bioMax       = 450
)

type User struct {
	ID        ID     `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

func (u *User) HasAnyFieldInvalid() (fields []string) {
	if !u.isValidUUID(u.ID) {
		fields = append(fields, UserFieldID)
	}

	if !u.isLenBetween(u.FirstName, firstNameMin, firstNameMax) {
		fields = append(fields, UserFieldFirstName)
	}

	if !u.isLenBetween(u.LastName, lastNameMin, lastNameMax) {
		fields = append(fields, UserFieldLastName)
	}

	if !u.isLenBetween(u.Biography, bioMin, bioMax) {
		fields = append(fields, UserFieldBiography)
	}

	return
}

func (u *User) isValidUUID(id ID) bool {
	if id == ID(uuid.Nil) {
		return false
	}

	_, err := uuid.Parse(id.ToString())
	return err == nil
}

func (u *User) isLenBetween(s string, min, max int) bool {
	l := len([]rune(s))
	return l >= min && l <= max
}
