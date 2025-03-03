package user

import (
	"oauth-server/app/model"
	"oauth-server/app/service"
	custom_graphql "oauth-server/graphql"
	"oauth-server/package/errors"
	"oauth-server/package/validator"

	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(custom_graphql.UUIDScalar),
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"phone_number": &graphql.Field{
			Type: graphql.String,
		},
		"is_active": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
		"created_at": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"updated_at": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
})

type userResolver struct {
	services service.ServiceCollections
}

func NewUserResolver(services service.ServiceCollections) custom_graphql.UserResolver {
	return &userResolver{services: services}
}

func (r *userResolver) CreateUser(params graphql.ResolveParams) (interface{}, error) {
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

	return r.services.UserSvc.Register(params.Context, &model.RegisterRequest{
		Email:           &email,
		PhoneNumber:     &phone,
		Password:        password,
		ConfirmPassword: confirmPassword,
	})
}
