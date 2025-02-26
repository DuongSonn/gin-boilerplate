package postgres_repository

import (
	"oauth-server/app/repository"

	"gorm.io/gorm"
)

func RegisterPostgresRepositories(db *gorm.DB) repository.RepositoryCollections {
	postgresUserRepo := NewPostgresUserRepository(db)
	postgresOAuthRepo := NewPostgresOAuthRepository(db)

	return repository.RepositoryCollections{
		UserRepo:  postgresUserRepo,
		OAuthRepo: postgresOAuthRepo,
	}
}
