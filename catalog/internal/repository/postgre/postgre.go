package postgre

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Maksim-Kot/Tech-store-catalog/config"
	"github.com/Maksim-Kot/Tech-store-catalog/internal/repository"
	"github.com/Maksim-Kot/Tech-store-catalog/pkg/model"

	_ "github.com/lib/pq"
)

type Repository struct {
	DB *sql.DB
}

func New(cfg config.DatabaseConfig) (*Repository, error) {
	db, err := sql.Open("postgres", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Repository{db}, nil
}

func (r *Repository) Close() error {
	return r.DB.Close()
}

func (r *Repository) Categories(ctx context.Context) ([]*model.Category, error) {
	query := `
		SELECT id, name
		FROM categories
		ORDER BY name`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) ProductsByCategoryID(ctx context.Context, id int64) ([]*model.Product, error) {
	query := `
		SELECT id, name, description, price, quantity, image_url, attributes, category_id
		FROM items
		WHERE category_id = $1
		ORDER BY id`

	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*model.Product{}

	for rows.Next() {
		var product model.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.ImageURL,
			&product.Attributes,
			&product.CategoryID,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, repository.ErrNotFound
	}

	return products, nil
}

func (r *Repository) ProductByID(ctx context.Context, id int64) (*model.Product, error) {
	if id < 1 {
		return nil, repository.ErrNotFound
	}

	query := `
		SELECT id, name, description, price, quantity, image_url, attributes, category_id
		FROM items
		WHERE id = $1`

	var product model.Product

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity,
		&product.ImageURL,
		&product.Attributes,
		&product.CategoryID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, repository.ErrNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (r *Repository) DecreaseProductQuantity(ctx context.Context, id int64, amount int32) error {
	queryExist := `SELECT quantity FROM items WHERE id = $1`
	var quantity int32

	err := r.DB.QueryRowContext(ctx, queryExist, id).Scan(&quantity)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrNotFound
		default:
			return err
		}
	}

	if quantity < amount {
		return repository.ErrNotEnough
	}

	query := `
		UPDATE items
		SET quantity = quantity - $3
		WHERE id = $1 AND quantity = $2`

	res, err := r.DB.ExecContext(ctx, query, id, quantity, amount)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repository.ErrEditConflict
	}

	return nil
}

func (r *Repository) IncreaseProductQuantity(ctx context.Context, id int64, amount int32) error {
	query := `
		UPDATE items
		SET quantity = quantity + $2
		WHERE id = $1`

	res, err := r.DB.ExecContext(ctx, query, id, amount)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *Repository) PutCategory(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id`

	err := r.DB.QueryRowContext(ctx, query, category.Name).Scan(&category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) PutProduct(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO items (name, description, price, quantity, image_url, attributes, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	args := []any{
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
		product.ImageURL,
		product.Attributes,
		product.CategoryID,
	}
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)

	if err != nil {
		return err
	}
	return nil
}
