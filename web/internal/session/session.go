package session

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/Maksim-Kot/Tech-store-web/config"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

type Manager interface {
	Put(ctx context.Context, key string, val any)
	PopString(ctx context.Context, key string) string
	LoadAndSave(http.Handler) http.Handler
}

type scsManager struct {
	sm *scs.SessionManager
}

// New creates and configures a new scs.SessionManager instance.
// If a non-nil *sql.DB is provided, it uses a MySQL-backed session store.
// If db is nil, the session manager will use the default in-memory store.
func New(db *sql.DB, cfg config.SessionConfig) (Manager, error) {
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

func (m *scsManager) PopString(ctx context.Context, key string) string {
	return m.sm.PopString(ctx, key)
}

func (m *scsManager) LoadAndSave(next http.Handler) http.Handler {
	return m.sm.LoadAndSave(next)
}
