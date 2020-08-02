package graph

//go:generate go run github.com/99designs/gqlgen

import (
	giteaCtx "code.gitea.io/gitea/modules/context"
)
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	giteaApiContext *giteaCtx.APIContext
}

