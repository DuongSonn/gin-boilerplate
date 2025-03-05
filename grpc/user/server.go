package user

import (
	context "context"
	"io"
	"log/slog"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	"oauth-server/package/errors"
	logger "oauth-server/package/log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceServer struct {
	postgresRepo repository.RepositoryCollections
	helpers      helper.HelperCollections

	UnimplementedUserServiceServer
}

func NewUserServiceServer(
	postgresRepo repository.RepositoryCollections,
	helpers helper.HelperCollections,
) UserServiceServer {
	return &userServiceServer{
		postgresRepo: postgresRepo,
		helpers:      helpers,
	}
}

func (s *userServiceServer) GetUser(ctx context.Context, data *GetUserRequest) (*GetUserResponse, error) {
	filter, err := s.buildFilter(data)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.GetMessage(errors.ErrCodeInternalServerError))
	}

	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, filter)
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

	users, err := s.postgresRepo.UserRepo.FindManyByFilter(ctx, nil, filter)
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
				"Send GetUsers",
				slog.String("ID", user.GetUserID()),
				slog.String("Error", err.Error()),
			)

			return status.Error(codes.Internal, errors.GetMessage(errors.ErrCodeInternalServerError))
		}
	}

	return nil
}

func (s *userServiceServer) CreateUsers(stream UserService_CreateUsersServer) error {
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&CreateUserResponse{Success: true})
		}
		if err != nil {
			logger.GetLogger().Info(
				"Recv CreateUsers",
				slog.String("Error", err.Error()),
			)

			return status.Error(codes.Internal, errors.GetMessage(errors.ErrCodeInternalServerError))
		}
		if _, err := s.helpers.UserHelper.CreateUser(stream.Context(), &model.RegisterRequest{
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Password:    user.GetPassword(),
		}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

	}
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

			return nil, err
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
