package postgres_repository

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) (*entity.User, error) {
	if tx != nil {
		return user, tx.WithContext(ctx).Create(&user).Error
	}

	return user, r.db.WithContext(ctx).Create(&user).Error
}

func (r *userRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) (*entity.User, error) {
	user.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return user, tx.WithContext(ctx).Save(&user).Error
	}

	return user, r.db.WithContext(ctx).Save(&user).Error
}

func (r *userRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&user).Error
	}

	return r.db.WithContext(ctx).Delete(&user).Error
}

func (r *userRepository) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	users []entity.User,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&users).Error
	}

	return r.db.WithContext(ctx).Create(&users).Error
}

func (r *userRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) (*entity.User, error) {
	var user *entity.User
	err := r.db.WithContext(ctx).First(&user).Error
	return user, err
}

func (r *userRepository) FindManyByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filer *repository.FindUserByFilter,
) ([]entity.User, error) {
	var user []entity.User
	err := r.buildFilter(ctx, tx, filer).Find(&user).Error
	return user, err
}

// -------------------------------------------------------------------------------
func (r *userRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.Email != nil && *filter.Email != "" {
		query = query.Scopes(whereBy(*filter.Email, "email"))
	}
	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Scopes(whereBy(*filter.PhoneNumber, "phone_number"))
	}
	if filter.ID != nil && *filter.ID != uuid.Nil {
		query = query.Scopes(whereBy[uuid.UUID](*filter.ID, "id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(whereBySlice(filter.IDs, "id"))
	}
	if filter.Emails != nil && len(filter.Emails) > 0 {
		query = query.Scopes(whereBySlice(filter.Emails, "email"))
	}
	if filter.PhoneNumbers != nil && len(filter.PhoneNumbers) > 0 {
		query = query.Scopes(whereBySlice(filter.PhoneNumbers, "phone_number"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}
	if filter.IsActive != nil {
		query = query.Scopes(whereBy(*filter.IsActive, "is_active"))
	}

	return query
}

func (r *userRepository) buildExistedFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.Email != nil && *filter.Email != "" {
		query = query.Scopes(orByText(*filter.Email, "email"))
	}
	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Scopes(orByText(*filter.PhoneNumber, "phone_number"))
	}
	if filter.ID != nil && *filter.ID != uuid.Nil {
		query = query.Scopes(orBy(*filter.ID, "id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(orBySlice(filter.IDs, "id"))
	}
	if filter.Emails != nil && len(filter.Emails) > 0 {
		query = query.Scopes(orBySlice(filter.Emails, "email"))
	}
	if filter.PhoneNumbers != nil && len(filter.PhoneNumbers) > 0 {
		query = query.Scopes(orBySlice(filter.PhoneNumbers, "phone_number"))
	}

	return query
}
