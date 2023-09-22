package middleware

import (
	"errors"
	"example.com/application/auth"
	"example.com/application/service"
	"fmt"
	"github.com/uptrace/bunrouter"
	"log"
	"net/http"
)

type Middleware struct {
	UserService *service.UserService
}

func NewMiddleware(userService *service.UserService) *Middleware {
	return &Middleware{
		UserService: userService,
	}
}

func (m *Middleware) AuthenticateMiddleware() func(bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			ctx := req.Context()
			token := req.Header.Get("x-token")
			if token == "" {
				return errors.New("x-token is empty")
			}

			user, err := m.UserService.GetUserByAuthToken(ctx, token)
			if err != nil {
				return err
			}
			if user == nil {
				return fmt.Errorf("user not found. token=%s", token)
			}

			ctx = auth.SetUserID(ctx, user.Id)
			req = req.WithContext(ctx)
			return next(w, req)
		}
	}
}

func (m *Middleware) CorsMiddleware() func(bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", "*")
			h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

			if req.Method == http.MethodOptions {
				return nil
			}
			return next(w, req)
		}
	}
}

func (m *Middleware) RecoverMiddleware() func(bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("recovered from panic: %v", r)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			return next(w, req)
		}
	}
}
