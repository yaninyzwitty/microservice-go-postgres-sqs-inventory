package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	Id        uuid.UUID `json:"id"`
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
