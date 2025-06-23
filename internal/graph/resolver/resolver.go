package resolver

import (
	"context"
	"fmt"
	"graphql-demo/internal/service"
	"log"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService    *service.UserService
	PostService    *service.PostService
	CommentService *service.CommentService
}

type ResolverArg struct {
	Key   string
	Value any
}

type ResolverArgs []ResolverArg

func (r ResolverArgs) String() string {
	parts := []string{}
	for _, arg := range r {
		if arg.Value != nil {
			value := reflect.ValueOf(arg.Value)
			if value.Kind() == reflect.Ptr && value.IsNil() {
				continue
			}
			if value.Kind() == reflect.Slice && value.IsNil() {
				continue
			}

			// Handle pointer vs non-pointer values
			var interfaceValue interface{}
			if value.Kind() == reflect.Ptr {
				interfaceValue = value.Elem().Interface()
			} else {
				interfaceValue = value.Interface()
			}

			parts = append(parts, fmt.Sprintf("%s: %v", arg.Key, interfaceValue))
		}
	}
	if len(parts) > 0 {
		return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
	}
	return ""
}

// LogResolverDepth is a helper function to get and log the current resolver's depth.
func LogResolverDepth(ctx context.Context, resolverName string, args ResolverArgs) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		log.Printf("WARNING: No FieldContext found for resolver '%s'. This might happen for root resolvers or in non-GraphQL contexts.\n", resolverName)
		return
	}

	// Calculate depth by traversing the FieldContext's Parent chain
	depth := -2
	currentFc := fc // Start from the current FieldContext
	path := currentFc.Path()
	for currentFc != nil {
		depth++
		currentFc = currentFc.Parent // This correctly traverses up the FieldContext chain
	}

	if depth < 0 {
		depth = 0
	}

	spaces := strings.Repeat(" ", depth*2)

	log.Printf("%s%s%s → %s\n", spaces, resolverName, args.String(), path.String())
}
