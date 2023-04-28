package repository

import (
	"backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllProducts() ([]*models.Product, error)
	OneProduct(id int) (*models.Product, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetHistoryByUser(userID int) ([]*models.HistoryGet, error)
	InsertUser(user models.User) error
	AddHistory(history *models.History) error
	SaveCoupon(coupon *models.Coupon) error
	GetCouponByCode(code string) (*models.Coupon, error)
}
