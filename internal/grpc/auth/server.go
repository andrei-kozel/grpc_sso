package auth

import (
	"context"

	ssov1 "github.com/andrei-kozel/grpc-protobuf/gen/go/sso"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, passwrd string) (userID int, err error)
	Logout(token string) (success bool, err error)
	IsAdmin(userId int) (isAdmin bool, err error)
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "password is required")
	}
	if req.GetAppId() == 0 {
		return nil, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Code(code.Code_INTERNAL), "internal erroe")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Regisrter(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "password is required")
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Code(code.Code_INTERNAL), "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: uint64(userID),
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if req.GetUserId() == emptyValue {
		return nil, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(int(req.GetUserId()))
	if err != nil {
		return nil, status.Error(codes.Code(code.Code_INTERNAL), "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "token is required")
	}

	success, err := s.auth.Logout(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Code(code.Code_INTERNAL), "internal error")
	}

	return &ssov1.LogoutResponse{
		Success: success,
	}, nil
}
