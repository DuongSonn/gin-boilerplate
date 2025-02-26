package user

import (
	"oauth-server/app/model"
	"oauth-server/app/service"
	"oauth-server/package/errors"
	"oauth-server/package/validator"

	"github.com/graphql-go/graphql"
)

type MutationResolver struct {
	UserService service.UserService
}

func NewMutationResolver(userService service.UserService) *MutationResolver {
	return &MutationResolver{UserService: userService}
}

func (r *MutationResolver) CreateUser(params graphql.ResolveParams) (interface{}, error) {
	phone, phoneOk := params.Args["phone"].(string)
	email, emailOk := params.Args["email"].(string)
	password, _ := params.Args["password"].(string)
	confirmPassword, _ := params.Args["confirm_password"].(string)

	if !phoneOk && !emailOk {
		return nil, errors.NewCustomError(
			errors.ErrCodeValidatorRequired,
			errors.GetCustomMessage(errors.ErrCodeValidatorRequired, "SĐT/Email"),
		)
	}
	if phoneOk && !validator.IsValidPhoneNumber(phone) {
		return nil, errors.NewCustomError(
			errors.ErrCodeValidatorFormat,
			errors.GetCustomMessage(errors.ErrCodeValidatorFormat, "SĐT"),
		)
	}
	if emailOk && !validator.IsValidEmail(email) {
		return nil, errors.NewCustomError(
			errors.ErrCodeValidatorFormat,
			errors.GetCustomMessage(errors.ErrCodeValidatorFormat, "Email"),
		)
	}
	if password != confirmPassword {
		return nil, errors.NewCustomError(
			errors.ErrCodeValidatorVerifiedData,
			errors.GetCustomMessage(errors.ErrCodeValidatorVerifiedData, "Mật khẩu"),
		)
	}

	return r.UserService.Register(params.Context, &model.RegisterRequest{
		Email:           &email,
		PhoneNumber:     &phone,
		Password:        password,
		ConfirmPassword: confirmPassword,
	})
}

func newUserMutation(resolver *MutationResolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "UserMutation",
		Fields: graphql.Fields{
			"create_user": &graphql.Field{
				Type: graphql.NewNonNull(userType),
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"phone_number": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"confirm_password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolver.CreateUser,
			},
		},
	})
}
