package teststore

import (
	"github.com/roman-kapanets/go-twitter-api/internal/app/model"
	"github.com/roman-kapanets/go-twitter-api/internal/app/store"
)

type Store struct {
	userRepository      *UserRepository
	tweetRepository     *TweetRepository
	subscribeRepository *SubscribeRepository
}

func New() *Store {
	return &Store{}
}

func (store *Store) User() store.UserRepository {

	if store.userRepository != nil {
		return store.userRepository
	}

	store.userRepository = &UserRepository{
		store: store,
		users: make(map[int]*model.User),
	}

	return store.userRepository
}

func (store *Store) Tweet() store.TweetRepository {

	if store.tweetRepository != nil {
		return store.tweetRepository
	}

	store.tweetRepository = &TweetRepository{
		store: store,
	}

	return store.tweetRepository
}

func (store *Store) Subscribe() store.SubscribeRepository {

	if store.subscribeRepository != nil {
		return store.subscribeRepository
	}

	store.subscribeRepository = &SubscribeRepository{
		store: store,
	}

	return store.subscribeRepository
}
