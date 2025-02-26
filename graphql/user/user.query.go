package user

import "github.com/graphql-go/graphql"

var userQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserQuery",
	Fields: graphql.Fields{
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
	},
})
