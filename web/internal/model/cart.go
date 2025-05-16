package model

type Item struct {
	ID       int64
	Name     string
	Quantity int32
}

type Cart struct {
	UserID int64
	Items  map[int64]Item
}
