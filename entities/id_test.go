package entities_test

import (
	"encoding/json"
	"go-rest-user-api/entities"
	"go-rest-user-api/utils"
	"testing"

	"github.com/google/uuid"
)

func TestID(t *testing.T) {
	t.Run("should return an error when invalid UUID string", func(t *testing.T) {
		newUUID := uuid.New()
		id := entities.ID(newUUID)

		utils.Assert(t, newUUID.String(), id.ToString())
	})

	t.Run("should receive a string when calling the MarshalJSON function", func(t *testing.T) {
		newUUID := uuid.New()
		id := entities.ID(newUUID)

		jsonBytes, err := json.Marshal(id)

		utils.Assert(t, `"`+newUUID.String()+`"`, string(jsonBytes))
		utils.Assert(t, err, nil)
	})

	t.Run("should return an unmarshal id match uuid original when called the function", func(t *testing.T) {
		newUUID := uuid.New()
		jsonData := `"` + newUUID.String() + `"`

		var id entities.ID
		err := json.Unmarshal([]byte(jsonData), &id)

		utils.Assert(t, err, nil)
		utils.Assert(t, entities.ID(newUUID), id)
	})

	t.Run("should an error when received as invalid UUID string", func(t *testing.T) {
		invalidJSON := `"invalid-uuid"`

		var id entities.ID
		err := json.Unmarshal([]byte(invalidJSON), &id)
		utils.Assert(t, err.Error(), "invalid UUID length: 12")
	})
}
