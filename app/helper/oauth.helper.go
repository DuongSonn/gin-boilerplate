package helper

import (
	"log/slog"
	"oauth-server/app/entity"
	"oauth-server/config"
	_jwt "oauth-server/package/jwt"
	logger "oauth-server/package/log"
	"oauth-server/utils"
)

type oauthHelper struct {
}

func newOauthHelper() OauthHelper {
	return &oauthHelper{}
}

func (h *oauthHelper) GenerateAccessToken(user entity.User) (string, error) {
	conf := config.GetConfiguration().Jwt

	payload := &_jwt.UserPayload{
		ID: user.ID,
	}
	accessToken, err := _jwt.GenerateToken(payload, conf.UserAccessTokenKey, utils.USER_ACCESS_TOKEN_IAT)
	if err != nil {
		logger.GetLogger().Info(
			"GenerateAccessToken",
			slog.Group(
				(entity.USER_TABLE_NAME),
				slog.String("email", *user.Email),
				slog.String("phone_number", *user.PhoneNumber),
			),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return accessToken, nil
}

func (h *oauthHelper) GenerateRefreshToken(user entity.User) (string, error) {
	conf := config.GetConfiguration().Jwt

	payload := &_jwt.UserPayload{
		ID: user.ID,
	}
	refreshToken, err := _jwt.GenerateToken(payload, conf.UserRefreshTokenKey, utils.USER_REFRESH_TOKEN_IAT)
	if err != nil {
		logger.GetLogger().Info(
			"GenerateRefreshToken",
			slog.Group(
				(entity.USER_TABLE_NAME),
				slog.String("email", *user.Email),
				slog.String("phone_number", *user.PhoneNumber),
			),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return refreshToken, nil
}
