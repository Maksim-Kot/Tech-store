package postgre

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Maksim-Kot/Tech-store-orders/config"
	"github.com/Maksim-Kot/Tech-store-orders/internal/repository"
	"github.com/Maksim-Kot/Tech-store-orders/pkg/model"

	_ "github.com/lib/pq"
)

const (
	StatusNew int64 = iota + 1
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

func (r *Repository) CreateOrder(ctx context.Context, userID int64, price float64, items []model.Item) (int64, error) {
	if len(items) == 0 {
		return 0, repository.ErrNotCreated
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int64

	orderQuery := `
		INSERT INTO orders (user_id, total_price, status_id)
		VALUES ($1, $2, $3)
		RETURNING id`

	err = tx.QueryRowContext(ctx, orderQuery, userID, price, StatusNew).Scan(&id)
	if err != nil {
		return 0, err
	}

	itemsQuery := `
		INSERT INTO order_items (order_id, item_id, quantity)
		VALUES ($1, $2, $3)`

	stmt, err := tx.PrepareContext(ctx, itemsQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, item := range items {
		if item.Quantity < 1 {
			return 0, repository.ErrNotCreated
		}
		_, err := stmt.ExecContext(ctx, id, item.ItemID, item.Quantity)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) OrderByID(ctx context.Context, id int64) (*model.Order, error) {
	if id < 1 {
		return nil, repository.ErrNotFound
	}

	query := `
		SELECT o.id, o.user_id, o.total_price, s.name AS status, o.created_at
		FROM orders o
		JOIN statuses s ON o.status_id = s.id
		WHERE o.id = $1`

	var order model.Order

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Price,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, repository.ErrNotFound
		default:
			return nil, err
		}
	}

	items, err := r.itemsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return &order, nil
}

func (r *Repository) OrdersByUserID(ctx context.Context, id int64) ([]*model.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.total_price, s.name, o.created_at
		FROM orders o
		JOIN statuses s ON o.status_id = s.id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order

	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Price,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, repository.ErrNotFound
	}

	for _, order := range orders {
		items, err := r.itemsByID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
	}

	return orders, nil
}

func (r *Repository) itemsByID(ctx context.Context, id int64) ([]model.Item, error) {
	query := `
		SELECT item_id, quantity
		FROM order_items
		WHERE order_id = $1`

	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item

	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ItemID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
