package session

import (
	"context"
	"database/sql"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/Maksim-Kot/Tech-store-web/config"
	"github.com/Maksim-Kot/Tech-store-web/internal/model"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

type Manager interface {
	Put(ctx context.Context, key string, val any)
	Get(ctx context.Context, key string) any
	GetInt64(ctx context.Context, key string) int64
	PopString(ctx context.Context, key string) string
	RenewToken(ctx context.Context) error
	Remove(ctx context.Context, key string)
	Exists(ctx context.Context, key string) bool
	LoadAndSave(http.Handler) http.Handler
}

type scsManager struct {
	sm *scs.SessionManager
}

// New creates and configures a new scs.SessionManager instance.
// If a non-nil *sql.DB is provided, it uses a MySQL-backed session store.
// If db is nil, the session manager will use the default in-memory store.
func New(db *sql.DB, cfg config.SessionConfig) (Manager, error) {
	gob.Register(model.Cart{})
	gob.Register(model.Item{})
	gob.Register(map[int64]model.Item{})

	sm := scs.New()
	lifetime, err := time.ParseDuration(cfg.Lifetime)
	if err != nil {
		return nil, err
	}
	sm.Lifetime = lifetime

	if db != nil {
		sm.Store = mysqlstore.New(db)
	}

	return &scsManager{sm}, nil
}

func (m *scsManager) Put(ctx context.Context, key string, val any) {
	m.sm.Put(ctx, key, val)
}

func (m *scsManager) Get(ctx context.Context, key string) any {
	return m.sm.Get(ctx, key)
}

func (m *scsManager) GetInt64(ctx context.Context, key string) int64 {
	return m.sm.GetInt64(ctx, key)
}

func (m *scsManager) PopString(ctx context.Context, key string) string {
	return m.sm.PopString(ctx, key)
}

func (m *scsManager) RenewToken(ctx context.Context) error {
	return m.sm.RenewToken(ctx)
}

func (m *scsManager) Remove(ctx context.Context, key string) {
	m.sm.Remove(ctx, key)
}

func (m *scsManager) Exists(ctx context.Context, key string) bool {
	return m.sm.Exists(ctx, key)
}

func (m *scsManager) LoadAndSave(next http.Handler) http.Handler {
	return m.sm.LoadAndSave(next)
}
