package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

type Subscribe struct {
	ID           int `json:"id"`
	Subscriber   int `json:"subscriber,omitempty"`
	SubscribedTo int `json:"subscribed_to,omitempty"`
}

func (s *Subscribe) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Subscriber, validation.Required),
		validation.Field(&s.SubscribedTo, validation.Required),
	)
}

func (s *Subscribe) Sanitize() {
	s.Subscriber = 0
	s.SubscribedTo = 0
}