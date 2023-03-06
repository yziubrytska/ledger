package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Accounts struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"ID,omitempty"`
	Email        string
	Balance      int64
	Transactions []*Transactions `gorm:"foreignKey:UserID"`
}

type Transactions struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"ID,omitempty"`
	UserID uuid.UUID
	Date   pgtype.Time
	Sum    int64
}

type PgTransactions struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"ID,omitempty"`
	UserID uuid.UUID
	Date   pgtype.Time
	Sum    pgtype.Numeric
}
