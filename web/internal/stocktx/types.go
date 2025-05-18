package stocktx

import "context"

type Item struct {
	ProductID int64
	Amount    int32
}

type Catalog interface {
	DecreaseProductQuantity(ctx context.Context, id int64, amount int32) error
	IncreaseProductQuantity(ctx context.Context, id int64, amount int32) error
}
