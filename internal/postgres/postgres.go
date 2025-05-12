package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDb struct {
	Pool *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, url string) *PostgresDb {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalln("Error connecting to PostgresDB", err.Error())
	}
	return &PostgresDb{Pool: pool}
}
