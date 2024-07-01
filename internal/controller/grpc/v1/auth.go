package v1

import (
	"context"
	"github.com/d1mitrii/authentication-service/internal/converter"
	"github.com/d1mitrii/authentication-service/internal/models"
	"github.com/d1mitrii/authentication-service/internal/services"
	desc "github.com/d1mitrii/authentication-service/pkg/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Register(context.Context, models.User) (int, error)
	Login(context.Context, models.User) (models.Token, error)
	RefreshSession(context.Context, string) (models.Token, error)
	Logout(context.Context, string) error
}

type Auth struct {
	desc.UnimplementedAuthV1Server
	service AuthService
}

func NewAuth(service AuthService) *Auth {
	return &Auth{
		service: service,
	}
}

func (a *Auth) Register(ctx context.Context, req *desc.RegisterRequest) (*desc.RegisterResponse, error) {
	user, err := converter.RegisterReqToUserModel(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userId, err := a.service.Register(ctx, user)
	if err != nil {
		switch err {
		case services.ErrHashing:
			return nil, status.Error(codes.Canceled, err.Error())
		case services.ErrUserAlreadyExist:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}
	return &desc.RegisterResponse{UserId: int64(userId)}, nil
}

func (a *Auth) Login(ctx context.Context, req *desc.LoginRequest) (*desc.Token, error) {
	user, err := converter.LoginReqToUserModel(req)
	if err != nil {
		return &desc.Token{}, status.Error(codes.InvalidArgument, err.Error())
	}
	token, err := a.service.Login(ctx, user)
	if err != nil {
		switch err {
		case services.ErrIncorrectPassword:
			return &desc.Token{}, status.Error(codes.InvalidArgument, err.Error())
		case services.ErrUserNotFound:
			return &desc.Token{}, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Aborted, "internal server error")
		}
	}
	return &desc.Token{
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
	}, nil
}

func (a *Auth) Refresh(ctx context.Context, req *desc.RefreshRequest) (*desc.Token, error) {
	if len(req.RefreshToken) == 0 {
		return &desc.Token{}, status.Error(codes.InvalidArgument, "empty refresh token provided")
	}
	token, err := a.service.RefreshSession(ctx, req.RefreshToken)
	if err != nil {
		return &desc.Token{}, status.Error(codes.NotFound, "refresh session not found")
	}
	return &desc.Token{
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
	}, nil
}

func (a *Auth) Logout(ctx context.Context, req *desc.LogoutRequest) (*desc.LogoutResponse, error) {
	if len(req.RefreshToken) == 0 {
		return &desc.LogoutResponse{}, status.Error(codes.InvalidArgument, "empty refresh token provided")
	}
	err := a.service.Logout(ctx, req.RefreshToken)
	if err != nil {
		return &desc.LogoutResponse{}, status.Error(codes.NotFound, "refresh session not found")
	}
	return &desc.LogoutResponse{
		Success: true,
	}, nil
}
