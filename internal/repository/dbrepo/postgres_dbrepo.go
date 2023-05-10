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
			id, product_name, price, description,
			created_at, updated_at
		from
		    products
		order by
		    product_name
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

func (m *PostresDBRepo) OneProduct(id int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, product_name, price, description, created_at, updated_at
			from products where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var product models.Product

	err := row.Scan(
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

	return &product, nil

	// get
	//query = `select `
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

func (m *PostresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

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

func (m *PostresDBRepo) GetHistoryByUser(userID int) ([]*models.HistoryGet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select histories.id, users.first_name, products.product_name, quantity, discount_applied,
			histories.created_at from histories left join products on histories.product_id = products.id
			left join users on histories.user_id = users.id where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var histories []*models.HistoryGet
	for rows.Next() {
		var history models.HistoryGet
		err := rows.Scan(
			&history.ID,
			&history.UserName,
			&history.ProductName,
			&history.Quantity,
			&history.Discount,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		histories = append(histories, &history)
	}

	query = `select id, product_name, price from history`

	return histories, nil
}

func (m *PostresDBRepo) InsertUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `insert into users (first_name, last_name, email, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, sql, user.FirstName, user.LastName, user.Email, user.Password, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *PostresDBRepo) AddHistory(history *models.History) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `insert into histories (product_id, user_id, quantity, discount_applied, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, sql, history.ProductID, history.UserID, history.Quantity, history.Discount, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *PostresDBRepo) SaveCoupon(coupon *models.Coupon) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `insert into coupons (code, rate, product_id, user_id, expired_at, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, sql, coupon.Code, coupon.Rate, coupon.ProductID, coupon.UserID, coupon.ExpireAt, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *PostresDBRepo) GetCouponByCode(code string) (*models.Coupon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from coupons where code = $1`

	var coupon models.Coupon
	row := m.DB.QueryRowContext(ctx, query, code)

	err := row.Scan(
		&coupon.ID,
		&coupon.Code,
		&coupon.Rate,
		&coupon.ProductID,
		&coupon.UserID,
		&coupon.Active,
		&coupon.ExpireAt,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (m *PostresDBRepo) DeactivateCoupon(couponCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update coupons set active = $1 where code = $2`

	_, err := m.DB.ExecContext(ctx, stmt, false, couponCode)
	if err != nil {
		return err
	}

	return nil
}
