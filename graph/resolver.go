package graph

import (
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
