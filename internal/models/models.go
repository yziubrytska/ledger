package models

import "github.com/google/uuid"

type Transaction struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"ID,omitempty"`
	Date        string    `json:"date"`
	Transaction string    `json:"transaction"`
}

type Transactions []Transaction
