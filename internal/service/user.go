package service

import (
	"graphql-demo/internal/graph/model"
	"graphql-demo/internal/safemap"
	"graphql-demo/internal/subscription"
)

type UserService = Service[int, *model.User]

func NewUserService() *UserService {
	return &UserService{
		Db:         safemap.New[int, *model.User](),
		CreatedSub: subscription.NewManager[*model.User](),
		DeletedSub: subscription.NewManager[int](),
	}
}
