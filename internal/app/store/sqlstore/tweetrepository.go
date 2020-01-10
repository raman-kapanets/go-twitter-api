package sqlstore

import "github.com/roman-kapanets/go-twitter-api/internal/app/model"

type TweetRepository struct {
	store *Store
}

func (t *TweetRepository) Create(tweet *model.Tweet) error {
	if err := tweet.Validate(); err != nil {
		return err
	}

	return t.store.db.QueryRow(
		"INSERT INTO tweets (user_id, message) VALUES ($1, $2) RETURNING id",
		tweet.UserId,
		tweet.Message,
	).Scan(&tweet.ID)
}

func (t *TweetRepository) FindBySubscription(userID int) (tweets []*model.Tweet, err error) {
	rows, err := t.store.db.Query("SELECT tw.id, tw.message FROM tweets tw INNER JOIN subscribes sb ON tw.user_id = sb.subscribed_to WHERE sb.subscriber = $1;", userID)
	if err != nil {
		return tweets, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var tweet model.Tweet
		err = rows.Scan(&tweet.ID, &tweet.Message)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}

	return tweets, nil
}
