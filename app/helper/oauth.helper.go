package helper

import (
	"log/slog"
	"oauth-server/app/entity"
	"oauth-server/config"
	logger "oauth-server/package/log"
	"oauth-server/utils"
)

type oauthHelper struct {
}

func NewOauthHelper() OauthHelper {
	return &oauthHelper{}
}

func (h *oauthHelper) GenerateAccessToken(user entity.User) (string, error) {
	conf := config.GetConfiguration().Jwt

	payload := &utils.UserPayload{
		ID: user.ID,
	}
	accessToken, err := utils.GenerateToken(payload, conf.UserAccessTokenKey, utils.USER_ACCESS_TOKEN_IAT)
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

	payload := &utils.UserPayload{
		ID: user.ID,
	}
	refreshToken, err := utils.GenerateToken(payload, conf.UserRefreshTokenKey, utils.USER_REFRESH_TOKEN_IAT)
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
