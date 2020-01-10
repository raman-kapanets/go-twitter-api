package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/roman-kapanets/go-twitter-api/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	tweetRepository *TweetRepository
	subscribeRepository *SubscribeRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) User() store.UserRepository {

	if store.userRepository != nil {
		return store.userRepository
	}

	store.userRepository = &UserRepository{
		store: store,
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