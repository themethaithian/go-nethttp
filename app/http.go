package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type contextHTTP struct {
	w          http.ResponseWriter
	r          *http.Request
	logHandler slog.Handler
}

func NewContextHttp(w http.ResponseWriter, r *http.Request) Context {
	return &contextHTTP{
		w: w,
		r: r,
	}
}

func (c *contextHTTP) Bind(v any) error {
	defer c.r.Body.Close()
	return json.NewDecoder(c.r.Body).Decode(v)
}

func (c *contextHTTP) Param(key string) string {
	return c.r.PathValue(key)
}

func (c *contextHTTP) OK(v any) {
	c.w.WriteHeader(http.StatusOK)
	if v == nil {
		return
	}

	err := json.NewEncoder(c.w).Encode(Response{
		Status: Success,
		Data:   v,
	})
	_ = err
}

func (c *contextHTTP) BadRequest(err error) {
	c.w.WriteHeader(http.StatusBadRequest)
	jsonErr := json.NewEncoder(c.w).Encode(Response{
		Status:  Fail,
		Message: err.Error(),
	})
	_ = jsonErr
}

func (c *contextHTTP) StoreError(err error) {
	c.w.WriteHeader(storeErrorStutas)
	jsonErr := json.NewEncoder(c.w).Encode(Response{
		Status:  Fail,
		Message: err.Error(),
	})
	_ = jsonErr
}

type RouterHTTP struct {
	mux          *http.ServeMux
	logger       *slog.Logger
	interceptors []middlewareFunc
}

func NewRouterHTTP() *RouterHTTP {
	r := http.NewServeMux()

	return &RouterHTTP{mux: r}
}

type middlewareFunc func(h http.Handler) http.Handler

func (router *RouterHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (r *RouterHTTP) Use(h ...middlewareFunc) {
	r.interceptors = append(r.interceptors, h...)
}

func (r *RouterHTTP) GET(path string, handler func(Context)) {
	r.mux.Handle(path, NewHTTPHandler(http.MethodGet, handler, r.interceptors, r.logger))
}

func (r *RouterHTTP) POST(path string, handler func(Context)) {
	r.mux.Handle(path, NewHTTPHandler(http.MethodPost, handler, r.interceptors, r.logger))
}

func NewHTTPHandler(method string, handler func(Context), interceptors []middlewareFunc, logger *slog.Logger) http.Handler {
	var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler(&contextHTTP{w: w, r: r, logHandler: logger.Handler().WithAttrs([]slog.Attr{slog.String("transaction-id", r.Header.Get("transaction-id"))})})
	})

	for _, interceptor := range interceptors {
		httpHandler = interceptor(http.Handler(httpHandler))
	}
	return httpHandler
}
