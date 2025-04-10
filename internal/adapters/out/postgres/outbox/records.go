package outbox

import (
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID          uuid.UUID
	Type        string
	Payload     []byte
	ProcessedAt sql.NullTime `gorm:"index"`
}

func (Message) TableName() string {
	return "outbox"
}
