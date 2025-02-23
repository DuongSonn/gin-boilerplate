package helper

import (
	"context"
	"log/slog"
	"oauth-server/app/entity"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	"oauth-server/package/errors"
	logger "oauth-server/package/log"
)

type userHelper struct {
	postgresRepo repository.RepositoryCollections
}

func newUserHelper(
	postgresRepo repository.RepositoryCollections,
) UserHelper {
	return &userHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *userHelper) CreateUser(ctx context.Context, data *model.RegisterRequest) error {
	// Check user exited
	existedUser, err := h.postgresRepo.UserRepo.FindManyByFilter(ctx, nil, &repository.FindUserByFilter{
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
	})
	if err != nil {
		logger.GetLogger().Info(
			"FindManyByFilter",
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
	user := entity.NewUser()
	user.PhoneNumber = data.PhoneNumber
	user.Email = data.Email
	user.Password = data.Password
	if err := h.postgresRepo.UserRepo.Create(ctx, nil, user); err != nil {
		logger.GetLogger().Info(
			"Create",
			slog.Group(
				(entity.USER_TABLE_NAME),
				slog.String("email", *data.Email),
				slog.String("phone_number", *data.PhoneNumber),
			),
			slog.String("error", err.Error()),
		)

		return errors.New(errors.ErrCodeInternalServerError)
	}

	return nil
}
