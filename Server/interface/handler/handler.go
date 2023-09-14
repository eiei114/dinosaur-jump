package _interface

import (
	"github.com/uptrace/bunrouter"
	"net/http"
)

// LoginHandle プレイヤーログイン
func LoginHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

// MoveHandle プレイヤー移動同期
func MoveHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

// OtherPlayerHandle 他のプレイヤーの情報同期
func OtherPlayerHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

// DestroyHandle プレイヤーゲームオーバー
func DestroyHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}
