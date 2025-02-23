package service

import (
	"oauth-server/app/helper"
	"oauth-server/app/repository"
)

type ServiceCollections struct {
	UserSvc  UserService
	OAuthSvc OAuthService
}

func RegisterServices(
	helpers helper.HelperCollections,

	postgresRepo repository.RepositoryCollections,
) ServiceCollections {
	userSvc := NewUserService(helpers, postgresRepo)
	oauthSvc := NewOAuthService(helpers, postgresRepo)

	return ServiceCollections{
		UserSvc:  userSvc,
		OAuthSvc: oauthSvc,
	}
}
