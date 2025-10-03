package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tsv-golang/internal/filter"
	"tsv-golang/internal/graph"
	"tsv-golang/internal/persistence"
	"tsv-golang/internal/route"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
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
	router.POST("/query", graphqlHandler())
	router.GET("/", playgroundHandler())

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

func loadPersisten() *persistence.Repositories {
	repo, err := persistence.NewRepositories() // load repositories
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

func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Server setup:
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
