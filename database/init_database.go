package database

import (
	"errors"
)

func InitDatabase(db Database) (Database, error) {
	if db == nil {
		return nil, errors.New("unable to initialize the database type implementation")
	}

	if err := db.StartStorage(); err != nil {
		return nil, err
	}

	return db, nil
}
