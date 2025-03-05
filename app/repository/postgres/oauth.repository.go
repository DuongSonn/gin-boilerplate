package postgres_repository

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	"time"

	"gorm.io/gorm"
)

type oAuthRepository struct {
	db *gorm.DB
}

func NewPostgresOAuthRepository(db *gorm.DB) repository.OAuthRepository {
	return &oAuthRepository{
		db,
	}
}

func (r *oAuthRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.OAuth,
) (*entity.OAuth, error) {
	if tx != nil {
		return oauth, tx.WithContext(ctx).Create(&oauth).Error
	}

	return oauth, r.db.WithContext(ctx).Create(&oauth).Error
}

func (r *oAuthRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.OAuth,
) (*entity.OAuth, error) {
	oauth.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return oauth, tx.WithContext(ctx).Save(&oauth).Error
	}

	return oauth, r.db.WithContext(ctx).Save(&oauth).Error
}

func (r *oAuthRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOAuthByFilter,
) (*entity.OAuth, error) {
	var data *entity.OAuth

	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.Token != nil {
		query = query.Where("token = ?", *filter.Token)
	}
	if filter.UserID != nil {
		query = query.Scopes(whereBy(*filter.UserID, "user_id"))
	}
	if filter.PlatForm != nil {
		query = query.Where("platform = ?", *filter.PlatForm)
	}

	err := query.First(&data).Error
	return data, err
}
