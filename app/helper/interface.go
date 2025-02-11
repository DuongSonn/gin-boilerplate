package helper

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/model"
)

type OauthHelper interface {
	GenerateAccessToken(user entity.User) (string, error)
	GenerateRefreshToken(user entity.User) (string, error)
}

type UserHelper interface {
	CreateUser(ctx context.Context, data *model.RegisterRequest) error
}
