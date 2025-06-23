package service

import (
	"graphql-demo/internal/safemap"
	"graphql-demo/internal/subscription"
	"slices"
)

type Entity[K comparable] interface {
	GetID() K
}

type Service[K comparable, V Entity[K]] struct {
	Db         *safemap.SafeMap[K, V]
	CreatedSub *subscription.Manager[V]
	DeletedSub *subscription.Manager[K]
}

func (s *Service[K, V]) Get(ids []K) []V {
	values := []V{}

	// default filter function
	filterFunc := func(value V) bool {
		return true
	}

	// filter by ids
	if len(ids) > 0 {
		filterFunc = func(value V) bool {
			return slices.Contains(ids, value.GetID())
		}
	} else if ids != nil {
		filterFunc = func(value V) bool {
			return false
		}
	}

	// all users, but possibly filtered by ids
	for value := range s.Db.Values() {
		if filterFunc(value) {
			values = append(values, value)
		}
	}

	return values
}

func (s *Service[K, V]) Set(value V) {
	s.Db.Set(value.GetID(), value)
	s.CreatedSub.Publish([]V{value})
}

func (s *Service[K, V]) Delete(id K) {
	s.Db.Delete(id)
	s.DeletedSub.Publish([]K{id})
}
