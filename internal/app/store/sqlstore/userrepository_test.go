package sqlstore_test

import (
	"github.com/roman-kapanets/go-twitter-api/internal/app/model"
	"github.com/roman-kapanets/go-twitter-api/internal/app/store"
	"github.com/roman-kapanets/go-twitter-api/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	email := "test@test.com"

	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = "test@test.com"
	u.Username = "test"

	s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	username := "testusername"

	_, err := s.User().FindByUsername(username)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u1 := model.TestUser(t)
	u1.Email = "test@test.com"
	u1.Username = "test"

	s.User().Create(u1)

	u2, err := s.User().FindByUsername(u1.Username)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	u1 := model.TestUser(t)
	s.User().Create(u1)

	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
