package order

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Location  Location   `gorm:"embedded;embeddedPrefix:location_"`
	Status    string     `gorm:"type:varchar(100)"`
	CourierID *uuid.UUID `gorm:"type:uuid"`
}

type Location struct {
	X int
	Y int
}
