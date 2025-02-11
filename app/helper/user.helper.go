package helper

import (
	"context"
	"log/slog"
	"oauth-server/app/entity"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/package/database"
	"oauth-server/package/errors"
	logger "oauth-server/package/log"
)

type userHelper struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewUserHelper(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserHelper {
	return &userHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *userHelper) CreateUser(ctx context.Context, data *model.RegisterRequest) error {
	// Check user exited
	existedUser, err := h.postgresRepo.PostgresUserRepo.FindUsersByFilter(ctx, nil, &repository.FindUserByFilter{
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
	})
	if err != nil {
		logger.GetLogger().Info(
			"FindUsersByFilter",
			slog.Group(
				(entity.USER_TABLE_NAME),
				slog.String("email", *data.Email),
				slog.String("phone_number", *data.PhoneNumber),
			),
			slog.String("error", err.Error()),
		)
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedUser) > 0 {
		return errors.New(errors.ErrCodeUserExisted)
	}

	// Create user
	tx := database.BeginPostgresTransaction()
	user := entity.NewUser()
	user.PhoneNumber = data.PhoneNumber
	user.Email = data.Email
	user.Password = data.Password
	if err := h.postgresRepo.PostgresUserRepo.CreateUser(ctx, tx, user); err != nil {
		logger.GetLogger().Info(
			"CreateUser",
			slog.Group(
				(entity.USER_TABLE_NAME),
				slog.String("email", *data.Email),
				slog.String("phone_number", *data.PhoneNumber),
			),
			slog.String("error", err.Error()),
		)

		tx.WithContext(ctx).Rollback()
		return errors.New(errors.ErrCodeInternalServerError)
	}
	tx.WithContext(ctx).Commit()

	return nil
}
