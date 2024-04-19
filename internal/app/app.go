package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/andrei-kozel/grpc_sso/internal/app/grpc"
	"github.com/andrei-kozel/grpc_sso/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	authService := auth.New(log, nil, nil, nil, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
