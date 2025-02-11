package user

import (
	context "context"
	"log/slog"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/package/errors"
	logger "oauth-server/package/log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceServer struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections

	UnimplementedUserServiceServer
}

func NewUserServiceServer(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserServiceServer {
	return &userServiceServer{
		postgresRepo: postgresRepo,
	}
}

func (s *userServiceServer) GetUser(ctx context.Context, data *GetUserRequest) (*GetUserResponse, error) {
	filter, err := s.buildFilter(data)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.GetMessage(errors.ErrCodeInternalServerError))
	}

	user, err := s.postgresRepo.PostgresUserRepo.FindUserByFilter(ctx, nil, filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.GetMessage(errors.ErrCodeUserNotFound))
	}

	return &GetUserResponse{
		Success: true,
		User: &User{
			Id:          user.GetUserID(),
			PhoneNumber: user.GetPhoneNumber(),
			Email:       user.GetEmail(),
			IsActive:    user.IsActive,
		},
	}, nil
}

func (s *userServiceServer) GetUsers(data *GetUserRequest, stream UserService_GetUsersServer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter, err := s.buildFilter(data)
	if err != nil {
		return status.Error(codes.Internal, errors.GetMessage(errors.ErrCodeInternalServerError))
	}

	users, err := s.postgresRepo.PostgresUserRepo.FindUsersByFilter(ctx, nil, filter)
	if err != nil {
		return status.Error(codes.NotFound, errors.GetMessage(errors.ErrCodeUserNotFound))
	}

	for _, user := range users {
		if err := stream.Send(&GetUserResponse{
			Success: true,
			User: &User{
				Id:          user.ID.String(),
				PhoneNumber: user.GetPhoneNumber(),
				Email:       user.GetEmail(),
				IsActive:    user.IsActive,
			},
		}); err != nil {
			logger.GetLogger().Info(
				"Send UserInfo",
				slog.String("ID", user.GetUserID()),
				slog.String("Error", err.Error()),
			)

			return status.Error(codes.Internal, errors.GetMessage(errors.ErrCodeInternalServerError))
		}
	}

	return nil
}

func (s *userServiceServer) CreateUsers(stream UserService_CreateUsersServer) error {
	return nil
}

// ---------------------------------------------------------------------------

func (s *userServiceServer) buildFilter(data *GetUserRequest) (*repository.FindUserByFilter, error) {
	filter := &repository.FindUserByFilter{}

	if data.Id != nil && data.GetId() != "" {
		id, err := uuid.Parse(data.GetId())
		if err != nil {
			logger.GetLogger().Info(
				"uuid.Parse",
				slog.String("ID", data.GetId()),
				slog.String("Error", err.Error()),
			)
		}

		filter.ID = &id
	}
	if data.PhoneNumber != nil && data.GetPhoneNumber() != "" {
		filter.PhoneNumber = data.PhoneNumber
	}
	if data.Email != nil && data.GetEmail() != "" {
		filter.Email = data.Email
	}

	return filter, nil
}
