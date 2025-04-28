package main

import (
	"context"
	"iotstarter/internal/api"
	"iotstarter/internal/auth"
	"iotstarter/internal/config"
	"iotstarter/internal/database"
	"iotstarter/internal/device"
	"iotstarter/internal/middleware"
	"iotstarter/internal/presentation"
	"iotstarter/internal/session"
	"iotstarter/internal/user"
	"iotstarter/internal/web"
	"time"
)

func main() {
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	apiAddr := config.MustLoadEnv("API_ADDR")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	db := database.NewPostgresPool(ctx, dbUrl)

	userRepo := user.NewPostgresRepository(ctx, db.Pool)
	sessionRepo := session.NewPostgresRepository(ctx, db.Pool)
	deviceRepo := device.NewPostgresRepository(ctx, db.Pool)

	userService := user.NewService(userRepo)
	sessionService := session.NewService(sessionRepo)
	authService := auth.NewService(userRepo, sessionRepo)
	deviceService := device.NewService(deviceRepo)

	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	deviceHandler := device.NewHandler(deviceService)

	htmlPresentation := presentation.NewHtmlPresentation()
	webPageHandler := web.NewHandler(htmlPresentation)

	middlewareStack := middleware.LoadMiddleware(sessionService)
	server := api.NewServer(
		apiAddr,
		middlewareStack,
		authHandler,
		userHandler,
		deviceHandler,
		webPageHandler,
	)
	server.Run("Dashboard")
}
