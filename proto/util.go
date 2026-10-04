package proto

import "uuid"

func NewDataCode() string {
	id := uuid.NewV7()
	return id.String()
}
