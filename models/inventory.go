package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID
	Name      string
	Quantity  uint
	Price     float64
	CreatedAt time.Time
}
