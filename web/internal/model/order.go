package model

import "time"

type Product struct {
	ID         int64
	Name       string
	Quantity   int32
	Price      float64
	TotalPrice float64
}

type Order struct {
	ID        int64
	Products  []*Product
	Price     float64
	Status    string
	CreatedAt time.Time
}
