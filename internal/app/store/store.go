package store

type Store interface {
	User() UserRepository
	Tweet() TweetRepository
	Subscribe() SubscribeRepository
}
