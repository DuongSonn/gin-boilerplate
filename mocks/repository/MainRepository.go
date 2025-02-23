package mocks

import (
	"oauth-server/app/repository"
	"testing"
)

func InitMockRepository(t *testing.T) repository.RepositoryCollections {
	return repository.RepositoryCollections{
		UserRepo:  NewUserRepository(t),
		OAuthRepo: NewOAuthRepository(t),
	}
}
