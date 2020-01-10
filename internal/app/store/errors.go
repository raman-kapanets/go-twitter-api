package store

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateEntity = errors.New("duplicate entity")
)
