package mysql_repository

import (
	"gin-boilerplate/app/repository"

	"gorm.io/gorm"
)

type MysqlRepositoryCollections struct {
	MysqlUserRepo  repository.UserRepository
	MysqlOAuthRepo repository.OAuthRepository
}

func RegisterMysqlRepositories(db *gorm.DB) MysqlRepositoryCollections {
	mysqlUserRepo := NewMysqlUserRepository(db)
	mysqlOAuthRepo := NewMysqlOAuthRepository(db)

	return MysqlRepositoryCollections{
		MysqlUserRepo:  mysqlUserRepo,
		MysqlOAuthRepo: mysqlOAuthRepo,
	}
}
