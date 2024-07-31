package main

import (
	"fmt"
	"net/http"

	"github.com/themethaithian/nethttp/router"
)

func main() {
  router := router.NewRouter()

  server := http.Server{
    Addr: ":8080",
    Handler: router.Route(),
  }

  fmt.Printf("Server listening on port :8080")
  err := server.ListenAndServe()
  if err != nil {
    panic(err)
  }
}
