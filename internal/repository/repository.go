package repository

import (
	"backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllProducts() ([]*models.Product, error)
	GetUserByEmail(email string) (*models.User, error)
}
