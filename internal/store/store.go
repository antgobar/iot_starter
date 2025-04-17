package main

import (
	"context"
	"fmt"
	"iotstarter/internal/measurement"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewStore(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Store{db: pool, ctx: ctx}, nil
}

func (s Store) Close() {
	s.db.Close()
}

func (s Store) RegisterDevice(location string) error {
	newDevice := measurement.Device{
		Location:  location,
		CreatedAt: time.Now().UTC(),
	}
	sql := `
        INSERT INTO devices (location, created_at)
        VALUES ($1, $2)
        RETURNING id
    `
	var deviceId int
	if err := s.db.QueryRow(s.ctx, sql, newDevice.Location, newDevice.CreatedAt).Scan(&deviceId); err != nil {
		return fmt.Errorf("failed to insert device %v: %w", newDevice, err)
	}
	return nil
}

func (s Store) GetDevices() ([]measurement.Device, error) {
	sql := `SELECT id, location, created_at FROM devices`
	rows, err := s.db.Query(s.ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var devices []measurement.Device
	for rows.Next() {
		var d measurement.Device
		if err := rows.Scan(&d.ID, &d.Location, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		devices = append(devices, d)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return devices, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set. Exiting.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store, err := NewStore(ctx, dbURL)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer store.Close()

	log.Println("Successfully connected to the database.")

	err = store.execSQLFile("sql/schema.sql")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	err = store.RegisterDevice("some-location")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	devices, err := store.GetDevices()
	if err != nil {
		log.Println(err.Error())
	}
	for _, d := range devices {
		log.Printf("Device: %+v\n", d)
	}
}

func (s Store) execSQLFile(path string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read SQL file: %w", err)
	}

	sql := string(sqlBytes)
	_, err = s.db.Exec(s.ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}
