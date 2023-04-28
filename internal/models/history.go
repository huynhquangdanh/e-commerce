package models

import "time"

type History struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	UserID    int       `json:"user_id"`
	Quantity  int       `json:"quantity"`
	Discount  float64   `json:"discount"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type HistoryGet struct {
	ID          int       `json:"id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	UserName    string    `json:"first_name"`
	Discount    int       `json:"discount"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
