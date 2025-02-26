package custom_graphql

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var UUIDScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name: "UUID",
	ParseValue: func(value interface{}) interface{} {
		if str, ok := value.(string); ok {
			if u, err := uuid.Parse(str); err == nil {
				return u
			}
		}

		return nil
	},
	Serialize: func(value interface{}) interface{} {
		if u, ok := value.(uuid.UUID); ok {
			return u.String()
		}
		return nil
	},

	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			if u, err := uuid.Parse(valueAST.Value); err == nil {
				return u
			}
		default:
			return nil
		}

		return nil
	},
})
