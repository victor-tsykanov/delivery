package courier

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Courier struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(500)"`
	Location  Location  `gorm:"embedded;embeddedPrefix:location_"`
	Transport Transport `gorm:"foreignKey:CourierID"`
	Status    string    `gorm:"type:varchar(100)"`
}

type Transport struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CourierID uuid.UUID `gorm:"type:uuid"`
	Name      string    `gorm:"type:varchar(500)"`
	Speed     int
}

type Location struct {
	X int
	Y int
}
