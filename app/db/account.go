package db

import (
	"errors"
	"fmt"

	uuid "github.com/gofrs/uuid"
)

var (
	ErrNoNameValid = errors.New("Name invalid")
)

// Account model example
type Account struct {
	ID   int       `json:"id" bson:"id"`
	Name string    `json:"name" bson:"name"`
	UUID uuid.UUID `json:"uuid" format:"uuid" bson:"uuid"`
}

func (a Account) String() string {
	return fmt.Sprintf("Account %d, Name %s", a.ID, a.Name)
}
