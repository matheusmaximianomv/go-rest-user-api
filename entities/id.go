package entities

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ID uuid.UUID

func (id ID) ToString() string {
	return uuid.UUID(id).String()
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	*id = ID(parsedUUID)
	return nil
}
