package service

import (
	"graphql-demo/internal/graph/model"
	"graphql-demo/internal/safemap"
	"graphql-demo/internal/subscription"
)

type PostService = Service[int, *model.Post]

func NewPostService() *PostService {
	return &PostService{
		Db:         safemap.New[int, *model.Post](),
		CreatedSub: subscription.NewManager[*model.Post](),
		DeletedSub: subscription.NewManager[int](),
	}
}
