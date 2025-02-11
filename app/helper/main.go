package helper

import (
	postgres_repository "oauth-server/app/repository/postgres"
)

type HelperCollections struct {
	OauthHelper OauthHelper
	UserHelper  UserHelper
}

func RegisterHelpers(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) HelperCollections {
	return HelperCollections{
		OauthHelper: NewOauthHelper(),
		UserHelper:  NewUserHelper(postgresRepo),
	}
}
