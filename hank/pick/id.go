package pick

import "github.com/google/uuid"

func NewDataCode() string {
	id := uuid.Must(uuid.NewV7())
	return id.String()
}
