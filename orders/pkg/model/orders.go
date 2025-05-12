package model

import "time"

type Item struct {
	ItemID   int64 `json:"item_id"`
	Quantity int32 `json:"quantity"`
}

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	Items     []Item    `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}
