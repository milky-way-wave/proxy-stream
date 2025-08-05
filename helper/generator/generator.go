package generator

import (
	"github.com/google/uuid"
)

func UUIDv7() string {
	id, _ := uuid.NewV7()

	return id.String()
}
