package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"fmt"
	"oauth-server/graphql/generated"
	"oauth-server/graphql/model"
)

// User is the resolver for the user field.
func (r *mutationResolver) User(ctx context.Context) (*model.UserMutation, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// OAuth is the resolver for the o_auth field.
func (r *mutationResolver) OAuth(ctx context.Context) (*model.OAuthMutation, error) {
	panic(fmt.Errorf("not implemented: OAuth - o_auth"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model.UserQuery, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// OAuth is the resolver for the o_auth field.
func (r *queryResolver) OAuth(ctx context.Context) (*model.OAuthQuery, error) {
	panic(fmt.Errorf("not implemented: OAuth - o_auth"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*entity.User, error) {
	panic(fmt.Errorf("not implemented: CreateUser - create_user"))
}
func (r *mutationResolver) UpdateUser(ctx context.Context, id uuid.UUID, input model.UpdateUserInput) (*entity.User, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - update_user"))
}
func (r *mutationResolver) UpdatePassword(ctx context.Context, id uuid.UUID, input model.UpdatePasswordInput) (*entity.User, error) {
	panic(fmt.Errorf("not implemented: UpdatePassword - update_password"))
}
func (r *mutationResolver) DeleteUser(ctx context.Context, id uuid.UUID) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteUser - delete_user"))
}
func (r *mutationResolver) ToggleUserActive(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	panic(fmt.Errorf("not implemented: ToggleUserActive - toggle_user_active"))
}
func (r *mutationResolver) CreateOAuth(ctx context.Context, input model.CreateOAuthInput) (*entity.OAuth, error) {
	panic(fmt.Errorf("not implemented: CreateOAuth - create_o_auth"))
}
func (r *mutationResolver) UpdateOAuthStatus(ctx context.Context, id uuid.UUID, status model.OAuthStatus) (*entity.OAuth, error) {
	panic(fmt.Errorf("not implemented: UpdateOAuthStatus - update_o_auth_status"))
}
func (r *mutationResolver) RevokeOAuth(ctx context.Context, id uuid.UUID) (bool, error) {
	panic(fmt.Errorf("not implemented: RevokeOAuth - revoke_o_auth"))
}
func (r *mutationResolver) RevokeAllUserOAuth(ctx context.Context, userID uuid.UUID) (bool, error) {
	panic(fmt.Errorf("not implemented: RevokeAllUserOAuth - revoke_all_user_o_auth"))
}
func (r *queryResolver) Users(ctx context.Context, limit *int32, offset *int32, isActive *bool) ([]*entity.User, error) {
	panic(fmt.Errorf("not implemented: Users - users"))
}
func (r *queryResolver) Me(ctx context.Context) (*entity.User, error) {
	panic(fmt.Errorf("not implemented: Me - me"))
}
func (r *queryResolver) OAuths(ctx context.Context, limit *int32, offset *int32, userID *uuid.UUID, platform *model.OAuthPlatform, status *model.OAuthStatus) ([]*entity.OAuth, error) {
	panic(fmt.Errorf("not implemented: OAuths - o_auths"))
}
func (r *queryResolver) ActiveOAuths(ctx context.Context, userID uuid.UUID) ([]*entity.OAuth, error) {
	panic(fmt.Errorf("not implemented: ActiveOAuths - active_o_auths"))
}
*/
