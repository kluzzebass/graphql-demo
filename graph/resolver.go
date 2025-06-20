package graph

import (
	"context"
	"log"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/kluzzebass/graphql-demo/graph/model"
	"github.com/kluzzebass/graphql-demo/safemap"
	"github.com/kluzzebass/graphql-demo/subscription"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserMap        *safemap.SafeMap[int, *model.User]
	PostMap        *safemap.SafeMap[int, *model.Post]
	PostCreatedSub *subscription.Manager[*model.Post]
	PostDeletedSub *subscription.Manager[int]
}

// LogResolverDepth is a helper function to get and log the current resolver's depth.
func LogResolverDepth(ctx context.Context, resolverName string) {
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

	log.Printf("%s%s → %s\n", spaces, resolverName, path.String())

	// You can still access other useful info directly from the current FieldContext:
	// fmt.Printf("  Field Name: %s\n", fc.Field.Name)
	// if fc.Index != nil {
	// 	fmt.Printf("  Index (if list item): %d\n", *fc.Index)
	// }
}
