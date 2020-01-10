package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

type Tweet struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id,omitempty"`
	Message string `json:"message"`
}

func (t *Tweet) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.UserId, validation.Required),
		validation.Field(&t.Message, validation.Required, validation.Length(10, 0)),
	)
}

func (t *Tweet) Sanitize() {
	t.UserId = 0
}
