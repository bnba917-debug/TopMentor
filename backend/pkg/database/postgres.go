package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) DB() *sql.DB {
	return p.db
}
