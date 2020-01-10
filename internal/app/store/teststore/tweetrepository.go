package teststore

import "github.com/roman-kapanets/go-twitter-api/internal/app/model"

type TweetRepository struct {
	store  *Store
	tweets map[int]*model.Tweet
}

func (t *TweetRepository) Create(tweet *model.Tweet) error {
	if err := tweet.Validate(); err != nil {
		return err
	}
	tweet.ID = len(t.tweets) + 1
	t.tweets[tweet.ID] = tweet
	tweet.ID = len(t.tweets)
	return nil
}

func (t *TweetRepository) FindBySubscription(userID int) (tweets []*model.Tweet, err error) {
	return tweets, nil
}
