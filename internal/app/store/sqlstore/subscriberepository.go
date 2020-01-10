package sqlstore

import "github.com/roman-kapanets/go-twitter-api/internal/app/model"

type SubscribeRepository struct {
	store *Store
}

func (sr *SubscribeRepository) Create(s *model.Subscribe) error {
	if err := s.Validate(); err != nil {
		return err
	}

	return sr.store.db.QueryRow(
		"INSERT INTO subscribes (subscriber, subscribed_to) VALUES ($1, $2) RETURNING id",
		s.Subscriber,
		s.SubscribedTo,
	).Scan(&s.ID)
}
