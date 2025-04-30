package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Maksim-Kot/Tech-store-web/config"
	"github.com/Maksim-Kot/Tech-store-web/internal/repository"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	DB *sql.DB
}

func New(cfg config.DatabaseConfig) (*Repository, error) {
	db, err := sql.Open("mysql", cfg.Dsn)
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

func (r *Repository) Insert(ctx context.Context, name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (name, email, hashed_password, created) 
		VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = r.DB.ExecContext(ctx, query, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return repository.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (r *Repository) Authenticate(ctx context.Context, email, password string) (int64, error) {
	var id int64
	var hashedPassword []byte

	query := `
		SELECT id, hashed_password 
		FROM users 
		WHERE email = ?`

	err := r.DB.QueryRowContext(ctx, query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repository.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, repository.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (r *Repository) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}
