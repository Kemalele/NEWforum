package services

import (
	"crypto/rand"

	uuid "github.com/satori/go.uuid"
)

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	newId, _ := uuid.NewV4()
	return newId.String()
}
