package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"fmt"
	"oauth-server/app/entity"
	"oauth-server/graphql/generated"
	"oauth-server/graphql/model"
	"time"
)

// Platform is the resolver for the platform field.
func (r *oAuthResolver) Platform(ctx context.Context, obj *entity.OAuth) (model.OAuthPlatform, error) {
	panic(fmt.Errorf("not implemented: Platform - platform"))
}

// Status is the resolver for the status field.
func (r *oAuthResolver) Status(ctx context.Context, obj *entity.OAuth) (model.OAuthStatus, error) {
	panic(fmt.Errorf("not implemented: Status - status"))
}

// ExpireAt is the resolver for the expire_at field.
func (r *oAuthResolver) ExpireAt(ctx context.Context, obj *entity.OAuth) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: ExpireAt - expire_at"))
}

// CreatedAt is the resolver for the created_at field.
func (r *oAuthResolver) CreatedAt(ctx context.Context, obj *entity.OAuth) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: CreatedAt - created_at"))
}

// UpdatedAt is the resolver for the updated_at field.
func (r *oAuthResolver) UpdatedAt(ctx context.Context, obj *entity.OAuth) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: UpdatedAt - updated_at"))
}

// LoginAt is the resolver for the login_at field.
func (r *oAuthResolver) LoginAt(ctx context.Context, obj *entity.OAuth) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: LoginAt - login_at"))
}

// User is the resolver for the user field.
func (r *oAuthResolver) User(ctx context.Context, obj *entity.OAuth) (*entity.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// OAuth returns generated.OAuthResolver implementation.
func (r *Resolver) OAuth() generated.OAuthResolver { return &oAuthResolver{r} }

type oAuthResolver struct{ *Resolver }
