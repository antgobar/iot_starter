package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	store := Store{db: pool}
	store.setUpTables(ctx)
	return &store, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) setUpTables(ctx context.Context) error {
	err := s.execSQLFile(ctx, "sql/schema.sql")
	if err != nil {
		return err
	}
	log.Println("Tables created")
	return nil
}

func (s *Store) execSQLFile(ctx context.Context, path string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read SQL file: %w", err)
	}

	sql := string(sqlBytes)
	_, err = s.db.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}
