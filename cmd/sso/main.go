package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrei-kozel/go-utils/utils/prettylog"
	"github.com/andrei-kozel/grpc_sso/internal/app"
	"github.com/andrei-kozel/grpc_sso/internal/config"
)

func main() {
	config := config.MustLoad()
	log := prettylog.SetupLoggger(config.Env)

	log.Info("Srarting application", slog.String("env", config.Env))

	application := app.New(log, config.GRPC.Port, config.StoragePath, config.TokenTTL)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))
	application.GRPCServer.Stop()
	log.Info("application stopped")
}
