package teststore

import "github.com/roman-kapanets/go-twitter-api/internal/app/model"

type SubscribeRepository struct {
	store *Store
	subscribes map[int]*model.Subscribe
}

func (sr *SubscribeRepository) Create(s *model.Subscribe) error  {
	if err := s.Validate(); err != nil {
		return err
	}
	s.ID = len(sr.subscribes) + 1
	sr.subscribes[s.ID] = s
	s.ID = len(sr.subscribes)
	return nil
}