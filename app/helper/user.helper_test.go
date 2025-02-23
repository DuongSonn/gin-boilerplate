package helper

import (
	"context"
	"io"
	"log/slog"
	"oauth-server/app/entity"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	mocks "oauth-server/mocks/repository"
	"oauth-server/package/errors"
	logger "oauth-server/package/log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

var (
	testEmail       = "test@gmail.com"
	testPhoneNumber = "0987654321"
	testUser        = entity.User{
		Email:       &testEmail,
		PhoneNumber: &testPhoneNumber,
	}
)

func TestMain(m *testing.M) {
	logger.SetLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

func Test_newUserHelper(t *testing.T) {
	type args struct {
		postgresRepo repository.RepositoryCollections
	}
	tests := []struct {
		name string
		args args
		want UserHelper
	}{
		{
			name: "Init Success",
			args: args{
				postgresRepo: mocks.InitMockRepository(t),
			},
			want: &userHelper{
				postgresRepo: mocks.InitMockRepository(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUserHelper(tt.args.postgresRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUserHelper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userHelper_CreateUser(t *testing.T) {
	type fields struct {
		postgresRepo repository.RepositoryCollections
	}
	type args struct {
		ctx  context.Context
		data *model.RegisterRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "User FindManyByFilter Error",
			fields: fields{
				postgresRepo: mocks.InitMockRepository(t),
			},
			args: args{
				ctx: context.Background(),
				data: &model.RegisterRequest{
					Email:       &testEmail,
					PhoneNumber: &testPhoneNumber,
					Password:    "password",
				},
			},
			wantErr: true,
			err:     errors.New(errors.ErrCodeInternalServerError),
		},
		{
			name: "User Existed",
			fields: fields{
				postgresRepo: mocks.InitMockRepository(t),
			},
			args: args{
				ctx: context.Background(),
				data: &model.RegisterRequest{
					Email:       &testEmail,
					PhoneNumber: &testPhoneNumber,
					Password:    "password",
				},
			},
			wantErr: true,
			err:     errors.New(errors.ErrCodeUserExisted),
		},
		{
			name: "User Create Error",
			fields: fields{
				postgresRepo: mocks.InitMockRepository(t),
			},
			args: args{
				ctx: context.Background(),
				data: &model.RegisterRequest{
					Email:       &testEmail,
					PhoneNumber: &testPhoneNumber,
					Password:    "password",
				},
			},
			wantErr: true,
			err:     errors.New(errors.ErrCodeInternalServerError),
		},
		{
			name: "User Create Success",
			fields: fields{
				postgresRepo: mocks.InitMockRepository(t),
			},
			args: args{
				ctx: context.Background(),
				data: &model.RegisterRequest{
					Email:       &testEmail,
					PhoneNumber: &testPhoneNumber,
					Password:    "password",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo, _ := tt.fields.postgresRepo.UserRepo.(*mocks.UserRepository)

			h := &userHelper{
				postgresRepo: tt.fields.postgresRepo,
			}

			switch tt.name {
			case "User FindManyByFilter Error":
				mockUserRepo.On("FindManyByFilter", tt.args.ctx, mock.Anything, mock.Anything).Return(nil, errors.New(errors.ErrCodeInternalServerError)).Once()
			case "User Existed":
				mockUserRepo.On("FindManyByFilter", tt.args.ctx, mock.Anything, mock.Anything).Return([]entity.User{
					testUser,
				}, nil)
			case "User Create Error":
				mockUserRepo.On("FindManyByFilter", tt.args.ctx, mock.Anything, mock.Anything).Return([]entity.User{}, nil).Once()
				mockUserRepo.On("Create", tt.args.ctx, mock.Anything, mock.Anything).Return(errors.New(errors.ErrCodeInternalServerError)).Once()
			case "User Create Success":
				mockUserRepo.On("FindManyByFilter", tt.args.ctx, mock.Anything, mock.Anything).Return([]entity.User{}, nil).Once()
				mockUserRepo.On("Create", tt.args.ctx, mock.Anything, mock.Anything).Return(nil).Once()
			}

			err := h.CreateUser(tt.args.ctx, tt.args.data)
			if !tt.wantErr != (err == nil) {
				t.Errorf("userHelper.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("userHelper.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && (err.Error() != tt.err.Error()) {
				t.Errorf("userHelper.CreateUser() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
