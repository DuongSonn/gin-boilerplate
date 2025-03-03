package user

import (
	"github.com/graphql-go/graphql"
)

func (r *userResolver) Mutation() *graphql.Object {
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
				Resolve: r.CreateUser,
			},
		},
	})
}
