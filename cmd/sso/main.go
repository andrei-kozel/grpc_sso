package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))
	application.GRPCServer.Stop()
	log.Info("Application stopped")
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
