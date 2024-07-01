package converter

import (
	"github.com/d1mitrii/authentication-service/internal/models"
	desc "github.com/d1mitrii/authentication-service/pkg/auth/v1"
)

// Convert api login request to user model for service layer
func LoginReqToUserModel(data *desc.LoginRequest) (models.User, error) {
	if len(data.Email) == 0 {
		return models.User{}, ErrEmptyEmail
	}
	if len(data.Password) == 0 {
		return models.User{}, ErrEmptyPassword
	}
	return models.User{
		Email:    data.Email,
		Password: data.Password,
	}, nil
}

// Convert api register request to user model for service layer
func RegisterReqToUserModel(data *desc.RegisterRequest) (models.User, error) {
	if len(data.Email) == 0 {
		return models.User{}, ErrEmptyEmail
	}
	if len(data.Password) == 0 {
		return models.User{}, ErrEmptyPassword
	}
	return models.User{
		Email:    data.Email,
		Password: data.Password,
	}, nil
}
