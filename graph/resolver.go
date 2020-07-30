package graph

import (
	"code.gitea.io/gitea/graph/model"
	giteaCtx "code.gitea.io/gitea/modules/context"
)
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	todos []*model.Todo
	giteaApiContext *giteaCtx.APIContext
}

