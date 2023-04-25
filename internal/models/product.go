package models

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
