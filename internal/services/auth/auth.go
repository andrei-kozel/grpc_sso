package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/andrei-kozel/grpc-protobuf/gen/go/sso"
	"github.com/andrei-kozel/grpc_sso/internal/domain/models"
	"github.com/andrei-kozel/grpc_sso/internal/lib/jwt"
	"github.com/andrei-kozel/grpc_sso/internal/services/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (id int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var ErrorInvalidCredentials = errors.New("invalid credentials")

func New(log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider, tokenTTL: tokenTTL,
	}
}

func (a *Auth) IsAdmin(userId int) (isAdmin bool, err error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("user_id", userId),
	)

	log.Info("checking if user is admin")

	isAdmin, err = a.userProvider.IsAdmin(context.Background(), int64(userId))
	if err != nil {
		log.Error("failed to check if user is admin", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user is admin")
	return isAdmin, nil
}

// Login implements auth.Auth.
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("app_id", fmt.Sprint(appID)),
	)

	log.Info("attempting to login a user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			a.log.Warn("user not found")

			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		a.log.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid password")

		return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrorAppNotFound) {
			a.log.Warn("app not found")
			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}
		a.log.Error("failed to get app", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// Logout implements auth.Auth.
func (a *Auth) Logout(token string) (success bool, err error) {
	panic("unimplemented")
}

// Register implements auth.Auth.
func (a *Auth) Register(ctx context.Context, email string, passwrd string) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(passwrd), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}
