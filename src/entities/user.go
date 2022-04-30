package entities

import (
	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID
	ExternalId     int64
	ExternalChatId int64
	UserName       string
	Balance        uint // in cents
}
