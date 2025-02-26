package user

import (
	custom_graphql "oauth-server/graphql"

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

var userSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: ,
	Mutation: ,
})