package models

import (
	"github.com/google/uuid"
	"time"
)

type Accounts struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"ID,omitempty"`
	Email        string
	Balance      int64
	Transactions []*Transactions `gorm:"foreignKey:UserID"`
}

type Transactions struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"ID,omitempty"`
	UserID uuid.UUID `gorm:"foreignKey:ID"`
	Date   time.Time
	Sum    int64
}
