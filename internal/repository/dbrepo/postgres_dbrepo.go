package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type PostresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostresDBRepo) AllProducts() ([]*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, name, coalesce(image, ''),
			created_at, updated_at
		from
		    products
		order by
		    name
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (m *PostresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
