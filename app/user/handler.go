package user

import "github.com/themethaithian/nethttp/app"

type Handler interface {
	CreateUser(ctx app.Context)
}

type handler struct{}

func NewHandler() Handler {
	return &handler{}
}
