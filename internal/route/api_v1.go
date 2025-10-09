package route

import (
	"tsv-golang/internal/graph"
	"tsv-golang/internal/handler"
	"tsv-golang/internal/persistence"

	handlerGraph "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func HandleApiV1(route *gin.RouterGroup, repo persistence.Repositories) {
	authHandler := handler.ApiAuthHandlerInit(repo)
	route.POST("auth/login", authHandler.Login)

	route.POST("/query", graphqlHandler(repo))
	route.GET("/", playgroundHandler())
}

func graphqlHandler(repo persistence.Repositories) gin.HandlerFunc {
	h := handlerGraph.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		&repo,
	}}))

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

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/api/v1/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
