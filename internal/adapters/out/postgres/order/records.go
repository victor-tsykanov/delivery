package order

import (
	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Location  Location   `gorm:"embedded;embeddedPrefix:location_"`
	Status    string     `gorm:"type:varchar(100)"`
	CourierID *uuid.UUID `gorm:"type:uuid"`
}

type Location struct {
	X int
	Y int
}
