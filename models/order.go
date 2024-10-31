package models

import "github.com/google/uuid"

type Order struct {
	ID       uuid.UUID `db:"id"`
	ItemID   uuid.UUID `db:"item_id"`
	Quantity int       `db:"quantity"`
	Status   string    `db:"status"`
}
