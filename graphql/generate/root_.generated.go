// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generate

import (
	"bytes"
	"context"
	"errors"
	"oauth-server/graphql/model"
	"sync/atomic"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	return &executableSchema{
		schema:     cfg.Schema,
		resolvers:  cfg.Resolvers,
		directives: cfg.Directives,
		complexity: cfg.Complexity,
	}
}

type Config struct {
	Schema     *ast.Schema
	Resolvers  ResolverRoot
	Directives DirectiveRoot
	Complexity ComplexityRoot
}

type ResolverRoot interface {
	Mutation() MutationResolver
	Query() QueryResolver
}

type DirectiveRoot struct {
}

type ComplexityRoot struct {
	Mutation struct {
		CreateUser func(childComplexity int, input model.NewUser) int
	}

	OAuth struct {
		ExpireAt func(childComplexity int) int
		ID       func(childComplexity int) int
		IP       func(childComplexity int) int
		Platform func(childComplexity int) int
		Status   func(childComplexity int) int
		Token    func(childComplexity int) int
	}

	Query struct {
		Users func(childComplexity int) int
	}

	User struct {
		Email       func(childComplexity int) int
		ID          func(childComplexity int) int
		IsActive    func(childComplexity int) int
		OAuth       func(childComplexity int) int
		PhoneNumber func(childComplexity int) int
	}
}

type executableSchema struct {
	schema     *ast.Schema
	resolvers  ResolverRoot
	directives DirectiveRoot
	complexity ComplexityRoot
}

func (e *executableSchema) Schema() *ast.Schema {
	if e.schema != nil {
		return e.schema
	}
	return parsedSchema
}

func (e *executableSchema) Complexity(typeName, field string, childComplexity int, rawArgs map[string]any) (int, bool) {
	ec := executionContext{nil, e, 0, 0, nil}
	_ = ec
	switch typeName + "." + field {

	case "Mutation.createUser":
		if e.complexity.Mutation.CreateUser == nil {
			break
		}

		args, err := ec.field_Mutation_createUser_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreateUser(childComplexity, args["input"].(model.NewUser)), true

	case "OAuth.expireAt":
		if e.complexity.OAuth.ExpireAt == nil {
			break
		}

		return e.complexity.OAuth.ExpireAt(childComplexity), true

	case "OAuth.id":
		if e.complexity.OAuth.ID == nil {
			break
		}

		return e.complexity.OAuth.ID(childComplexity), true

	case "OAuth.ip":
		if e.complexity.OAuth.IP == nil {
			break
		}

		return e.complexity.OAuth.IP(childComplexity), true

	case "OAuth.platform":
		if e.complexity.OAuth.Platform == nil {
			break
		}

		return e.complexity.OAuth.Platform(childComplexity), true

	case "OAuth.status":
		if e.complexity.OAuth.Status == nil {
			break
		}

		return e.complexity.OAuth.Status(childComplexity), true

	case "OAuth.token":
		if e.complexity.OAuth.Token == nil {
			break
		}

		return e.complexity.OAuth.Token(childComplexity), true

	case "Query.users":
		if e.complexity.Query.Users == nil {
			break
		}

		return e.complexity.Query.Users(childComplexity), true

	case "User.email":
		if e.complexity.User.Email == nil {
			break
		}

		return e.complexity.User.Email(childComplexity), true

	case "User.id":
		if e.complexity.User.ID == nil {
			break
		}

		return e.complexity.User.ID(childComplexity), true

	case "User.isActive":
		if e.complexity.User.IsActive == nil {
			break
		}

		return e.complexity.User.IsActive(childComplexity), true

	case "User.oAuth":
		if e.complexity.User.OAuth == nil {
			break
		}

		return e.complexity.User.OAuth(childComplexity), true

	case "User.phoneNumber":
		if e.complexity.User.PhoneNumber == nil {
			break
		}

		return e.complexity.User.PhoneNumber(childComplexity), true

	}
	return 0, false
}

func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	opCtx := graphql.GetOperationContext(ctx)
	ec := executionContext{opCtx, e, 0, 0, make(chan graphql.DeferredResult)}
	inputUnmarshalMap := graphql.BuildUnmarshalerMap(
		ec.unmarshalInputNewUser,
	)
	first := true

	switch opCtx.Operation.Operation {
	case ast.Query:
		return func(ctx context.Context) *graphql.Response {
			var response graphql.Response
			var data graphql.Marshaler
			if first {
				first = false
				ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
				data = ec._Query(ctx, opCtx.Operation.SelectionSet)
			} else {
				if atomic.LoadInt32(&ec.pendingDeferred) > 0 {
					result := <-ec.deferredResults
					atomic.AddInt32(&ec.pendingDeferred, -1)
					data = result.Result
					response.Path = result.Path
					response.Label = result.Label
					response.Errors = result.Errors
				} else {
					return nil
				}
			}
			var buf bytes.Buffer
			data.MarshalGQL(&buf)
			response.Data = buf.Bytes()
			if atomic.LoadInt32(&ec.deferred) > 0 {
				hasNext := atomic.LoadInt32(&ec.pendingDeferred) > 0
				response.HasNext = &hasNext
			}

			return &response
		}
	case ast.Mutation:
		return func(ctx context.Context) *graphql.Response {
			if !first {
				return nil
			}
			first = false
			ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
			data := ec._Mutation(ctx, opCtx.Operation.SelectionSet)
			var buf bytes.Buffer
			data.MarshalGQL(&buf)

			return &graphql.Response{
				Data: buf.Bytes(),
			}
		}

	default:
		return graphql.OneShot(graphql.ErrorResponse(ctx, "unsupported GraphQL operation"))
	}
}

type executionContext struct {
	*graphql.OperationContext
	*executableSchema
	deferred        int32
	pendingDeferred int32
	deferredResults chan graphql.DeferredResult
}

func (ec *executionContext) processDeferredGroup(dg graphql.DeferredGroup) {
	atomic.AddInt32(&ec.pendingDeferred, 1)
	go func() {
		ctx := graphql.WithFreshResponseContext(dg.Context)
		dg.FieldSet.Dispatch(ctx)
		ds := graphql.DeferredResult{
			Path:   dg.Path,
			Label:  dg.Label,
			Result: dg.FieldSet,
			Errors: graphql.GetErrors(ctx),
		}
		// null fields should bubble up
		if dg.FieldSet.Invalids > 0 {
			ds.Result = graphql.Null
		}
		ec.deferredResults <- ds
	}()
}

func (ec *executionContext) introspectSchema() (*introspection.Schema, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapSchema(ec.Schema()), nil
}

func (ec *executionContext) introspectType(name string) (*introspection.Type, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapTypeFromDef(ec.Schema(), ec.Schema().Types[name]), nil
}

var sources = []*ast.Source{
	{Name: "../schema/user.graphqls", Input: `scalar UUID
scalar Int64

type User {
    id: UUID!
    phoneNumber: String
    email: String
    isActive: Boolean!
    oAuth: OAuth
}

enum OAuthPlatform {
    WEB
    MOBILE
}
enum OAuthStatus {
    ACTIVE
    INACTIVE
    BLOCKED
}
type OAuth {
    id: UUID!
    ip: String!
    platform: OAuthPlatform!
    token: String!
    status: OAuthStatus!
    expireAt: Int64!
}

type Query {
    users: [User!]!
}

input NewUser {
    password: String!
    confirmPassword: String!
    email: String
    phoneNumber: String
}
type Mutation {
    createUser(input: NewUser!): User!
}
`, BuiltIn: false},
}
var parsedSchema = gqlparser.MustLoadSchema(sources...)
