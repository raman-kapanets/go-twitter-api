package store

import (
	"github.com/roman-kapanets/go-twitter-api/internal/app/model"
)

type UserRepository interface {
	Create(user *model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	FindByUsername(string) (*model.User, error)
}

type TweetRepository interface {
	Create(tweet *model.Tweet) error
	FindBySubscription(int) ([]*model.Tweet, error)
}

type SubscribeRepository interface {
	Create(*model.Subscribe) error
}