package models

type Purchase struct {
	ProductID int    `json:"product_id"`
	Coupon    string `json:"coupon"`
	Quantity  int    `json:"quantity"`
}
