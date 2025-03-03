package custom_graphql

import (
	"github.com/graphql-go/graphql"
)

func NewSchema(resolvers []interface{}) (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    mergeQueries(resolvers),
		Mutation: mergeMutations(resolvers),
	})
}

func mergeQueries(resolvers []interface{}) *graphql.Object {
	fields := graphql.Fields{}
	for _, resolver := range resolvers {
		if r, ok := resolver.(Resolver); ok {
			query := r.Query()
			for name, field := range query.Fields() {
				fields[name] = &graphql.Field{
					Type: field.Type,
					// Args:        field.Args,
					Resolve:   field.Resolve,
					Subscribe: field.Subscribe,
					// Deprecation: field.Deprecation,
					Description: field.Description,
				}
			}
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: fields,
	})
}

func mergeMutations(resolvers []interface{}) *graphql.Object {
	fields := graphql.Fields{}
	for _, resolver := range resolvers {
		if r, ok := resolver.(Resolver); ok {
			mutation := r.Mutation()
			for name, field := range mutation.Fields() {
				fields[name] = &graphql.Field{
					Type: field.Type,
					// Args:        field.Args,
					Resolve:   field.Resolve,
					Subscribe: field.Subscribe,
					// Deprecation: field.Deprecation,
					Description: field.Description,
				}
			}
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: fields,
	})
}
