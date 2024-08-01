package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/themethaithian/nethttp/app"
	"github.com/themethaithian/nethttp/app/user"
	"github.com/themethaithian/nethttp/config"
)

func main() {
	r := app.NewRouterHTTP()

	userHandler := user.NewHandler()

	r.POST("/create_user", userHandler.CreateUser)

	server := http.Server{
		Addr:              ":" + config.Val.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		d := time.Duration(5 * time.Second)
		fmt.Printf("shutting down init %s ...", d)
		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			fmt.Printf("HTTP server Shutdown: %s", err.Error())
		}
		close(idleConnsClosed)
	}()

	fmt.Println(":" + config.Val.Port + " is serve")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("HTTP server ListenAndServe: %s", err.Error())
		return
	}

	<-idleConnsClosed
	fmt.Println("gracefully")
}
