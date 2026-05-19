package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)



func (app *application) mount() http.Handler {
   
   r := chi.NewRouter()

  // A good base middleware stack
  r.Use(middleware.RequestID)
  r.Use(middleware.RealIP)
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)

  // Set a timeout value on the request context (ctx), that will signal
  // through ctx.Done() that the request has timed out and further
  // processing should be stopped.
  r.Use(middleware.Timeout(60 * time.Second))

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hi"))
  })
}

type application struct {
	config config
}
type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}