package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"fmt"
	"oauth-server/graphql/generate"
	"oauth-server/graphql/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented: Users - users"))
}

// Mutation returns generate.MutationResolver implementation.
func (r *Resolver) Mutation() generate.MutationResolver { return &mutationResolver{r} }

// Query returns generate.QueryResolver implementation.
func (r *Resolver) Query() generate.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
