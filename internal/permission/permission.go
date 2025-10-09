package permission

import (
	"context"
	"tsv-golang/internal/repository"

	"github.com/99designs/gqlgen/graphql"
)

func HasPermission(repo *repository.Repositories) func(ctx context.Context, obj interface{}, next graphql.Resolver, action string) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, action string) (interface{}, error) {
		return next(ctx)
	}
}
