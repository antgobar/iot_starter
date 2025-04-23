package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, url string) *PostgresStore {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalln("Error connecting to PostgresDB", err.Error())
	}
	store := PostgresStore{db: pool}

	err = store.setUpTables(ctx)
	if err != nil {
		log.Fatalln("Error setting up tables", err.Error())
	}
	return &store
}

func (s *PostgresStore) Close() {
	s.db.Close()
}

func (s *PostgresStore) setUpTables(ctx context.Context) error {
	err := s.execSQLFile(ctx, "sql/schema.sql")
	if err != nil {
		return err
	}
	log.Println("Tables created")
	return nil
}

func (s *PostgresStore) execSQLFile(ctx context.Context, path string) error {
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
