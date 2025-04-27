package model

import "encoding/json"

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Price       float64         `json:"price"`
	Quantity    int32           `json:"quantity"`
	ImageURL    string          `json:"image_url,omitempty"`
	Attributes  json.RawMessage `json:"attributes"`
	CategoryID  int64           `json:"category_id"`
}
