package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"streamgo/core"
	"streamgo/logger"
	"streamgo/route"
)

func main() {

	mux := http.NewServeMux()
	log := logger.New()

	service := core.NewService()

	route.Handle(mux, log, service)

	s := http.Server{
		Addr:              ":8000",
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 10,
	}

	log.Println("Stream server started :8000")
	go func() {
		if e := s.ListenAndServe(); e != nil {
			log.Fatal(e)
		}

	}()

	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Interrupt)
	signal.Notify(sigC, os.Kill)

	log.Println("Terminated", <-sigC)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
