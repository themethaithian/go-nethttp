package app

type Context interface {
	Bind(v any) error
	Param(key string) string
	OK(v any)
	BadRequest(err error)
	StoreError(err error)
}

type HandlerFunc func(Context)
