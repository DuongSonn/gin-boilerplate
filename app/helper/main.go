package helper

import (
	"oauth-server/app/repository"
)

type HelperCollections struct {
	OauthHelper OauthHelper
	UserHelper  UserHelper
}

func RegisterHelpers(
	postgresRepo repository.RepositoryCollections,
) HelperCollections {
	return HelperCollections{
		OauthHelper: newOauthHelper(),
		UserHelper:  newUserHelper(postgresRepo),
	}
}
