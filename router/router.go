package router

import "net/http"

type Router interface{
  Route() *http.ServeMux
}

type router struct{}

func NewRouter() Router {
	return &router{}
}

func (r *router) Route() *http.ServeMux {
  handler := http.NewServeMux()

  return handler
}
