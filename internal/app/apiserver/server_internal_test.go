package apiserver

import (
	"bytes"
	"encoding/json"
	"github.com/roman-kapanets/go-twitter-api/internal/app/model"
	"github.com/roman-kapanets/go-twitter-api/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestServer_authMiddleware(t *testing.T) {
//	secretKey := "secret"
//	store := teststore.New()
//	u := model.TestUser(t)
//	store.User().Create(u)
//
//	expirationTime := time.Now().Add(5 * time.Minute)
//	claims := &Claims{
//		ID: u.ID,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: expirationTime.Unix(),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(secretKey)
//	assert.NoError(t, err)
//
//	testCases := []struct {
//		name         string
//		token  string
//		expectedCode int
//	}{
//		{
//			name: "authenticated",
//			token: tokenString,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "not authenticated",
//			token:  "some_invalid_token",
//			expectedCode: http.StatusUnauthorized,
//		},
//	}
//
//	s := newServer(store, secretKey)
//
//	mw := s.authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//	}))
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			rec := httptest.NewRecorder()
//			req, _ := http.NewRequest(http.MethodGet, "/", nil)
//			//cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
//			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "token", tc.token))
//			mw.ServeHTTP(rec, req)
//			assert.Equal(t, tc.expectedCode, rec.Code)
//		})
//	}
//}

func TestServer_HandleUsersRegister(t *testing.T) {
	s := newServer(teststore.New(), "secret_key")
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "test@test.com",
				"password": "password",
				"username": "testuser",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]interface{}{
				"email":    "invalid",
				"password": "short",
				"username": "ts",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/register", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUsersLogin(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	store.User().Create(u)
	s := newServer(store, "secret_key")
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]interface{}{
				"email":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]interface{}{
				"email":    u.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/login", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
