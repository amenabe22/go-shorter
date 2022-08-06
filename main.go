package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "shortner/api"
	"shortner/core"
	rr "shortner/repository/redis"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func getRepo() core.RedirectRepository {
	// env variables
	// redisURL := os.Getenv("REDIS_URL")
	redisURL := "redis://127.0.0.1:6379"
	repo, err := rr.NewRedisRepository(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func main() {
	repo := getRepo()
	service := core.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	// added core middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	// create buffered channel with 2 buffer space
	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8080")
		errs <- http.ListenAndServe("127.0.0.1:8080", r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}
