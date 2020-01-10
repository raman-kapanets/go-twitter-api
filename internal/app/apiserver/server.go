package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/roman-kapanets/go-twitter-api/internal/app/model"
	"github.com/roman-kapanets/go-twitter-api/internal/app/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	ctxKeyUser ctxKey = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errIncorrectUsername        = errors.New("incorrect nickname")
	errSubscribeToSelf          = errors.New("you can't subscribe to yourself")
)

type ctxKey int8

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
	appKey []byte
}

type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func newServer(store store.Store, appKey string) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
		appKey: []byte(appKey),
	}

	s.ConfigureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ConfigureRouter() {
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/register", s.handleUsersRegister()).Methods(http.MethodPost)
	s.router.HandleFunc("/login", s.handleUsersLogin()).Methods(http.MethodPost)

	s.router.HandleFunc("/subscribe", s.authMiddleware(s.handleSubscribe())).Methods(http.MethodPost)

	s.router.HandleFunc("/tweets", s.authMiddleware(s.handleTweetsCreate())).Methods(http.MethodPost)
	s.router.HandleFunc("/tweets", s.authMiddleware(s.handleTweetsFetch())).Methods(http.MethodGet)
}

func (s *server) handleSubscribe() http.HandlerFunc {
	type request struct {
		Nickname string `json:"nickname"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByUsername(req.Nickname)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errIncorrectUsername)
			return
		}

		currentUser := r.Context().Value(ctxKeyUser).(*model.User)

		if u.ID == currentUser.ID {
			s.error(w, r, http.StatusUnprocessableEntity, errSubscribeToSelf)
			return
		}

		subscribe := &model.Subscribe{
			Subscriber:   currentUser.ID,
			SubscribedTo: u.ID,
		}

		if err := s.store.Subscribe().Create(subscribe); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		subscribe.Sanitize()

		s.respond(w, r, http.StatusOK, subscribe)
	}
}

func (s *server) handleTweetsCreate() http.HandlerFunc {
	type request struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		currentUser := r.Context().Value(ctxKeyUser).(*model.User)

		tweet := &model.Tweet{
			UserId:  currentUser.ID,
			Message: req.Message,
		}

		if err := s.store.Tweet().Create(tweet); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		tweet.Sanitize()

		s.respond(w, r, http.StatusOK, tweet)
	}
}

func (s *server) handleTweetsFetch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value(ctxKeyUser).(*model.User)
		tweets, err := s.store.Tweet().FindBySubscription(currentUser.ID)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, tweets)
	}
}

func (s *server) handleUsersLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type result struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			ID: u.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(s.appKey)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		s.respond(w, r, http.StatusOK, &result{Token: tokenString})
	}
}
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				s.error(w, r, http.StatusUnauthorized, err)
				return
			}

			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		tknStr := c.Value

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return s.appKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				s.error(w, r, http.StatusUnauthorized, err)
				return
			}
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if !tkn.Valid {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		u, err := s.store.User().Find(claims.ID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))

	}
}

func (s *server) handleUsersRegister() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Username: req.Username,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
