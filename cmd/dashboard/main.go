package main

import (
	"context"
	"iotstarter/internal/api"
	"iotstarter/internal/auth"
	"iotstarter/internal/config"
	"iotstarter/internal/device"
	"iotstarter/internal/measurement"
	"iotstarter/internal/middleware"
	"iotstarter/internal/pages"
	"iotstarter/internal/postgres"
	"iotstarter/internal/presentation"
	"iotstarter/internal/session"
	"iotstarter/internal/user"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbUrl := config.MustLoadEnv("DATABASE_URL")
	apiAddr := config.MustLoadEnv("API_ADDR")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	db := postgres.NewPostgresPool(ctx, dbUrl)

	userRepo := user.NewPostgresRepository(ctx, db.Pool)
	sessionRepo := session.NewPostgresRepository(ctx, db.Pool)
	deviceRepo := device.NewPostgresRepository(ctx, db.Pool)
	measurementRepo := measurement.NewPostgresRepository(ctx, db.Pool)

	userService := user.NewService(userRepo)
	sessionService := session.NewService(sessionRepo)
	authService := auth.NewService(userRepo, sessionRepo)
	deviceService := device.NewService(deviceRepo)
	measurementService := measurement.NewService(measurementRepo)
	htmlPresenter := presentation.NewHtmlPresenter()

	deviceHandler := device.NewHandler(deviceService, htmlPresenter)
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	webPageHandler := pages.NewHandler(htmlPresenter)
	measurementHandler := measurement.NewHandler(measurementService)

	middlewareStack := middleware.LoadMiddleware(sessionService)
	server := api.NewServer(
		apiAddr,
		middlewareStack,
		authHandler,
		userHandler,
		deviceHandler,
		webPageHandler,
		measurementHandler,
	)
	server.Run("Dashboard")
}
