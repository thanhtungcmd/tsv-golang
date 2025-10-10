package route

import (
	"tsv-golang/internal/direction"
	"tsv-golang/internal/graph"
	"tsv-golang/internal/repository"
	service2 "tsv-golang/internal/service"

	handlerGraph "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func HandleApiV1(route *gin.RouterGroup, repo repository.Repositories) {
	authHandler := ApiAuthHandlerInit(repo)
	route.POST("auth/login", authHandler.Login)

	route.POST("/query", graphqlHandler(repo))
	route.GET("/", playgroundHandler())
}

func graphqlHandler(repo repository.Repositories) gin.HandlerFunc {
	service := service2.NewService(&repo)

	cfg := graph.Config{Resolvers: &graph.Resolver{
		Repositories: &repo,
		Service:      service,
	}}
	cfg.Directives.HasPermission = direction.HasPermission(&repo)
	cfg.Directives.Validate = direction.Validate

	h := handlerGraph.New(graph.NewExecutableSchema(cfg))

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
