package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/ARM-auto-place/users/internal/config"
	"github.com/ARM-auto-place/users/internal/service"
	"github.com/ARM-auto-place/users/internal/storage/postgres"
	customhandlers "github.com/ARM-auto-place/users/internal/transport/custom-handels"
	"github.com/ARM-auto-place/users/internal/transport/jsonRPC/externalapi"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}()

	err = godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}
	cfg, err := config.Values()
	if err != nil {
		return
	}
	logger := cfg.Logger()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	users, err := postgres.New(ctx, cfg.Users)
	if err != nil {
		logger.Error().Msg("failed to connect to users postgres")
		return
	}
	defer func() {
		users.GetMasterConn().Close()
		users.GetReplicaConn().Close()
	}()
	cars, err := postgres.New(ctx, cfg.Cars)
	if err != nil {
		logger.Error().Msg("failed to connect to cars postgres")
		return
	}
	defer func() {
		cars.GetMasterConn().Close()
		cars.GetReplicaConn().Close()
	}()

	storage := postgres.NewUsersStorage(users, cars)
	svc := service.NewUsersService(logger, storage)
	app := externalapi.New(logger,
		externalapi.SetReadBufferSize(cfg.MaxRequestHeaderSize),
		externalapi.MaxBodySize(cfg.MaxRequestBodySize),
		externalapi.ReadTimeout(cfg.ReadTimeout),
		externalapi.WriteTimeout(cfg.WriteTimeout),
		externalapi.WithRequestID("X-Request-Id"),
		externalapi.Use(cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool { return slices.Contains(cfg.AllowedOrigins, origin) },
			AllowMethods:     "GET,POST,OPTIONS",
			AllowHeaders:     "Content-Type,Accept,X-Request-Id",
			AllowCredentials: true,
		})),
		externalapi.Use(customhandlers.LogRequests),
		externalapi.Use(customhandlers.Authenticate(cfg.SessionKey, cfg.AllowedOrigins)),
		externalapi.ExternalAPI(externalapi.NewExternalAPI(svc)),
	)
	app.ServeHealth(cfg.HealthBind, map[string]string{"status": "ok"})

	listen := make(chan error, 1)
	go func() {
		logger.Info().Str("address", cfg.ServiceBind).Msg("HTTP server starting")
		listen <- app.Fiber().Listen(cfg.ServiceBind)
	}()

	select {
	case err = <-listen:
		return
	case <-ctx.Done():
		logger.Info().Msg("server stopping")
	}

	shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	done := make(chan struct{})
	go func() {
		app.Shutdown()
		close(done)
	}()
	select {
	case <-done:
	case <-shutdown.Done():
		err = shutdown.Err()
	}
}
