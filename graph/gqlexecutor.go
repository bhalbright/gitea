package graph

import (
	"code.gitea.io/gitea/graph/generated"
	giteaCtx "code.gitea.io/gitea/modules/context"

	handler2 "github.com/99designs/gqlgen/graphql/handler"
)

func GraphQL(ctx *giteaCtx.APIContext) {
	config := generated.Config{
		Resolvers: &Resolver{
			giteaApiContext: ctx,
		},
	}
	handler := handler2.NewDefaultServer(generated.NewExecutableSchema(config))
	handler.ServeHTTP(ctx.Resp, ctx.Req.Request)
}
