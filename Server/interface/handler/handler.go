package _interface

import (
	"github.com/uptrace/bunrouter"
	"net/http"
)

func BunHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}
