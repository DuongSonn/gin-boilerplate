package custom_graphql

import "github.com/graphql-go/graphql"

type Resolver interface {
	Mutation() *graphql.Object
	Query() *graphql.Object
}

type UserResolver interface {
	Resolver
	CreateUser(p graphql.ResolveParams) (interface{}, error)
}
