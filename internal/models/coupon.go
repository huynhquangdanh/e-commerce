package models

import "time"

type Coupon struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	UserID    int       `json:"user_id"`
	Code      string    `json:"code"`
	Rate      float64   `json:"rate"`
	ExpireAt  time.Time `json:"expire_at"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
