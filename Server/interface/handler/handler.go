package _interface

import (
	"example.com/application/service"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: *userService}
}

// LoginHandle プレイヤーログイン
func (u *UserHandler) LoginHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {

		return nil
	}
}

// MoveHandle プレイヤー移動同期
func (u *UserHandler) MoveHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {

		return nil
	}
}

// OtherPlayerHandle 他のプレイヤーの情報同期
func (u *UserHandler) OtherPlayerHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

// DestroyHandle プレイヤーゲームオーバー
func (u *UserHandler) DestroyHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}
