package main

import (
	"log/slog"

	"github.com/andrei-kozel/grpc_sso/internal/app"
	"github.com/andrei-kozel/grpc_sso/internal/config"
	"github.com/andrei-kozel/grpc_sso/internal/lib/prettylog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	log := setupLoggger(config.Env)
	log.Info("Srarting application", slog.String("env", config.Env))

	application := app.New(log, config.GRPC.Port, config.StoragePath, config.TokenTTL)
	application.GRPCServer.MustRun()
}

func setupLoggger(env string) *slog.Logger {
	log := slog.New(prettylog.NewHandler(nil))

	switch env {
	case envLocal:
		log = slog.New(prettylog.NewHandler(&slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(prettylog.NewHandler(&slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(prettylog.NewHandler(&slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
