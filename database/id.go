package database

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ID uuid.UUID

func (i ID) ToString() string {
	return uuid.UUID(i).String()
}

// Implementando o método MarshalJSON para serialização customizada
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String()) // Converte o UUID para string
}

// Implementando o método UnmarshalJSON para deserialização customizada
func (id *ID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*id = ID(parsedUUID) // Atribui o UUID ao tipo ID
	return nil
}
