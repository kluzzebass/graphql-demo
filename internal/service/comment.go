package service

import (
	"graphql-demo/internal/graph/model"
	"graphql-demo/internal/safemap"
	"graphql-demo/internal/subscription"
)

type CommentService = Service[int, *model.Comment]

func NewCommentService() *CommentService {
	return &CommentService{
		Db:         safemap.New[int, *model.Comment](),
		CreatedSub: subscription.NewManager[*model.Comment](),
		DeletedSub: subscription.NewManager[int](),
	}
}
