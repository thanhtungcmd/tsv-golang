package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tsv-golang/internal/filter"
	"tsv-golang/internal/repository"
	"tsv-golang/internal/route"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load env
	loadEnv()

	// load persisten
	repo := loadPersisten()

	// load api
	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	apiV1.Use(filter.TrackIdFilter)
	apiV1.Use(filter.AuthFilter)
	route.HandleApiV1(apiV1, *repo)

	// server
	s := &http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(any(err))
	}
	fmt.Printf("<------Server started------>")
}

func loadPersisten() *repository.Repositories {
	repo, err := repository.NewRepositories() // load repositories
	if err != nil {
		panic(any(err))
	}

	return repo
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
		panic(any("[EMERGENCY] Can't load .env file"))
	}

	gin.SetMode(os.Getenv("GIN_MODE"))
}
